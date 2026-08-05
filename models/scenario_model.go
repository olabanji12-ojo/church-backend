package models

type ScenarioOption struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	BadgeLabel string `json:"badge_label"`
}

type ScenarioQuestion struct {
	ID         string           `json:"id"`
	Pillar     string           `json:"pillar"`
	Question   string           `json:"question"`
	Options    []ScenarioOption `json:"options"`
	IsOnboard  bool             `json:"is_onboard"`
}

type AnswerSubmissionRequest struct {
	QuestionID string `json:"question_id"`
	OptionID   string `json:"option_id"`
}
