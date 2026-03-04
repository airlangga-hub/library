package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (r *repository) SendEmail(to, subject, textPart string) error {
	url := "https://api.mailjet.com/v3.1/send"

	payload := MailjetRequest{
		Messages: []MessageRequest{
			{
				From: Person{
					Email: r.MailjetSender,
					Name:  "Library FTGO 14",
				},
				To: []Person{
					{
						Email: to,
						Name:  "Library Renter",
					},
				},
				Subject:  subject,
				TextPart: textPart,
			},
		},
	}
	
	ppayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("repo.SendEmail: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(ppayload))
	if err != nil {
		return fmt.Errorf("repo.SendEmail: %w", err)
	}

	req.SetBasicAuth(r.MailjetUsername, r.MailjetPassword)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("repo.SendEmail: %w", err)
	}
	defer res.Body.Close()
	
	var mailjetResp MailjetResponse
	if err := json.NewDecoder(req.Body).Decode(&mailjetResp); err != nil {
		return fmt.Errorf("repo.SendEmail: %w", err)
	}
	
	if mailjetResp.StatusCode >= 400 {
		return fmt.Errorf("repo.SendEmail: %s", mailjetResp.ErrorMessage)
	}

	return nil
}
