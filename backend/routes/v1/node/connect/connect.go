package connect

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Tag   string `json:"tag"`
	Token string `json:"token"`
}

// Route: /node/connect
func Connect(c *fiber.Ctx) error {

	// Parse request
	var req connectRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if !verify.InfoLocals(c).HasPermission(verify.PermissionUseServices) {
		return util.FailedRequest(c, localization.ErrorNoPermission, nil)
	}

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	currentSessionId, err := verify.InfoLocals(c).GetSessionUUID()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	tk := req.Token

	var acc database.Account
	if err := database.DBConn.Preload("Sessions").Where("id = ?", accId).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorAccountNotFound, nil)
	}

	// Check if account has key set
	if database.DBConn.Where("id = ?", acc.ID).Find(&database.PublicKey{}).Error != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Get the most recent session
	var mostRecent database.Session = database.Session{
		LastConnection: time.Unix(0, 10),
	}
	var sessionIds []string
	for _, session := range acc.Sessions {
		sessionIds = append(sessionIds, session.ID.String())

		if session.LastConnection.After(mostRecent.LastConnection) {
			mostRecent = session
		}
	}

	var currentSession database.Session
	if err := database.DBConn.Where("id = ?", currentSessionId).Take(&currentSession).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	if currentSession.Token != tk {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Get the app
	var application database.App
	if err := database.DBConn.Where("tag = ?", req.Tag).Take(&application).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	// Get lowest load node
	var lowest database.Node

	// Connect to the same node if possible
	if mostRecent.Node != 0 {
		if err := database.DBConn.Model(&database.Node{}).Where("app_id = ? AND status = ? AND id = ?", application.ID, database.StatusStarted, mostRecent.Node).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorNotSetup, nil)
		}
	} else {
		if err := database.DBConn.Model(&database.Node{}).Where("app_id = ? AND status = ?", application.ID, database.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorNotSetup, nil)
		}
	}

	// Ping node (to see if it's online)
	if err := lowest.SendPing(); err != nil {

		// Set the node to error
		nodes.TurnOff(&lowest, database.StatusError)
		return util.FailedRequest(c, localization.ErrorNode, err)
	}

	// Generate a jwt token for the node
	token, err := util.ConnectionToken(accId, currentSessionId.String(), lowest.ID)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate a jwt token with session information
	sessionInformationToken, err := util.SessionInformationToken(accId, sessionIds)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	currentSession.LastConnection = time.Now()
	currentSession.Node = lowest.ID
	currentSession.App = application.ID
	if err := database.DBConn.Save(&currentSession).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Save node
	if err := database.DBConn.Save(&lowest).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
		"id":      lowest.ID,
		"token":   token,
		"s_info":  sessionInformationToken,
	})
}
