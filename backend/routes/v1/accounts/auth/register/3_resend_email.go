package register_routes

import (
	"time"

	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/mail"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/resend_email (SSR)
func resendEmail(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate the token
	state, msg := validateToken(req.Token, 3)
	if msg != nil {
		return integration.FailedRequest(c, msg, nil)
	}

	// Rate limit the amount of email sending
	if time.Since(state.LastEmail) < time.Minute*5 {
		duration := time.Minute*5 - time.Since(state.LastEmail)
		return c.JSON(ssr.PopupResponse(c, localization.DialogTitleError, localization.AuthRegisterCodeEmailCooldown(int64(duration.Seconds()))))
	}

	// Validate the email
	valid, normalizedEmail := standards.CheckEmail(req.Email)
	if !valid {
		return integration.FailedRequest(c, localization.ErrorEmailInvalid, nil)
	}

	// Resend the email with a new code
	state.Mutex.Lock()
	state.Email = normalizedEmail
	state.EmailCode = auth.GenerateToken(6)
	state.LastEmail = time.Now()
	state.Mutex.Unlock()
	if err := mail.SendEmail(state.Email, localization.Locale(c), mail.EmailVerification, state.EmailCode); err != nil {
		return integration.FailedRequest(c, localization.ErrorMail, err)
	}

	// Open a popup telling the user the email was successfully resent
	return c.JSON(ssr.PopupResponse(c, localization.DialogTitleSuccess, localization.RegisterResendEmailDescription))
}
