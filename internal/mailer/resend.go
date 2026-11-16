package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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

	// Sandbox mode (enabled whenever ENV is not production) skips the Resend
	// call entirely and only logs the email that would have been sent.
	//
	// This is needed because the Resend account has no verified sending domain,
	// so it can only deliver to the API key owner's address; sending to any
	// other recipient would fail and abort the registration flow.
	//
	// Remove or relax this once a domain is verified in Resend.

	// TODO: DELETE COMMENT AND UNCOMMENT SWITCH BELOW WHEN DOMAIN IS VERIFIED
	// COMMENTED -> ALLOWS SENDING TO API HOLDER'S EMAIL ONLY
	// UNCOMMENTED -> PRODUCTION ENVIRONMENT: ALLOWS SENDING TO ANY EMAIL
	// 			   -> DEVELOPMENT ENVIRONMENT: SANDBOX MODE (LOGS EMAIL INSTEAD OF SENDING)
	if isSandbox {
		log.Printf("[sandbox] would send %q to %s", subject.String(), email)
		return 200, nil
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", FromName, m.fromEmail),
		To:      []string{email},
		Html:    body.String(),
		Subject: subject.String(),
	}

	var retryErr error

	for i := 0; i < maxRetries; i++ {
		_, retryErr = m.client.Emails.Send(params)
		if retryErr != nil {
			time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
			continue
		}
		return 200, nil
	}
	return -1, fmt.Errorf("failed to send email after %d attempts, error: %v", maxRetries, retryErr)
}
