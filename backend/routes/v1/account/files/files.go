package files

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var fileRepoType string = "local"
var s3Client *s3.Client
var bucketName string
var uploader *manager.Uploader
var disabled = false

// Constants
const repoTypeR2 = "r2"
const repoTypeLocal = "local"

// Configuration
var maxUploadSize int64 = 10      // 10 MB
var maxTotalStorage int64 = 1_000 // 1 GB
var saveLocation = ""
var urlPath = ""

func Unencrypted(router fiber.Router) {

	if os.Getenv("FILE_REPO") == "" {
		util.Log.Println("If you want to enable file uploading, please specify a path for those files using the FILE_REPO env variable!")
		disabled = true
	} else {
		disabled = false
		saveLocation = os.Getenv("FILE_REPO")
		if !strings.HasSuffix(saveLocation, "/") {
			saveLocation = saveLocation + "/"
		}
	}

	if os.Getenv("BASE_PATH") == "" {
		util.Log.Println("If you want to enable file uploading, please specify the domain of the server in the BASE_PATH env variable (without https:// or http:// (that's specified in the PROTOCOL env variable, https:// by default), you can specify a port if needed)")
		disabled = true
	} else {
		urlPath = os.Getenv("PROTOCOL") + os.Getenv("BASE_PATH")
	}

	if !disabled {
		maxUploadSize = GetIntEnv("MAX_UPLOAD_SIZE", maxUploadSize) * 1_000_000
		maxTotalStorage = GetIntEnv("MAX_TOTAL_STORAGE", maxTotalStorage) * 1_000_000
	}

	// Autorized by using a normal JWT token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS512,
			Key:    []byte(util.JWT_SECRET),
		},

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			if util.IsExpired(c) {
				return util.InvalidRequest(c)
			}

			return c.Next()
		},

		// Error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			util.Log.Println(err.Error())

			// Return error message
			return c.SendStatus(401)
		},
	}))

	if os.Getenv("FILE_REPO_TYPE") != "" {
		fileRepoType = os.Getenv("FILE_REPO_TYPE")

		// Connect to r2 if the file repo is r2
		if fileRepoType == "r2" {
			connectToR2()
		}
	}

	router.Post("/upload", uploadFile)
}

func UnencryptedUnauthorized(router fiber.Router) {
	router.Post("/download/:id", downloadFile)
}

func GetIntEnv(key string, standard int64) int64 {
	envValue := os.Getenv(key)
	if envValue == "" {
		return standard
	} else {
		envInt, err := strconv.Atoi(envValue)
		if err != nil {
			util.Log.Println("ERROR: Couldn't read", key, ". Please set it if you want to modify the option. Default value:", standard)
			return standard
		}

		return int64(envInt)
	}
}

func Authorized(router fiber.Router) {

	// Setup file routes
	router.Post("/delete", deleteFile)
	router.Post("/list", listFiles)
	router.Post("/change_tag", changeFileTag)
	router.Post("/info", fileInfo)
	router.Post("/storage", getStorageUsage)
}

func CountTotalStorage(accId uuid.UUID) (int64, error) {

	// Get total storage (coalesce is important cause otherwise we get null)
	var totalStorage int64
	rq := database.DBConn.Model(&account.CloudFile{}).Where("account = ?", accId).Select("coalesce(sum(size), 0)").Scan(&totalStorage)
	if rq.Error != nil {
		if rq.RowsAffected > 0 {
			return 0, nil
		}
		return 0, rq.Error
	}

	return totalStorage, nil
}

func connectToR2() {
	bucketName = os.Getenv("FILE_REPO_BUCKET")
	var url = os.Getenv("FILE_REPO")
	var accessKeyId = os.Getenv("FILE_REPO_KEY_ID")
	var accessKeySecret = os.Getenv("FILE_REPO_KEY")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to R2 and make a new uploader
	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(url)
	})
	uploader = manager.NewUploader(s3Client)

	// Make sure the API works
	_, err = s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: aws.Int32(10),
	})
	if err != nil {
		util.Log.Fatal(err)
	}

	util.Log.Println("Successfully connected to Cloudflare R2.")
}
