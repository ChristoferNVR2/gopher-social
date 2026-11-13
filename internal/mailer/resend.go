package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/resend/resend-go/v4"
)

type ResendMailer struct {
	fromEmail string
	apiKey    string
	client    *resend.Client
}

func NewResend(apiKey, fromEmail string) (*ResendMailer, error) {
	if apiKey == "" {
		return nil, errors.New("resend API key is required")
	}

	client := resend.NewClient(apiKey)
	return &ResendMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}, nil
}

func (m ResendMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	// TODO: Uncomment this when we are ready to send emails in production
	//if isSandbox {
	//	log.Printf("[sandbox] would send %q to %s", subject.String(), email)
	//	return nil
	//}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", FromName, m.fromEmail),
		To:      []string{email},
		Html:    body.String(),
		Subject: subject.String(),
	}

	var retryErr error

	for i := 0; i < maxRetries; i++ {
		_, err := m.client.Emails.Send(params)
		if err != nil {
			time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
			continue
		}
		return 200, nil
	}
	return -1, fmt.Errorf("failed to send email after %d attempts, error: %v", maxRetries, retryErr)
}
