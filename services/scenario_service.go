package services

import (
	"math"

	"github.com/olabanji12-ojo/church-backend/models"
)

type ScenarioService struct{}

func NewScenarioService() *ScenarioService {
	return &ScenarioService{}
}

// GetPredefinedScenarios returns the list of all covenant scenario questions
func (s *ScenarioService) GetPredefinedScenarios() []models.ScenarioQuestion {
	return []models.ScenarioQuestion{
		{
			ID:        "q_boundaries_1",
			Pillar:    "Moral Anchor & Boundaries",
			Question:  "A close colleague of the opposite sex invites you for a late-night private drink to discuss a bad day. How do you handle this?",
			IsOnboard: true,
			Options: []models.ScenarioOption{
				{ID: "opt_a", Text: "Decline or suggest coffee during office hours—late-night private drinks cross a personal boundary.", BadgeLabel: "🛡️ Proactive Boundaries"},
				{ID: "opt_b", Text: "Accept, but inform my partner immediately so everything is transparent.", BadgeLabel: "👁️ Open Transparency"},
				{ID: "opt_c", Text: "Go and focus on work—trust means not worrying about strict rules.", BadgeLabel: "🤝 High Trust Flexibility"},
			},
		},
		{
			ID:        "q_conflict_1",
			Pillar:    "Conflict & Grace",
			Question:  "When you and your partner have a heated disagreement before a social event, what is your instinct?",
			IsOnboard: true,
			Options: []models.ScenarioOption{
				{ID: "opt_a", Text: "Pause and address it calmly right away, even if it makes us late.", BadgeLabel: "🕊️ Immediate Resolution"},
				{ID: "opt_b", Text: "Put it on hold gracefully, enjoy the event, and discuss it in private later.", BadgeLabel: "⏳ Grateful Delay"},
				{ID: "opt_c", Text: "Take space alone first to process thoughts before speaking.", BadgeLabel: "💭 Thoughtful Space"},
			},
		},
		{
			ID:        "q_stewardship_1",
			Pillar:    "Stewardship & Finances",
			Question:  "When an unexpected financial bonus or blessing comes in, what is your immediate priority?",
			IsOnboard: true,
			Options: []models.ScenarioOption{
				{ID: "opt_a", Text: "Save or invest the majority to build long-term family security.", BadgeLabel: "⚖️ Future Security"},
				{ID: "opt_b", Text: "Set aside a portion for giving/tithing, and invest/save the rest.", BadgeLabel: "🌱 Kingdom Stewardship"},
				{ID: "opt_c", Text: "Use a portion for a shared experience or personal goal, then save the remainder.", BadgeLabel: "🎉 Balanced Joy"},
			},
		},
		{
			ID:        "q_pacing_1",
			Pillar:    "Pacing & Intentionality",
			Question:  "In the first 3 to 6 months of dating, what matters most to you in determining if this person is 'The One'?",
			IsOnboard: true,
			Options: []models.ScenarioOption{
				{ID: "opt_a", Text: "Observing how they handle stress, conflict, and treat service workers/family.", BadgeLabel: "🔍 Character Observer"},
				{ID: "opt_b", Text: "Experiencing deep, unfiltered conversations about faith, vision, and future goals.", BadgeLabel: "📖 Vision Alignment"},
				{ID: "opt_c", Text: "Seeing how well our daily rhythms, laughter, and companionship naturally flow together.", BadgeLabel: "☀️ Natural Harmony"},
			},
		},
		{
			ID:        "q_accountability_2",
			Pillar:    "Moral Anchor & Boundaries",
			Question:  "If you found yourself in a confusing or compromising situation by mistake, what is your immediate natural reaction?",
			IsOnboard: false,
			Options: []models.ScenarioOption{
				{ID: "opt_a", Text: "Tell my partner right away, even if it causes temporary hurt or awkwardness.", BadgeLabel: "💎 Uncompromising Honesty"},
				{ID: "opt_b", Text: "Handle it privately, cut off the situation, and only share if necessary.", BadgeLabel: "🔒 Private Resolution"},
				{ID: "opt_c", Text: "Speak to a pastor/mentor first to understand the situation before addressing it.", BadgeLabel: "📜 Wise Counsel Seeking"},
			},
		},
		{
			ID:        "q_forgiveness_2",
			Pillar:    "Conflict & Grace",
			Question:  "If your partner makes a genuine error or crosses a boundary but shows total repentance, where does your heart lean?",
			IsOnboard: false,
			Options: []models.ScenarioOption{
				{ID: "opt_a", Text: "Commit to working through it via counseling; true covenant means grace and restoration.", BadgeLabel: "✝️ Grace & Restoration"},
				{ID: "opt_b", Text: "Need strict space and time; forgiveness is immediate, but trust must be completely rebuilt.", BadgeLabel: "⏳ Patient Restoration"},
				{ID: "opt_c", Text: "A crossed boundary compromises the core covenant permanently.", BadgeLabel: "🛡️ Firm Non-Negotiable"},
			},
		},
	}
}

