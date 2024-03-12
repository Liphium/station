package integration

import (
	"bufio"
	"strconv"
)

func setupClusters(res map[string]interface{}, scanner *bufio.Scanner) int64 {

	var clusterId int64
	for clusterId == 0 {

		Log.Println("3. Please select a cluster:")

		clusters := res["clusters"].([]interface{})
		for i, cluster := range clusters {

			cluster := cluster.(map[string]interface{})
			Log.Println(i, cluster["name"], " (", cluster["id"], ") - ", cluster["country"])
		}

		scanner.Scan()
		input, err := strconv.Atoi(scanner.Text())

		if err != nil {
			Log.Println("Please enter a valid number.")
			continue
		}

		if input < 0 || input >= len(clusters) {
			Log.Println("Please enter a valid number.")
			continue
		}

		clusterId = int64(clusters[input].(map[string]interface{})["id"].(float64))
	}

	return clusterId
}
