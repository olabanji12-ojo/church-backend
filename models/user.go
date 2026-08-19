package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"password,omitempty"` // Omit from JSON responses for security
	FirstName    string             `bson:"first_name" json:"first_name"`
	LastName     string             `bson:"last_name" json:"last_name"`
	DateOfBirth  time.Time          `bson:"dob" json:"dob"`
	Gender       string             `bson:"gender" json:"gender"`
	InterestedIn string             `bson:"interested_in" json:"interested_in"`

	// Faith Profile
	Denomination   string `bson:"denomination" json:"denomination"`
	ChurchAssembly string `bson:"church_assembly" json:"church_assembly"`
	ChurchFreq     string `bson:"church_freq" json:"church_freq"`
	PrayerFreq   string `bson:"prayer_freq" json:"prayer_freq"`
	BibleFreq    string `bson:"bible_freq" json:"bible_freq"`
	Intention    string `bson:"intention" json:"intention"`
	// Preferences (FiltersScreen)
	MinAgePref           int    `bson:"min_age_pref" json:"min_age_pref"`
	MaxAgePref           int    `bson:"max_age_pref" json:"max_age_pref"`
	MaxDistance          int    `bson:"max_distance" json:"max_distance"` // kept for potential future use
	PreferredDenomination string `bson:"preferred_denomination" json:"preferred_denomination"`
	PreferredChurchFreq  string `bson:"preferred_church_freq" json:"preferred_church_freq"`

	// Media & Settings
	Photos       []string             `bson:"photos" json:"photos"`
	PushToken    string               `bson:"push_token" json:"push_token"` // For push notifications
	BlockedUsers []primitive.ObjectID `bson:"blocked_users" json:"blocked_users"`
	IsVerified   bool                 `bson:"is_verified" json:"is_verified"`
	IsGuest    bool      `bson:"is_guest" json:"is_guest"`     // Identifies temporary shadow users
	Bio        string    `bson:"bio" json:"bio"`
	ProfileEmbedding []float32 `bson:"profile_embedding,omitempty" json:"profile_embedding,omitempty"`
	PartnerPrefText      string    `bson:"partner_pref_text" json:"partner_pref_text"`
	PartnerPrefEmbedding []float32 `bson:"partner_pref_embedding,omitempty" json:"partner_pref_embedding,omitempty"`

	// Scenario Matching & Badges
	ScenarioAnswers map[string]string `bson:"scenario_answers,omitempty" json:"scenario_answers,omitempty"`
	UnlockedBadges  []string          `bson:"unlocked_badges,omitempty" json:"unlocked_badges,omitempty"`

	// Genotype & Medical Profile
	Genotype             string `bson:"genotype,omitempty" json:"genotype,omitempty"`                         // "AA", "AS", "AC", "SS", "SC", or "Unknown"
	StrictGenotypeFilter bool   `bson:"strict_genotype_filter,omitempty" json:"strict_genotype_filter,omitempty"` // Enforce medical safeguard filter

	// Calculated Dynamic Fields for Candidate Profiles
	MatchScore       int      `bson:"-" json:"match_score,omitempty"`
	SharedBadges     []string `bson:"-" json:"shared_badges,omitempty"`
	IcebreakerPrompt string   `bson:"-" json:"icebreaker_prompt,omitempty"`
	GenotypeStatus   string   `bson:"-" json:"genotype_status,omitempty"`  // "compatible", "incompatible", "unverified"
	GenotypeWarning  string   `bson:"-" json:"genotype_warning,omitempty"` // Medical warning message if incompatible

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
