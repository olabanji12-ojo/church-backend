package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SeedUserData struct {
	FirstName      string
	LastName       string
	Email          string
	Password       string
	Gender         string
	InterestedIn   string
	DOBStr         string
	Denomination   string
	ChurchAssembly string
	ChurchFreq     string
	PrayerFreq     string
	BibleFreq      string
	Intention      string
	Genotype       string
	Bio            string
	Photos         []string
	ScenarioAnswers map[string]string
}

func main() {
	log.Println("🚀 Seeding 5 Male and 5 Female profiles for Covenant defense walkthrough...")

	config.LoadEnv()
	db := database.ConnectDB()
	collection := db.Collection("users")
	embedService := services.NewEmbeddingService()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash default password: %v", err)
	}

	seedUsers := []SeedUserData{
		// --- 5 FEMALES ---
		{
			FirstName:      "Grace",
			LastName:       "Adebayo",
			Email:          "grace.adebayo@covenant.app",
			Password:       "Password123!",
			Gender:         "female",
			InterestedIn:   "male",
			DOBStr:         "1998-04-12",
			Denomination:   "Pentecostal",
			ChurchAssembly: "RCCG City of David",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AA",
			Bio:            "Passionate about youth ministry, worship music, and living a faith-centered life. Looking for a God-fearing partner to build a home with.",
			Photos:         []string{"https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_a",
				"q_stewardship_1": "opt_b",
				"q_pacing_1":      "opt_b",
			},
		},
		{
			FirstName:      "Joy",
			LastName:       "Okonkwo",
			Email:          "joy.okonkwo@covenant.app",
			Password:       "Password123!",
			Gender:         "female",
			InterestedIn:   "male",
			DOBStr:         "1999-09-25",
			Denomination:   "Anglican",
			ChurchAssembly: "House on the Rock",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Weekly",
			Intention:      "Marriage",
			Genotype:       "AA",
			Bio:            "Architect by day, choir member by heart. Deeply rooted in grace, love, and purposeful kingdom living.",
			Photos:         []string{"https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_b",
				"q_stewardship_1": "opt_a",
				"q_pacing_1":      "opt_a",
			},
		},
		{
			FirstName:      "Esther",
			LastName:       "Danjuma",
			Email:          "esther.danjuma@covenant.app",
			Password:       "Password123!",
			Gender:         "female",
			InterestedIn:   "male",
			DOBStr:         "1997-01-18",
			Denomination:   "Baptist",
			ChurchAssembly: "Daystar Christian Centre",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AS",
			Bio:            "Software engineer & Sunday school teacher. Seeking an intentional Christian gentleman for covenant marriage.",
			Photos:         []string{"https://images.unsplash.com/photo-1524504388940-b1c1722653e1?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_b",
				"q_conflict_1":     "opt_a",
				"q_stewardship_1": "opt_b",
				"q_pacing_1":      "opt_b",
			},
		},
		{
			FirstName:      "Blessing",
			LastName:       "Nwachukwu",
			Email:          "blessing.nwachukwu@covenant.app",
			Password:       "Password123!",
			Gender:         "female",
			InterestedIn:   "male",
			DOBStr:         "2000-06-30",
			Denomination:   "Pentecostal",
			ChurchAssembly: "Elevation Church",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AA",
			Bio:            "Medical practitioner. Passionate about community outreach, gospel music, and spiritual growth with a godly spouse.",
			Photos:         []string{"https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_a",
				"q_stewardship_1": "opt_a",
				"q_pacing_1":      "opt_c",
			},
		},
		{
			FirstName:      "Hannah",
			LastName:       "Bello",
			Email:          "hannah.bello@covenant.app",
			Password:       "Password123!",
			Gender:         "female",
			InterestedIn:   "male",
			DOBStr:         "1996-11-05",
			Denomination:   "Methodist",
			ChurchAssembly: "Covenant Nation",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Weekly",
			Intention:      "Marriage",
			Genotype:       "AS",
			Bio:            "Event planner and interior designer. Believer in prayer, honor, and intentional kingdom relationships.",
			Photos:         []string{"https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_c",
				"q_stewardship_1": "opt_b",
				"q_pacing_1":      "opt_a",
			},
		},

		// --- 5 MALES ---
		{
			FirstName:      "David",
			LastName:       "Emmanuel",
			Email:          "david.emmanuel@covenant.app",
			Password:       "Password123!",
			Gender:         "male",
			InterestedIn:   "female",
			DOBStr:         "1995-03-14",
			Denomination:   "Pentecostal",
			ChurchAssembly: "RCCG City of David",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AA",
			Bio:            "Financial analyst & worship leader. Loving God, honoring family, and seeking a godly partner for life.",
			Photos:         []string{"https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_a",
				"q_stewardship_1": "opt_b",
				"q_pacing_1":      "opt_b",
			},
		},
		{
			FirstName:      "Samuel",
			LastName:       "Ezekiel",
			Email:          "samuel.ezekiel@covenant.app",
			Password:       "Password123!",
			Gender:         "male",
			InterestedIn:   "female",
			DOBStr:         "1996-08-10",
			Denomination:   "Baptist",
			ChurchAssembly: "Daystar Christian Centre",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AS",
			Bio:            "Civil engineer & media team volunteer. Believes in purposeful leadership, open communication, and grace.",
			Photos:         []string{"https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_b",
				"q_conflict_1":     "opt_a",
				"q_stewardship_1": "opt_b",
				"q_pacing_1":      "opt_a",
			},
		},
		{
			FirstName:      "Daniel",
			LastName:       "Ogundipe",
			Email:          "daniel.ogundipe@covenant.app",
			Password:       "Password123!",
			Gender:         "male",
			InterestedIn:   "female",
			DOBStr:         "1997-12-04",
			Denomination:   "Pentecostal",
			ChurchAssembly: "Elevation Church",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AA",
			Bio:            "Entrepreneur and tech enthusiast. Seeking a virtuous woman with whom to build a strong Christian home.",
			Photos:         []string{"https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_b",
				"q_stewardship_1": "opt_a",
				"q_pacing_1":      "opt_b",
			},
		},
		{
			FirstName:      "Caleb",
			LastName:       "Martins",
			Email:          "caleb.martins@covenant.app",
			Password:       "Password123!",
			Gender:         "male",
			InterestedIn:   "female",
			DOBStr:         "1994-07-22",
			Denomination:   "Anglican",
			ChurchAssembly: "House on the Rock",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Weekly",
			Intention:      "Marriage",
			Genotype:       "AA",
			Bio:            "Legal practitioner and Bible study coordinator. Focused on faith, integrity, and lifelong growth.",
			Photos:         []string{"https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_a",
				"q_stewardship_1": "opt_b",
				"q_pacing_1":      "opt_a",
			},
		},
		{
			FirstName:      "Joshua",
			LastName:       "Chukwu",
			Email:          "joshua.chukwu@covenant.app",
			Password:       "Password123!",
			Gender:         "male",
			InterestedIn:   "female",
			DOBStr:         "1998-02-14",
			Denomination:   "Methodist",
			ChurchAssembly: "Covenant Nation",
			ChurchFreq:     "Weekly",
			PrayerFreq:     "Daily",
			BibleFreq:      "Daily",
			Intention:      "Marriage",
			Genotype:       "AS",
			Bio:            "Product designer and youth mentor. Loving Christ, valuing honesty, and ready for a purposeful relationship.",
			Photos:         []string{"https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?auto=format&fit=crop&w=800&q=80"},
			ScenarioAnswers: map[string]string{
				"q_boundaries_1":   "opt_a",
				"q_conflict_1":     "opt_c",
				"q_stewardship_1": "opt_c",
				"q_pacing_1":      "opt_c",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	createdUsers := make(map[string]models.User)

	for _, uData := range seedUsers {
		dob, _ := time.Parse("2006-01-02", uData.DOBStr)

		userObj := models.User{
			Email:                 strings.ToLower(uData.Email),
			PasswordHash:          string(passwordHash),
			FirstName:             uData.FirstName,
			LastName:              uData.LastName,
			Gender:                uData.Gender,
			InterestedIn:          uData.InterestedIn,
			DateOfBirth:           dob,
			Denomination:          uData.Denomination,
			ChurchAssembly:        uData.ChurchAssembly,
			ChurchFreq:            uData.ChurchFreq,
			PrayerFreq:            uData.PrayerFreq,
			BibleFreq:             uData.BibleFreq,
			Intention:             uData.Intention,
			Genotype:              uData.Genotype,
			StrictGenotypeFilter: false,
			MinAgePref:            18,
			MaxAgePref:            50,
			PreferredDenomination: "Any",
			PreferredChurchFreq:  "Any",
			Bio:                   uData.Bio,
			Photos:                uData.Photos,
			ScenarioAnswers:       uData.ScenarioAnswers,
			IsVerified:            true,
			IsGuest:               false,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		// Compute Embeddings
		profText := embedService.GenerateUserText(&userObj)
		if profEmb, err := embedService.GetEmbedding(profText); err == nil {
			userObj.ProfileEmbedding = profEmb
		} else {
			log.Printf("⚠️ Embedding failed for %s: %v", userObj.Email, err)
		}

		partnerText := embedService.GeneratePartnerPreferenceText(&userObj)
		if partnerEmb, err := embedService.GetEmbedding(partnerText); err == nil {
			userObj.PartnerPrefEmbedding = partnerEmb
		} else {
			log.Printf("⚠️ Partner embedding failed for %s: %v", userObj.Email, err)
		}

		// Upsert into MongoDB by Email
		var existing models.User
		err := collection.FindOne(ctx, bson.M{"email": userObj.Email}).Decode(&existing)
		if err == nil {
			userObj.ID = existing.ID
			_, err = collection.ReplaceOne(ctx, bson.M{"_id": existing.ID}, userObj)
			if err != nil {
				log.Printf("❌ Failed to replace user %s: %v", userObj.Email, err)
			} else {
				log.Printf("🔄 Updated profile & embeddings for: %s %s (%s)", userObj.FirstName, userObj.LastName, userObj.Email)
			}
		} else {
			userObj.ID = primitive.NewObjectID()
			_, err = collection.InsertOne(ctx, userObj)
			if err != nil {
				log.Printf("❌ Failed to insert user %s: %v", userObj.Email, err)
			} else {
				log.Printf("✨ Registered new user & generated embeddings for: %s %s (%s)", userObj.FirstName, userObj.LastName, userObj.Email)
			}
		}

		createdUsers[userObj.Email] = userObj
		time.Sleep(300 * time.Millisecond) // HF rate limit delay
	}

	// RUN MATCHING TEST BETWEEN ALL MALES AND ALL FEMALES
	fmt.Println("\n==========================================================================")
	fmt.Println("📊 COVENANT ALGORITHM MATCH MATRIX (5 MALES vs 5 FEMALES)")
	fmt.Println("==========================================================================")

	males := []string{"david.emmanuel@covenant.app", "samuel.ezekiel@covenant.app", "daniel.ogundipe@covenant.app", "caleb.martins@covenant.app", "joshua.chukwu@covenant.app"}
	females := []string{"grace.adebayo@covenant.app", "joy.okonkwo@covenant.app", "esther.danjuma@covenant.app", "blessing.nwachukwu@covenant.app", "hannah.bello@covenant.app"}

	scenarioSvc := services.NewScenarioService()
	genotypeSvc := services.NewGenotypeService()

	for _, mEmail := range males {
		mUser := createdUsers[mEmail]
		fmt.Printf("\n👨🏻‍💼 MALE PROFILE: %s %s (%s, Genotype: %s, Assembly: %s)\n", mUser.FirstName, mUser.LastName, mUser.Denomination, mUser.Genotype, mUser.ChurchAssembly)
		fmt.Println("--------------------------------------------------------------------------")

		type MatchEval struct {
			FemaleName     string
			MatchScore     int
			GenotypeStatus string
			GenotypeWarn   string
			SharedBadges   []string
			Icebreaker     string
		}

		var evals []MatchEval

		for _, fEmail := range females {
			fUser := createdUsers[fEmail]
			score, sharedBadges, icebreaker := scenarioSvc.CalculateCompatibility(mUser, fUser)
			status, warn := genotypeSvc.EvaluateCompatibility(mUser.Genotype, fUser.Genotype)

			if status == "incompatible" {
				sharedBadges = append([]string{"🤝 Friendship & Fellowship Only"}, sharedBadges...)
			}

			evals = append(evals, MatchEval{
				FemaleName:     fmt.Sprintf("%s %s (%s, Genotype: %s)", fUser.FirstName, fUser.LastName, fUser.Denomination, fUser.Genotype),
				MatchScore:     score,
				GenotypeStatus: status,
				GenotypeWarn:   warn,
				SharedBadges:   sharedBadges,
				Icebreaker:     icebreaker,
			})
		}

		sort.Slice(evals, func(i, j int) bool {
			return evals[i].MatchScore > evals[j].MatchScore
		})

		for rank, e := range evals {
			fmt.Printf("   [%d] Candidate: %s\n", rank+1, e.FemaleName)
			fmt.Printf("       💚 Covenant Compatibility Score: %d%%\n", e.MatchScore)
			fmt.Printf("       🧬 Genotype Status: %s %s\n", strings.ToUpper(e.GenotypeStatus), e.GenotypeWarn)
			fmt.Printf("       🏷️ Shared Badges: %s\n", strings.Join(e.SharedBadges, ", "))
			fmt.Printf("       💡 Icebreaker Insight: %s\n\n", e.Icebreaker)
		}
	}

	fmt.Println("==========================================================================")
	fmt.Println("🎉 Seeding and Matching Verification Complete!")
	fmt.Println("==========================================================================")
}
