package connect

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
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

	if !util.Permission(c, util.PermissionUseServices) {
		return util.FailedRequest(c, "no.permission", nil)
	}

	// Get account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	currentSessionId := util.GetSession(c)
	tk := req.Token

	var acc account.Account
	if err := database.DBConn.Preload("Sessions").Where("id = ?", accId).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, "not.found", nil)
	}

	// Check if account has key set
	if database.DBConn.Where("id = ?", acc.ID).Find(&account.PublicKey{}).Error != nil {
		return util.FailedRequest(c, "no.key", nil)
	}

	// Get the most recent session
	var mostRecent account.Session = account.Session{
		LastConnection: time.Unix(0, 10),
	}
	var sessionIds []string
	for _, session := range acc.Sessions {
		sessionIds = append(sessionIds, session.ID)

		if session.LastConnection.After(mostRecent.LastConnection) {
			mostRecent = session
		}
	}

	var currentSession account.Session
	if err := database.DBConn.Where("id = ?", currentSessionId).Take(&currentSession).Error; err != nil {
		return util.FailedRequest(c, "not.found", err)
	}

	if currentSession.Token != tk {
		return util.FailedRequest(c, "invalid.token", nil)
	}

	// Get the app
	var application app.App
	if err := database.DBConn.Where("tag = ?", req.Tag).Take(&application).Error; err != nil {
		return util.FailedRequest(c, "invalid.app", err)
	}

	// Get lowest load node
	var lowest node.Node

	// Connect to the same node if possible
	if mostRecent.Node != 0 {
		if err := database.DBConn.Model(&node.Node{}).Where("app_id = ? AND status = ? AND id = ?", application.ID, node.StatusStarted, mostRecent.Node).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, "not.setup", nil)
		}
	} else {
		if err := database.DBConn.Model(&node.Node{}).Where("app_id = ? AND status = ?", application.ID, node.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, "not.setup", nil)
		}
	}

	// Ping node (to see if it's online)
	if err := lowest.SendPing(); err != nil {

		// Set the node to error
		nodes.TurnOff(&lowest, node.StatusError)
		return util.FailedRequest(c, util.ErrorNode, err)
	}

	// Generate a jwt token for the node
	token, err := util.ConnectionToken(accId, currentSessionId, lowest.ID)
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Generate a jwt token with session information
	sessionInformationToken, err := util.SessionInformationToken(accId, sessionIds)
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	currentSession.LastConnection = time.Now()
	currentSession.Node = lowest.ID
	currentSession.App = application.ID
	if err := database.DBConn.Save(&currentSession).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Save node
	if err := database.DBConn.Save(&lowest).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
		"id":      lowest.ID,
		"token":   token,
		"s_info":  sessionInformationToken,
	})
}
