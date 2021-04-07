package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"google.golang.org/api/option"
)

var cdnApp *firebase.App
var cdnAuth *auth.Client
var cdnFirestore *firestore.Client
var cdnSpaces *minio.Client
var cdnConfig *Config

func main() {
	setUpRoutes()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
		log.Fatal(err)
	}

	cdnConfig = &Config{
		SpacesConfig: SpacesConfig{
			SpacesAccessKey: os.Getenv("SPACES_ACCESS_KEY"),
			SpacesSecretKey: os.Getenv("SPACES_SECRET_KEY"),
			SpacesEndpoint:  os.Getenv("SPACES_ENDPOINT"),
			SpacesUrl:       os.Getenv("SPACES_URL"),
			SpacesCdn:       os.Getenv("SPACES_CDN_URL"),
			SpacesName:      os.Getenv("SPACES_NAME"),
			SpacesRegion:    os.Getenv("SPACES_REGION"),
		},
		CdnEndpoint: os.Getenv("CDN_ENDPOINT"),
		AccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	setUpSpaces()
	setUpFirebase()
	setUpFirebaseFirestore()
	setUpFirebaseAuth()
}

func setUpRoutes() {
	server := fiber.New()

	server.Static("/", "../client/public")

	api := server.Group("/api")

	// api.Get("/user", authorize, getUserRoute)
	// api.Post("/user", createUserRoute)

	api.Post("/upload", uploadFileRoute)
	// api.Get("/ws", authorize, getWebSocket)

	// single files
	api.Get("/file/:id", getFileRoute)
	api.Get("/file/:id/info", getFileInfoRoute)
	// api.Patch("/file/:id/info", authorize, updateFileInfoRoute)
	// api.Delete("/file/:id", authorize, deleteFileRoute)

	// server.Get("/oembed/:id", getOGEmbedRoute)

	// folder of files
	// api.Get("/files", getFilesRoute)
	// api.Patch("/files", authorize, updateFilesRoute)
	// api.Delete("/files", authorize, deleteFilesRoute)

	log.Fatal(server.Listen(":3000"))
}

func setUpFirebase() {
	options := option.WithCredentialsFile("./service-account.json")
	ctx := context.Background()

	fbApp, err := firebase.NewApp(ctx, nil, options)
	cdnApp = fbApp

	if err != nil {
		log.Printf("Could not connect to Firebase")
		log.Fatal(err)
		return
	}

	log.Printf("Connected to Firebase")
}

func setUpFirebaseFirestore() {
	ctx := context.Background()
	fbStore, err := cdnApp.Firestore(ctx)
	cdnFirestore = fbStore

	if err != nil {
		log.Printf("Could not connect to Firebase Firestore")
		log.Fatal(err)
		return
	}

	log.Printf("Connected to Firebase Firestore")
}

func setUpFirebaseAuth() {
	ctx := context.Background()
	fbAuth, err := cdnApp.Auth(ctx)
	cdnAuth = fbAuth

	if err != nil {
		log.Printf("Could not connect to Firebase Auth")
		log.Fatal(err)
		return
	}

	log.Printf("Connected to Firebase Auth")
}

func setUpSpaces() {
	sps, err := minio.New(
		cdnConfig.SpacesConfig.SpacesEndpoint,
		cdnConfig.SpacesConfig.SpacesAccessKey,
		cdnConfig.SpacesConfig.SpacesSecretKey,
		false,
	)

	cdnSpaces = sps

	if err != nil {
		log.Printf("Could not connect to DigitalOcean Spaces")
		log.Fatal(err)
		return
	}

	log.Printf("Connected to DigitalOcean Spaces")

	bucketName := cdnConfig.SpacesConfig.SpacesName

	err = cdnSpaces.MakeBucket(bucketName, cdnConfig.SpacesConfig.SpacesRegion)
	if err != nil {
		exists, bucketExistErr := cdnSpaces.BucketExists(bucketName)

		if bucketExistErr == nil && exists {
			log.Printf("We already own bucket \"%s\"", bucketName)
		} else {
			log.Printf("Unable to create bucket \"%s\"", bucketName)
			log.Fatalln(err)
			return
		}
	} else {
		log.Printf("Successfully created bucket \"%s\"", bucketName)
	}
}
