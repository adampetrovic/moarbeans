package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/adampetrovic/moarbeans/internal/database"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"regexp"
	"time"
)

type EmailPoller struct {
	db         *database.DB
	gmailSvc   *gmail.Service
	searchQuery string
}

func NewEmailPoller(db *database.DB) (*EmailPoller, error) {
	config, err := google.ConfigFromJSON([]byte(`{
		// Your Gmail API credentials here
	}`), gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	// TODO: Implement OAuth2 token retrieval
	client := config.Client(context.Background(), nil)
	
	svc, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	return &EmailPoller{
		db:         db,
		gmailSvc:   svc,
		searchQuery: "from:woodroaster.com.au subject:\"Magic link\"",
	}, nil
}

func (p *EmailPoller) Start() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		if err := p.pollEmails(); err != nil {
			fmt.Printf("Error polling emails: %v\n", err)
		}
	}
}

func (p *EmailPoller) pollEmails() error {
	msgs, err := p.gmailSvc.Users.Messages.List("me").Q(p.searchQuery).Do()
	if err != nil {
		return err
	}

	for _, msg := range msgs.Messages {
		message, err := p.gmailSvc.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue
		}

		link := p.extractMagicLink(message)
		if link != "" {
			if err := p.processLoginLink(link); err != nil {
				fmt.Printf("Error processing login link: %v\n", err)
			}
		}

		// Mark message as read
		p.gmailSvc.Users.Messages.Modify("me", msg.Id, &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{"UNREAD"},
		}).Do()
	}

	return nil
}

func (p *EmailPoller) extractMagicLink(message *gmail.Message) string {
	var body string
	for _, part := range message.Payload.Parts {
		if part.MimeType == "text/plain" {
			data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
			body = string(data)
			break
		}
	}

	// Extract link using regex (adjust pattern based on actual email format)
	re := regexp.MustCompile(`https://woodroaster\.com\.au/login/\S+`)
	matches := re.FindString(body)
	return matches
}

func (p *EmailPoller) processLoginLink(link string) error {
	// TODO: Implement logic to visit the link and extract session token
	// Store session in database
	return nil
} 