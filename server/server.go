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

	server.Get("/:id.:ext", getFileRoute)
	server.Get("/oembed/:id.:ext", getOGEmbedRoute)

	server.Static("/", "../client/public")

	api := server.Group("/api")

	// api.Get("/user", authorize, getUserRoute) // auth
	// api.Post("/user", createUserRoute) // auth

	api.Post("/upload", dummyMiddleware, uploadFileRoute) // auth
	// api.Get("/ws", authorize, getWebSocket) // auth

	// single files
	api.Get("/file/:id", getFileInfoRoute)
	// no reason to update files just yet
	// api.Patch("/file/:id/info", authorize, updateFileInfoRoute) // auth
	api.Delete("/file/:id", dummyMiddleware, deleteFileRoute) // auth

	// folder of files
	api.Post("/folder", dummyMiddleware, createFolderRoute) // auth
	api.Get("/folder/:id", getFolderRoute)
	api.Patch("/folder/:id", dummyMiddleware, updateFolderRoute)  // auth
	api.Delete("/folder/:id", dummyMiddleware, deleteFolderRoute) // auth

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
