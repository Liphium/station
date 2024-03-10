package integration

import (
	"os"
	"strconv"

	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

var ClusterID uint = 0

func GetCurrent() (id int64, token string, app uint, domain string) {

	res, err := PostRequest("/node/this", fiber.Map{
		"node":  NODE_ID,
		"token": NODE_TOKEN,
	})

	if err != nil {
		Log.Println("Backend is currently offline!")
		os.Exit(1)
	}

	if !res["success"].(bool) {
		Log.Println("This node may not be registered..")
		os.Exit(1)
	}

	JwtSecret = res["jwt_secret"].(string)
	n := res["node"].(map[string]interface{})
	ClusterID = uint(res["cluster"].(float64))

	return int64(n["id"].(float64)), n["token"].(string), uint(n["app"].(float64)), n["domain"].(string)
}

func SetOnline() map[string]interface{} {

	res, err := PostRequest("/node/status/online", fiber.Map{
		"id":    NODE_ID,
		"token": NODE_TOKEN,
	})

	if err != nil {
		Log.Println("Backend is currently offline!")
		os.Exit(1)
	}

	if !res["success"].(bool) {
		Log.Println("This node may not be registered..")
		os.Exit(1)
	}

	return res
}

func ReportOffline(node pipes.Node) {

	Log.Println("Outgoing event stream to node", node.ID, "disconnected.")

	// Convert node id
	nodeID, _ := strconv.Atoi(node.ID)

	_, err := PostRequest("/node/status/offline", fiber.Map{
		"node":  nodeID,
		"token": node.Token,
	})

	if err != nil {
		Log.Println("Failed to report offline status. Is the backend online?")
	}
}
