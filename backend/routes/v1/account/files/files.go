package files

import (
	"context"
	"log"
	"os"
	"time"

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
)

var bucketName string
var client *s3.Client
var uploader *manager.Uploader
var disabled = false

// Configuration
const maxUploadSize = 10_000_000       // 10 MB
const maxFavoriteStorage = 500_000_000 // 500 MB
const maxTotalStorage = 1_000_000_000  // 1 GB

func Unencrypted(router fiber.Router) {

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

			log.Println(err.Error())

			// Return error message
			return c.SendStatus(401)
		},
	}))

	router.Post("/upload", uploadFile)
}

func Authorized(router fiber.Router) {
	url := os.Getenv("R2_URL")
	bucketName = os.Getenv("R2_BUCKET")
	accessKeyId := os.Getenv("R2_CLIENT_ID")
	accessKeySecret := os.Getenv("R2_CLIENT_SECRET")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: url,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		disabled = true
		log.Println("Failed to connect to R2. File integration disabled.")
		log.Fatal(err)
	}

	// Setup uploader
	client = s3.NewFromConfig(cfg)

	log.Println("Checking R2 connection..")
	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Println("R2 NOT WORKING")
		panic(err)
	}
	log.Println("Successfully connected to R2.")

	// Setup file routes
	router.Post("/delete", deleteFile)
	router.Post("/list", listFiles)
	router.Post("/favorite", favoriteFile)
	router.Post("/unfavorite", unfavoriteFile)
	router.Post("/info", fileInfo)
}

func CountTotalStorage(accId string) (int64, error) {

	// Get total storage (coalesce is important cause otherwise we get null)
	var totalStorage int64
	unix := time.Now().Add(-time.Hour * 24 * 30).UnixMilli()
	rq := database.DBConn.Model(&account.CloudFile{}).Where("account = ? AND (created_at > ? OR favorite = ?)", accId, unix, true).Select("coalesce(sum(size), 0)").Scan(&totalStorage)
	if rq.Error != nil {
		if rq.RowsAffected > 0 {
			return 0, nil
		}
		return 0, rq.Error
	}

	return totalStorage, nil
}

func CountFavoriteStorage(accId string) (int64, error) {

	// Get favorite storage (coalesce is important cause otherwise we get null)
	var favoriteStorage int64
	if err := database.DBConn.Model(&account.CloudFile{}).Where("account = ? AND favorite = ?", accId, true).Select("coalesce(sum(size), 0)").Scan(&favoriteStorage).Error; err != nil {
		return 0, err
	}

	return favoriteStorage, nil
}
