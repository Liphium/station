package register_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/mail"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Route: /account/auth/register/invite (SSR)
func checkInvite(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token  string `json:"token"`
		Invite string `json:"invite"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the token
	state, msg := validateToken(req.Token, 2)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Validate the invite
	inviteId, err := uuid.Parse(req.Invite)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorInviteInvalid, nil)
	}
	var invite account.Invite
	if err := database.DBConn.Where("id = ?", inviteId).Take(&invite).Error; err != nil {

		// Send an invalid invite error if the record wasn't found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.FailedRequest(c, localization.ErrorInviteInvalid, nil)
		}

		return util.FailedRequest(c, localization.ErrorServer, nil)
	}

	// Add the invite to the state
	state.Mutex.Lock()
	state.Invite = req.Invite
	state.EmailCode = auth.GenerateToken(6)
	state.Mutex.Unlock()

	// Send them an email code
	if err := mail.SendEmail(state.Email, localization.Locale(c), mail.EmailVerification, state.EmailCode); err != nil {
		return util.FailedRequest(c, localization.ErrorMail, err)
	}

	// Upgrade the token for the next step
	if msg := upgradeToken(req.Token, 3); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Return the email validate form
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterCodeTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Text{
			Text:  localization.RegisterCodeDescription,
			Style: ssr.TextStyleDescription,
		},
		ssr.Input{
			Placeholder: localization.AuthStartEmailPlaceholder,
			Value:       state.Email,
			Name:        "email",
		},
		ssr.Input{
			Placeholder: localization.RegisterCodePlaceholder,
			Name:        "code",
		},
		ssr.SubmitButton{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/register/email_code",
		},
		ssr.Button{
			Label: localization.AuthResendEmailButton,
			Path:  "/account/auth/register/resend_email",
		},
	}))
}
