package app

import (
	"fmt"
	"net/http"
	"net/mail"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"

	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"

	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
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
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			log.Error().Msgf("could not parse form %v", err)
			renderErrorPage(c, "unknown", fmt.Errorf("could not parse form: %w", err))
			return
		}

		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")
		recaptcha_response := c.Request.FormValue("g-recaptcha-response")

		// Make the request to Google's API only if user
		// configured recatpcha settings
		if (len(common.Settings.RecaptchaSecret) > 0) && (len(common.Settings.RecaptchaSiteKey) > 0) {
			ctx := c.Request.Context()
			err := verifyRecaptchaEnterprise(ctx, "gocms-1750166214215", common.Settings.RecaptchaSiteKey, recaptcha_response, "contact_submit")
			if err != nil {
				// El error ya se loguea dentro de verifyRecaptcha si es necesario o aquí en renderErrorPage
				renderErrorPage(c, email, err)
				return
			}
		}

		err := validateEmail(email)
		if err != nil {
			renderErrorPage(c, email, err)
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			renderErrorPage(c, email, fmt.Errorf("name too long (200 chars max)"))
			return
		}

		if len(message) > 10000 {
			// El mensaje de error original decía 1000, pero el código verifica 10000.
			renderErrorPage(c, email, fmt.Errorf("message too long (10000 chars max)"))
			return
		}

		if renderErr := TemplRender(c, http.StatusOK, views.MakeContactSuccess(email, name)); renderErr != nil {
			log.Error().Err(renderErr).Msg("Failed to render contact success page")
			c.String(http.StatusInternalServerError, "An error occurred while processing your request.")
		}
	}
}

// TODO : This is a duplicate of the index handler... abstract
func contactHandler(c *gin.Context, db database.Database) ([]byte, error) {
	return renderHtml(c, views.MakeContactPage(common.Settings.AppNavbar.Links, common.Settings.RecaptchaSiteKey, common.Settings.AppNavbar.Dropdowns))
}
