package services

import (
	"fmt"

	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	messageRepository *repositories.MessageRepository
	matchRepository   *repositories.MatchRepository
	userRepository    *repositories.UserRepository
	notifService      NotificationService
}

func NewMessageService(msgRepo *repositories.MessageRepository, matchRepo *repositories.MatchRepository, userRepo *repositories.UserRepository, notifService NotificationService) *MessageService {
	return &MessageService{
		messageRepository: msgRepo,
		matchRepository:   matchRepo,
		userRepository:    userRepo,
		notifService:      notifService,
	}
}

// SendMessage ensures the match exists, then saves the message
func (ms *MessageService) SendMessage(matchID, senderID primitive.ObjectID, content, messageType string) (*models.Message, error) {
	// Default to text if not provided
	if messageType == "" {
		messageType = "text"
	}
	
	newMsg := models.Message{
		ID:       primitive.NewObjectID(),
		MatchID:  matchID,
		SenderID: senderID,
		Content:  content,
		Type:     messageType,
		IsRead:   false,
	}

	err := ms.messageRepository.CreateMessage(newMsg)
	if err != nil {
		return nil, err
	}

	return &newMsg, nil
}

// SendMessageToTargetUser finds the match between sender and target, and sends the message
func (ms *MessageService) SendMessageToTargetUser(senderID, targetID primitive.ObjectID, content, messageType string) (*models.Message, error) {
	// 1. Find the match between these two users
	match, err := ms.matchRepository.FindMatchBetweenUsers(senderID, targetID)
	if err != nil {
		return nil, fmt.Errorf("match not found between users")
	}

	// 2. Save the message using the found matchID
	newMsg, err := ms.SendMessage(match.ID, senderID, content, messageType)
	if err != nil {
		return nil, err
	}

	// 3. Send Push Notification (Fire and forget via Goroutine so it doesn't block WS)
	go ms.sendPushNotification(senderID, targetID, content, messageType)

	return newMsg, nil
}

func (ms *MessageService) sendPushNotification(senderID, targetID primitive.ObjectID, content, messageType string) {
	// 1. Get target user to find their PushToken
	targetUser, err := ms.userRepository.FindUserByID(targetID)
	if err != nil || targetUser.PushToken == "" {
		return // No token, can't push
	}

	// 2. Get sender user for the notification title
	senderUser, err := ms.userRepository.FindUserByID(senderID)
	if err != nil {
		return
	}

	// 3. Prepare message body
	bodyText := content
	if messageType == "prayer" {
		bodyText = "🙏 Sent you a prayer"
	} else if messageType == "scripture" {
		bodyText = "📖 Shared a scripture with you"
	}

	// 4. Delegate to the highly-scalable NotificationService interface
	ms.notifService.SendPush(targetUser.PushToken, senderUser.FirstName, bodyText, "/app/chat")
}

// GetChatHistory fetches old messages for a chat screen
func (ms *MessageService) GetChatHistory(matchID primitive.ObjectID) ([]models.Message, error) {
	return ms.messageRepository.GetMessagesForMatch(matchID)
}

// GetChatHistoryByTarget fetches old messages between two users
func (ms *MessageService) GetChatHistoryByTarget(user1, user2 primitive.ObjectID) ([]models.Message, error) {
	match, err := ms.matchRepository.FindMatchBetweenUsers(user1, user2)
	if err != nil {
		// If no match found, they just haven't talked yet, return empty array instead of error
		return []models.Message{}, nil
	}
	return ms.messageRepository.GetMessagesForMatch(match.ID)
}
