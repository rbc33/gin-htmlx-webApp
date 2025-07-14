package app

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"

	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"

	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"github.com/mailersend/mailersend-go"
)

func verifyRecaptchaEnterprise(ctx context.Context, projectID, recaptchaKey, token, expectedAction string) error {
	client, err := recaptcha.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating reCAPTCHA client: %w", err)
	}
	defer client.Close()

	event := &recaptchapb.Event{
		Token:   token,
		SiteKey: recaptchaKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
	}

	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", projectID),
	}

	response, err := client.CreateAssessment(ctx, request)
	if err != nil {
		return fmt.Errorf("error calling CreateAssessment: %w", err)
	}

	if !response.TokenProperties.Valid {
		return fmt.Errorf("invalid token: %v", response.TokenProperties.InvalidReason)
	}

	if response.TokenProperties.Action != expectedAction {
		return fmt.Errorf("unexpected action: got %q, want %q", response.TokenProperties.Action, expectedAction)
	}

	log.Info().Msgf("reCAPTCHA Enterprise validation succeeded. Score: %.2f", response.RiskAnalysis.Score)
	for _, reason := range response.RiskAnalysis.Reasons {
		log.Info().Msgf("Risk reason: %s", reason.String())
	}

	// You can add threshold check here if desired, e.g.:
	if response.RiskAnalysis.Score < 0.5 {
		return fmt.Errorf("low reCAPTCHA score: %.2f", response.RiskAnalysis.Score)
	}

	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// renderErrorPage renderiza la página de error. HTMX reemplazará el contenido del target.
func renderErrorPage(c *gin.Context, email string, originalErr error) {
	log.Error().Err(originalErr).Str("email_attempted", email).Msg("Contact form submission error")
	if renderErr := TemplRender(c, http.StatusOK, views.MakeContactFailure(email, originalErr.Error())); renderErr != nil {
		log.Error().Err(renderErr).Msg("Failed to render contact failure page")
		// Fallback si renderizar la página de error específica falla
		c.String(http.StatusInternalServerError, "An error occurred processing your request, and we could not display the specific error message.")
	}
}

func makeContactFormHandler() func(*gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("Error loading .env file")
	}

	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			log.Error().Msgf("could not parse form %v", err)
			renderErrorPage(c, "unknown", fmt.Errorf("could not parse form: %w", err))
			return
		}

		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		subject := c.Request.FormValue("subject") // <<< Make sure your HTML form has this field
		message := c.Request.FormValue("message")
		recaptcha_response := c.Request.FormValue("g-recaptcha-response")

		// --- 1. reCAPTCHA Enterprise verification ---
		if (len(common.Settings.RecaptchaSecret) > 0) && (len(common.Settings.RecaptchaSiteKey) > 0) {
			ctx := c.Request.Context()
			err := verifyRecaptchaEnterprise(ctx, "gocms-1750166214215", common.Settings.RecaptchaSiteKey, recaptcha_response, "contact_submit")
			if err != nil {
				renderErrorPage(c, email, err)
				return
			}
		}

		// --- 2. Input Validations (Email, Name, Subject, Message) ---
		if err := validateEmail(email); err != nil {
			renderErrorPage(c, email, err)
			return
		}
		if len(name) > 200 {
			renderErrorPage(c, email, fmt.Errorf("name too long (200 chars max)"))
			return
		}
		if len(subject) == 0 {
			renderErrorPage(c, email, fmt.Errorf("subject cannot be empty"))
			return
		}
		if len(subject) > 250 { // Example limit for subject
			renderErrorPage(c, email, fmt.Errorf("subject too long (250 chars max)"))
			return
		}
		if len(message) > 10000 {
			renderErrorPage(c, email, fmt.Errorf("message too long (10000 chars max)"))
			return
		}

		// --- NEW: Email Sending Logic via MailerSend API ---
		mailersendAPIKey := os.Getenv("MAILERSEND_API_KEY")
		recipientEmail := os.Getenv("RECIPIENT_EMAIL")            // Your email address to receive contact messages
		verifiedSenderEmail := os.Getenv("VERIFIED_SENDER_EMAIL") // An email address you've verified in MailerSend

		if mailersendAPIKey == "" || recipientEmail == "" || verifiedSenderEmail == "" {
			log.Error().Msg("Missing one or more MailerSend environment variables. Please set MAILERSEND_API_KEY, RECIPIENT_EMAIL, and VERIFIED_SENDER_EMAIL.")
			renderErrorPage(c, email, fmt.Errorf("server email configuration is incomplete. Please try again later."))
			return
		}

		// Create an instance of the mailersend client
		ms := mailersend.NewMailersend(mailersendAPIKey)

		// Create a context with a timeout for the API call
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second) // Adjust timeout as needed
		defer cancel()                                                          // Ensures the context is cancelled when the function exits

		// Define the 'From' sender (must be a verified sender in MailerSend)
		from := mailersend.From{
			Name:  name,                // The user's name from the form
			Email: verifiedSenderEmail, // Your verified sender email in MailerSend
		}

		// Define the recipient (your email address)
		recipients := []mailersend.Recipient{
			{
				Email: recipientEmail,
				// Name: "Your Name for Contact Form", // Optional: If you want a specific name for yourself
			},
		}

		// Create the email message
		msg := ms.Email.NewMessage()

		msg.SetFrom(from)
		msg.SetRecipients(recipients)
		msg.SetSubject("GoCMS Contact: " + subject)
		msg.SetReplyTo(mailersend.Recipient{Email: email})

		htmlBody := fmt.Sprintf(`
            <p><strong>Name:</strong> %s</p>
            <p><strong>Email (from contact form):</strong> %s</p>
            <p><strong>Subject:</strong> %s</p>
            <hr>
            <p><strong>Message:</strong></p>
            <p style="white-space: pre-wrap;">%s</p>
        `, name, email, subject, message)
		msg.SetHTML(htmlBody)

		textBody := fmt.Sprintf("Name: %s\nEmail (from contact form): %s\nSubject: %s\n\nMessage:\n%s",
			name, email, subject, message)
		msg.SetText(textBody)

		res, err := ms.Email.Send(ctx, msg)
		if err != nil {
			bodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				log.Error().Err(err).Msg("Failed to read MailerSend API response body")
			} else {
				log.Error().Int("status_code", res.StatusCode).Bytes("body", bodyBytes).Msg("MailerSend API response error details")
			}
			res.Body.Close()

			log.Info().Str("from", fmt.Sprintf("%s <%s>", name, email)).Str("to", recipientEmail).Msgf("Contact form email sent successfully via MailerSend API (Message ID: %s)", res.Header.Get("X-Message-Id"))
			// --- END NEW: Email Sending Logic ---

			if renderErr := TemplRender(c, http.StatusOK, views.MakeContactSuccess(email, name)); renderErr != nil {
				log.Error().Err(renderErr).Msg("Failed to render contact success page")
				c.String(http.StatusInternalServerError, "An error occurred while processing your request.")
			}
		}
	}
}

// TODO : This is a duplicate of the index handler... abstract
func contactHandler(c *gin.Context, db database.Database) ([]byte, error) {
	return renderHtml(c, views.MakeContactPage(common.Settings.AppNavbar.Links, common.Settings.RecaptchaSiteKey, common.Settings.AppNavbar.Dropdowns))
}
