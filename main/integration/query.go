package integration

import (
	"os"
	"strconv"

	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

func GetCurrent(identifier string) (id int64, token string, app uint, domain string, err error) {

	res, err := PostRequestBackend("/node/this", fiber.Map{
		"node":  Nodes[identifier].NodeId,
		"token": Nodes[identifier].NodeToken,
	})

	if err != nil {
		return 0, "", 0, "", err
	}

	if !res["success"].(bool) {
		Log.Println("ERROR: This node may not be registered..")
		os.Exit(1)
	}

	JwtSecret = res["jwt_secret"].(string)
	n := res["node"].(map[string]interface{})

	return int64(n["id"].(float64)), n["token"].(string), uint(n["app"].(float64)), n["domain"].(string), nil
}

func SetOnline(identifier string) (map[string]interface{}, error) {

	res, err := PostRequestBackend("/node/status/online", fiber.Map{
		"id":    Nodes[identifier].NodeId,
		"token": Nodes[identifier].NodeToken,
	})

	if err != nil {
		return map[string]interface{}{}, err
	}

	if !res["success"].(bool) {
		Log.Println("This node may not be registered..")
		os.Exit(1)
	}

	return res, nil
}

func ReportOffline(node pipes.Node) {

	Log.Println("Outgoing event stream to node", node.ID, "disconnected.")

	// Convert node id
	nodeID, _ := strconv.Atoi(node.ID)

	_, err := PostRequestBackend("/node/status/offline", fiber.Map{
		"node":  nodeID,
		"token": node.Token,
	})

	if err != nil {
		Log.Println("Failed to report offline status. Is the backend online?")
	}
}
