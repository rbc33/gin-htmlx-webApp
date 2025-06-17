package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

const RECAPTCHA_VERIFY_URL string = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaResponse struct {
	Success   bool    `json:"success"`
	Score     float32 `json:"score"`
	Timestamp string  `json:"challenge_ts"`
	Hostname  string  `json:"hostname"`
	Action    string  `json:"action,omitempty"` // Añadir Action para v3
}

func verifyRecaptcha(recaptcha_secret string, recaptcha_response string) error {
	// Validate that the recaptcha response was actually
	// not a bot by checking the success rate
	recaptcha_response_data, err := http.PostForm(RECAPTCHA_VERIFY_URL, url.Values{
		"secret":   {recaptcha_secret},
		"response": {recaptcha_response},
	})
	if err != nil {
		err_str := fmt.Sprintf("could not do recaptcha post request: %s", err)
		return fmt.Errorf("%s: %s", err_str, err)
	}
	defer recaptcha_response_data.Body.Close()

	if recaptcha_response_data.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid recaptcha response: %s", recaptcha_response_data.Status)
	}
	var recaptcha_answer RecaptchaResponse
	recaptcha_response_data_buffer, _ := io.ReadAll(recaptcha_response_data.Body)
	err = json.Unmarshal(recaptcha_response_data_buffer, &recaptcha_answer)
	if err != nil {
		return fmt.Errorf("could not parse recaptcha response: %s", err)
	}

	// Para reCAPTCHA v3, `Success` debe ser true y `Score` debe estar por encima de un umbral.
	// También es recomendable verificar `Hostname` y `Action`.
	expectedAction := "contact_submit" // Debe coincidir con la acción en el frontend
	// Asegúrate de que common.Settings.AppDomain esté configurado si quieres validar el hostname.
	isHostnameMismatch := common.Settings.AppDomain != "" && recaptcha_answer.Hostname != common.Settings.AppDomain

	if !recaptcha_answer.Success ||
		recaptcha_answer.Score < 0.5 || // Ajusta el umbral según sea necesario (ej. 0.5 o 0.7)
		isHostnameMismatch ||
		recaptcha_answer.Action != expectedAction {
		log.Warn().Msgf("Validación de reCAPTCHA v3 fallida. Success: %v, Score: %.2f, Hostname: %s (esperado: %s), Action: %s (esperada: %s), Timestamp: %s",
			recaptcha_answer.Success,
			recaptcha_answer.Score,
			recaptcha_answer.Hostname,
			common.Settings.AppDomain, // Asegúrate que common.Settings.AppDomain está configurado
			recaptcha_answer.Action,
			expectedAction,
			recaptcha_answer.Timestamp)
		return fmt.Errorf("could not validate recaptcha") // Mensaje genérico para el usuario
	}
	log.Info().Msgf("Validación de reCAPTCHA v3 exitosa. Score: %.2f, Hostname: %s, Action: %s", recaptcha_answer.Score, recaptcha_answer.Hostname, recaptcha_answer.Action)
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
			err := verifyRecaptcha(common.Settings.RecaptchaSecret, recaptcha_response)
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
	return renderHtml(c, views.MakeContactPage(common.Settings.AppNavbar.Links, common.Settings.RecaptchaSiteKey))
}
