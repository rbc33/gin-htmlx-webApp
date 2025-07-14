package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailersend/mailersend-go"
	"github.com/rs/zerolog/log"
)

var APIKey = "MAILERSEND_API_KEY"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("Error loading .env file")
	}
	ms := mailersend.NewMailersend(APIKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subject := "Subject"
	text := "Greetings from the team, you got this message through MailerSend."
	html := "Greetings from the team, you got this message through MailerSend."

	from := mailersend.From{
		Name:  "Your Name",
		Email: os.Getenv("VERIFIED_SENDER_EMAIL"),
	}

	recipients := []mailersend.Recipient{
		{
			Name:  "Recipient",
			Email: os.Getenv("RECIPIENT_EMAIL"),
		},
	}

	variables := []mailersend.Variables{
		{
			Email: "recipient@email.com",
			Substitutions: []mailersend.Substitution{
				{
					Var:   "foo",
					Value: "bar",
				},
			},
		},
	}

	tags := []string{"foo", "bar"}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)
	message.SetSubstitutions(variables)
	message.SetTags(tags)

	res, _ := ms.Email.Send(ctx, message)

	fmt.Printf(res.Header.Get("X-Message-Id"))

}
