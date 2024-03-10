package manage

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

type newRequest struct {
	Token           string  `json:"token"`
	Cluster         uint    `json:"cluster"` // Cluster ID
	App             uint    `json:"app"`     // App ID
	Domain          string  `json:"domain"`
	PeformanceLevel float32 `json:"performance_level"`
}

func newNode(c *fiber.Ctx) error {

	// Parse body to add request
	var req newRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if token is valid
	var ct node.NodeCreation
	if err := database.DBConn.Where("token = ?", req.Token).Take(&ct).Error; err != nil {
		return util.FailedRequest(c, "invalid", nil)
	}

	if req.Cluster == 0 || req.Domain == "" {
		return util.FailedRequest(c, "invalid", nil)
	}

	if len(req.Domain) < 3 {
		return util.FailedRequest(c, "invalid.domain", nil)
	}

	var cluster node.Cluster
	if err := database.DBConn.Where("id = ?", req.Cluster).Take(&cluster).Error; err != nil {
		return util.FailedRequest(c, "invalid", nil)
	}

	var app app.App
	if err := database.DBConn.Take(&app, req.App).Error; err != nil {
		return util.FailedRequest(c, "invalid", nil)
	}

	// Create node
	var created node.Node = node.Node{
		AppID:           req.App,
		ClusterID:       req.Cluster,
		Token:           auth.GenerateToken(300),
		Domain:          req.Domain,
		Load:            0,
		PeformanceLevel: req.PeformanceLevel,
		Status:          1,
	}

	if err := database.DBConn.Create(&created).Error; err != nil {
		return util.FailedRequest(c, "invalid.domain", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   created.Token,
		"cluster": cluster.Country,
		"app":     app.Name,
		"id":      created.ID,
	})
}
