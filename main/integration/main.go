package integration

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Log = log.New(os.Stdout, "node-integration ", log.Flags())
var FilePath = ""

// * Credentials for this node
var NodePrivateKey *rsa.PrivateKey // Private key for current node
var NodePublicKey *rsa.PublicKey   // Public key for current node

type NodeData struct {
	NodeToken string
	NodeId    uint
	AppId     uint
}

const IdentifierChatNode = "chat"
const IdentifierSpaceNode = "space"

var Nodes map[string]NodeData = make(map[string]NodeData)

func Setup(identifier string) bool {

	// Setup environment
	err := godotenv.Load()
	if err != nil {
		Log.Println("Error loading .env file")
		return false
	}

	if os.Getenv("PROTOCOL") == "" {
		Log.Println("Please set PROTOCOL in the .env file to 'https://' if you are using HTTPS.")
		os.Setenv("PROTOCOL", "http://")
	} else {
		Protocol = os.Getenv("PROTOCOL")
	}

	scanner := bufio.NewScanner(os.Stdin)

	// Get the public and private key of this node
	NodePublicKey, err = UnpackageRSAPublicKey(os.Getenv("TC_PUBLIC_KEY"))
	if err != nil {
		panic("Couldn't unpackage public key. Required for v1 API. Please set TC_PUBLIC_KEY in your environment variables or .env file.")
	}

	NodePrivateKey, err = UnpackageRSAPrivateKey(os.Getenv("TC_PRIVATE_KEY"))
	if err != nil {
		panic("Couldn't unpackage private key. Required for v1 API. Please set TC_PRIVATE_KEY in your environment variables or .env file.")
	}

	// Ask for testing if required
	testing := os.Getenv("TESTING")
	if testing == "true" {
		Testing = true
	} else {

		Log.Println("Do you want to run this node in testing mode? (y/n)")

		scanner.Scan()
		input := scanner.Text()
		Testing = input == "y"
	}

	if Testing {
		TestingToken = os.Getenv("TESTING_SECRET")

		if TestingToken == "" {
			Log.Println("Testing mode enabled. Please set TESTING_SECRET in the .env file if you wish to use the feature.")
		} else {
			Log.Println("Testing mode enabled. Testing token: " + TestingToken)
		}
	}

	// Check if already setup
	_, ok := Nodes[identifier]
	if ok {
		return true
	}

	input := ""
	if os.Getenv("DEFAULT_FILE") != "" {
		Log.Println("Using default file from .env: " + os.Getenv("DEFAULT_FILE"))
		input = os.Getenv("DEFAULT_FILE")
	} else {
		Log.Println("Please provide the file name of the data file (e.g. data (.node will be appended automatically))")
		scanner.Scan()
		input = scanner.Text()
	}

	Log.Println("Continuing in normal mode..")

	if os.Getenv("NODE_ENV") == "" {
		Log.Println("Please set NODE_ENV to the path for data files in the .env file.")
		return false
	}
	FilePath = os.Getenv("NODE_ENV")

	if readData(FilePath+"/"+input+".node", identifier) {

		Log.Println("Grabbing server public key..")
		err = grabServerPublicKey()
		if err != nil {
			return false
		}

		Log.Println("Ready to start.")
		return true
	}

	var creationToken, nodeDomain string
	Log.Println("No data file found. Please enter the following information:")

	Log.Println("1. Base Path (e.g. http://localhost:3000)")
	scanner.Scan()
	BasePath = scanner.Text()

	Log.Println("Grabbing server public key..")
	err = grabServerPublicKey()
	if err != nil {
		return false
	}

	Log.Println("2. Creation Token (Received from a creation request in the admin panel)")
	scanner.Scan()
	creationToken = scanner.Text()

	Log.Println("Getting clusters..")
	res, err := PostRequest("/node/manage/clusters", map[string]interface{}{
		"token": creationToken,
	})

	if err != nil {
		Log.Println("Your creation token is invalid.")
		return false
	}

	clusterId := setupClusters(res, scanner)

	Log.Println("4. App id (e.g. 1)")
	scanner.Scan()
	appId, err := strconv.Atoi(scanner.Text())

	if err != nil {
		Log.Println("Please enter a valid number.")
		return false
	}

	Log.Println("5. The domain of this node (e.g. node-1.example.com)")
	scanner.Scan()
	nodeDomain = scanner.Text()

	Log.Println("6. The performance level (relative to all other nodes) of this node (e.g. 0.75)")
	scanner.Scan()
	performanceLevel, err := strconv.ParseFloat(scanner.Text(), 64)

	if err != nil {
		Log.Println("Please enter a valid number.")
		return false
	}

	Log.Println("Creating node..")

	res, err = PostRequest("/node/manage/new", map[string]interface{}{
		"token":             creationToken,
		"domain":            nodeDomain,
		"performance_level": performanceLevel,
		"app":               appId,
		"cluster":           clusterId,
	})

	if err != nil {
		Log.Println("Error while creating node.")
		return false
	}

	if !res["success"].(bool) {
		Log.Println("Error while creating node. Please check your input.")
		return false
	}

	Log.Println("Node created successfully.")

	// Setup
	Nodes[identifier] = NodeData{
		NodeToken: res["token"].(string),
		NodeId:    uint(res["id"].(float64)),
		AppId:     uint(appId),
	}

	// Save data to file
	f, err := os.Create(FilePath + "/" + input + ".node")
	if err != nil {
		Log.Println("Error while creating data file. err:", err.Error())
		return false
	}
	defer f.Close()

	// Write data to file
	f.WriteString(BasePath + "\n")
	f.WriteString(Nodes[identifier].NodeToken + "\n")
	f.WriteString(fmt.Sprintf("%d", Nodes[identifier].NodeId) + "\n")

	Log.Println("Data saved to file.")

	return true
}

func readData(path string, identifier string) bool {
	Log.Println(path)

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	BasePath = scanner.Text()

	scanner.Scan()
	nodeToken := scanner.Text()

	scanner.Scan()
	nodeId, err := strconv.Atoi(scanner.Text())

	Nodes[identifier] = NodeData{
		NodeToken: nodeToken,
		NodeId:    uint(nodeId),
	}

	if err != nil {
		Log.Println("Error while reading data file.")
		return false
	}

	return true
}
