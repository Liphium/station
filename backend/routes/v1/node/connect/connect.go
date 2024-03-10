package connect

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/nodes"
	"time"

	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Cluster uint   `json:"cluster"`
	App     uint   `json:"app"`
	Token   string `json:"token"`
}

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
	accId := util.GetAcc(c)
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
		Token:          "-1",
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

	// Get lowest load node
	var lowest node.Node

	// Connect to the same node if possible
	if mostRecent.Node != 0 {
		if err := database.DBConn.Model(&node.Node{}).Where("cluster_id = ? AND app_id = ? AND status = ? AND id = ?", req.Cluster, req.App, node.StatusStarted, mostRecent.Node).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, "not.setup", nil)
		}
	} else {
		if err := database.DBConn.Model(&node.Node{}).Where("cluster_id = ? AND app_id = ? AND status = ?", req.Cluster, req.App, node.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
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
	currentSession.App = req.App
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
