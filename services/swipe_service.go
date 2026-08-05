package services

import (
	"fmt"
	"sort"

	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SwipeService struct {
	userRepository    *repositories.UserRepository
	matchRepository   *repositories.MatchRepository
	messageRepository *repositories.MessageRepository
	notifService      NotificationService
}

func NewSwipeService(userRepo *repositories.UserRepository, matchRepo *repositories.MatchRepository, msgRepo *repositories.MessageRepository, notifService NotificationService) *SwipeService {
	return &SwipeService{
		userRepository:    userRepo,
		matchRepository:   matchRepo,
		messageRepository: msgRepo,
		notifService:      notifService,
	}
}

// SwipeRight handles the logic when a user likes another user
func (ss *SwipeService) SwipeRight(actorID, targetID primitive.ObjectID) (*models.Match, error) {
	// 1. Check if a match document already exists (Target swiped first?)
	existingMatch, err := ss.matchRepository.FindMatchBetweenUsers(actorID, targetID)
	
	if err == nil && existingMatch != nil {
		// Document exists. Was it pending? 
		if existingMatch.Status == "pending" {
			// IT'S A MATCH!
			err = ss.matchRepository.UpdateMatchStatus(existingMatch.ID, "matched")
			if err != nil {
				return nil, err
			}
			existingMatch.Status = "matched"
			
			// Trigger Match Notification in background
			go ss.sendMatchNotification(actorID, targetID)

			return existingMatch, nil
		}
		return existingMatch, nil // Already matched or rejected
	}

	// 2. No document exists. Create a new "pending" match
	newMatch := models.Match{
		ID:     primitive.NewObjectID(),
		Users:  []primitive.ObjectID{actorID, targetID},
		Status: "pending",
	}

	err = ss.matchRepository.CreateMatch(newMatch)
	if err != nil {
		return nil, err
	}

	// Trigger Like Notification in background
	go ss.sendLikeNotification(actorID, targetID)

	return &newMatch, nil
}

// SwipeLeft handles the logic when a user passes on another user
func (ss *SwipeService) SwipeLeft(actorID, targetID primitive.ObjectID) (*models.Match, error) {
	// 1. Check if a match document already exists (Target swiped first?)
	existingMatch, err := ss.matchRepository.FindMatchBetweenUsers(actorID, targetID)
	
	if err == nil && existingMatch != nil {
		// Update status to rejected
		err = ss.matchRepository.UpdateMatchStatus(existingMatch.ID, "rejected")
		if err != nil {
			return nil, err
		}
		existingMatch.Status = "rejected"
		return existingMatch, nil
	}

	// 2. No document exists. Create a new "rejected" match
	newMatch := models.Match{
		ID:     primitive.NewObjectID(),
		Users:  []primitive.ObjectID{actorID, targetID},
		Status: "rejected",
	}

	err = ss.matchRepository.CreateMatch(newMatch)
	if err != nil {
		return nil, err
	}

	return &newMatch, nil
}

// GetDiscoveryFeed returns a list of potential matches for the user
func (ss *SwipeService) GetDiscoveryFeed(userID primitive.ObjectID) ([]models.User, error) {
	// 1. Fetch all match documents the user is involved in
	interactions, err := ss.matchRepository.GetAllUserInteractions(userID)
	if err != nil {
		logrus.Warn("Could not fetch user interactions, proceeding without exclusions: ", err)
	}

	// 2. Extract the target user IDs from these interactions
	var excludedUserIDs []primitive.ObjectID
	for _, interaction := range interactions {
		// If the match is pending and the current user is NOT the initiator, do not exclude them!
		if interaction.Status == "pending" && len(interaction.Users) > 0 && interaction.Users[0] != userID {
			continue
		}
		for _, uID := range interaction.Users {
			if uID != userID {
				excludedUserIDs = append(excludedUserIDs, uID)
			}
		}
	}

	// 3. Fetch current user to get their preferences
	currentUser, err := ss.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 4. Try Vector Search first if user has embedding populated
	var users []models.User
	var vectorErr error
	if len(currentUser.ProfileEmbedding) == 384 {
		users, vectorErr = ss.userRepository.FindPotentialMatchesVector(currentUser, excludedUserIDs, 20)
	}

	// 5. Fallback to standard query if vector search is unavailable/fails, or user has no embedding
	if len(currentUser.ProfileEmbedding) != 384 || vectorErr != nil {
		logrus.Info("Falling back to standard database matching query.")
		users, err = ss.userRepository.FindPotentialMatches(currentUser, excludedUserIDs, 20)
		if err != nil {
			return nil, err
		}
	}

	// 6. Enrich candidate profiles with Scenario Match Scores, Shared Badges & Icebreakers
	scenarioSvc := NewScenarioService()
	for i := range users {
		matchScore, sharedBadges, icebreaker := scenarioSvc.CalculateCompatibility(*currentUser, users[i])
		users[i].MatchScore = matchScore
		users[i].SharedBadges = sharedBadges
		users[i].IcebreakerPrompt = icebreaker
	}

	// 7. Sort candidates in descending order of Covenant Match Score
	sort.Slice(users, func(i, j int) bool {
		return users[i].MatchScore > users[j].MatchScore
	})

	return users, nil
}

// GetMatches returns a list of MatchResponse models representing the user's successful matches and their last message
func (ss *SwipeService) GetMatches(userID primitive.ObjectID) ([]models.MatchResponse, error) {
	// 1. Get all match documents where status == "matched"
	matches, err := ss.matchRepository.GetUserMatches(userID)
	if err != nil {
		return nil, err
	}

	// 2. Extract the IDs of the matched users and map them to their Match document
	var matchedUserIDs []primitive.ObjectID
	matchMap := make(map[primitive.ObjectID]primitive.ObjectID) // Maps TargetUserID -> MatchID
	
	for _, match := range matches {
		for _, uID := range match.Users {
			if uID != userID {
				matchedUserIDs = append(matchedUserIDs, uID)
				matchMap[uID] = match.ID
			}
		}
	}

	// 3. If no matches, return empty array
	if len(matchedUserIDs) == 0 {
		return []models.MatchResponse{}, nil
	}

	// 4. Fetch the User objects
	users, err := ss.userRepository.GetUsersByIDs(matchedUserIDs)
	if err != nil {
		return nil, err
	}

	// 5. Construct MatchResponse array by fetching the last message for each match
	var responses []models.MatchResponse
	for _, user := range users {
		matchID := matchMap[user.ID]
		lastMsg, _ := ss.messageRepository.GetLastMessageForMatch(matchID) // Ignore error, nil is fine
		
		responses = append(responses, models.MatchResponse{
			User:        user,
			LastMessage: lastMsg,
		})
	}

	return responses, nil
}

// Unmatch removes a match and all of its associated chat history
func (ss *SwipeService) Unmatch(actorID, targetID primitive.ObjectID) error {
	// 1. Find the match document
	match, err := ss.matchRepository.FindMatchBetweenUsers(actorID, targetID)
	if err != nil {
		return err // Could be ErrNoDocuments
	}

	// 2. Delete all messages for this match
	err = ss.messageRepository.DeleteMessagesForMatch(match.ID)
	if err != nil {
		logrus.Errorf("Failed to delete messages for match %s: %v", match.ID.Hex(), err)
		// We continue anyway to ensure the match is deleted
	}

	// 3. Delete the match document itself
	return ss.matchRepository.DeleteMatch(match.ID)
}

func (ss *SwipeService) sendMatchNotification(actorID, targetID primitive.ObjectID) {
	// 1. Send push to targetID (notifying target user that they matched with actor)
	targetUser, err := ss.userRepository.FindUserByID(targetID)
	if err == nil && targetUser.PushToken != "" {
		actorUser, err := ss.userRepository.FindUserByID(actorID)
		if err == nil {
			ss.notifService.SendPush(
				targetUser.PushToken,
				"It's a Match! 💚",
				fmt.Sprintf("You and %s are now connected! Tap to chat.", actorUser.FirstName),
				"/app/matches",
			)
		}
	}

	// 2. Also send push to actorID (notifying actor user that they matched with target)
	actorUser, err := ss.userRepository.FindUserByID(actorID)
	if err == nil && actorUser.PushToken != "" {
		targetUser, err := ss.userRepository.FindUserByID(targetID)
		if err == nil {
			ss.notifService.SendPush(
				actorUser.PushToken,
				"It's a Match! 💚",
				fmt.Sprintf("You and %s are now connected! Tap to chat.", targetUser.FirstName),
				"/app/matches",
			)
		}
	}
}

func (ss *SwipeService) sendLikeNotification(actorID, targetID primitive.ObjectID) {
	targetUser, err := ss.userRepository.FindUserByID(targetID)
	if err != nil || targetUser.PushToken == "" {
		return
	}
	ss.notifService.SendPush(
		targetUser.PushToken,
		"New Connection Request 🙏",
		"Someone is interested in connecting with you. Open Church-Match to see!",
		"/app/discover",
	)
}