// CalculateCompatibility calculates overall score (0-100), shared badges, and icebreaker prompt
func (s *ScenarioService) CalculateCompatibility(userA, userB models.User) (int, []string, string) {
	scenarios := s.GetPredefinedScenarios()

	scenarioScoreSum := 0.0
	mutuallyAnsweredCount := 0
	var sharedBadges []string
	var topIcebreaker string

	// Map scenarios by ID for fast lookup
	scenarioMap := make(map[string]models.ScenarioQuestion)
	for _, q := range scenarios {
		scenarioMap[q.ID] = q
	}

	for qID, ansA := range userA.ScenarioAnswers {
		ansB, existsB := userB.ScenarioAnswers[qID]
		if !existsB {
			continue
		}

		mutuallyAnsweredCount++
		question, qExists := scenarioMap[qID]

		if ansA == ansB {
			scenarioScoreSum += 1.0

			// Extract badge label if available
			if qExists {
				for _, opt := range question.Options {
					if opt.ID == ansA {
						sharedBadges = append(sharedBadges, opt.BadgeLabel)
						if topIcebreaker == "" {
							topIcebreaker = "You both chose: '" + opt.Text + "'"
						}
						break
					}
				}
			}
		} else {
			// Partial compatibility score for non-identical answers
			scenarioScoreSum += 0.5

			// Build complementary teaser insight if topIcebreaker is empty
			if qExists && topIcebreaker == "" {
				var optAText, optBText string
				for _, opt := range question.Options {
					if opt.ID == ansA {
						optAText = opt.Text
					}
					if opt.ID == ansB {
						optBText = opt.Text
					}
				}
				if optAText != "" && optBText != "" {
					candidateName := userB.FirstName
					if candidateName == "" {
						candidateName = "your match"
					}
					topIcebreaker = "Complementary Insight: You chose '" + optAText + "', while " + candidateName + " chose '" + optBText + "'—both value proactive peace."
				}
			}
		}
	}

	scenarioScoreFinal := 75.0 // Base fallback score if no mutual scenarios answered yet
	if mutuallyAnsweredCount > 0 {
		scenarioScoreFinal = (scenarioScoreSum / float64(mutuallyAnsweredCount)) * 100.0
	}

	// Calculate Vector Cosine Similarity Score (0 - 100)
	vectorScoreFinal := 75.0
	if len(userA.PartnerPrefEmbedding) == 384 && len(userB.ProfileEmbedding) == 384 {
		simAB := float64(CosineSimilarity(userA.PartnerPrefEmbedding, userB.ProfileEmbedding))
		simBA := float64(CosineSimilarity(userB.PartnerPrefEmbedding, userA.ProfileEmbedding))
		mutualSim := (simAB + simBA) / 2.0
		vectorScoreFinal = math.Max(0.0, math.Min(100.0, mutualSim*100.0))
	}

	// Calculate Faith Attribute Compatibility Score (0 - 100)
	faithScoreFinal := 70.0
	if userA.Denomination != "" && userB.Denomination != "" {
		if userA.Denomination == userB.Denomination {
			faithScoreFinal += 15.0
		}
	}
	if userA.ChurchFreq != "" && userB.ChurchFreq != "" {
		if userA.ChurchFreq == userB.ChurchFreq {
			faithScoreFinal += 10.0
		}
	}
	if userA.ChurchAssembly != "" && userB.ChurchAssembly != "" {
		if userA.ChurchAssembly == userB.ChurchAssembly {
			faithScoreFinal += 10.0
		}
	}

	// Weighted Hybrid Score: 50% Scenario + 30% Vector + 20% Faith
	compositeScore := (0.50 * scenarioScoreFinal) + (0.30 * vectorScoreFinal) + (0.20 * faithScoreFinal)
	finalMatchScore := int(math.Round(compositeScore))

	if finalMatchScore > 99 {
		finalMatchScore = 99
	}
	if finalMatchScore < 50 {
		finalMatchScore = 50
	}

	if topIcebreaker == "" {
		topIcebreaker = "Ask about their perspective on shared faith values!"
	}

	return finalMatchScore, sharedBadges, topIcebreaker
}

func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dotProduct, normA, normB float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
