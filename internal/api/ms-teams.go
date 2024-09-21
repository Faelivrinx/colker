package api

// AdaptiveCard represents the full adaptive card JSON structure
type AdaptiveCard struct {
	Type        string       `json:"type"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment represents the attachment in the adaptive card
type Attachment struct {
	ContentType string  `json:"contentType"`
	Content     Content `json:"content"`
}

// Content represents the content of the attachment (the AdaptiveCard itself)
type Content struct {
	Type    string `json:"type"`
	Body    []Body `json:"body"`
	Schema  string `json:"$schema"`
	Version string `json:"version"`
}

// Body represents each block in the body of the AdaptiveCard, in this case a TextBlock
type Body struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook?tabs=newteams%2Cdotnet
func New(message string) *AdaptiveCard {
	return &AdaptiveCard{
		Type: "message",
		Attachments: []Attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: Content{
					Type: "AdaptiveCard",
					Body: []Body{
						{
							Type: "TextBlock",
							Text: message,
						},
					},
					Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
					Version: "1.0",
				},
			},
		},
	}
}
