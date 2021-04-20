package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var cdnApp *firebase.App
var cdnAuth *auth.Client
var cdnFirestore *firestore.Client
var cdnS3Config *aws.Config
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
		CdnEndpoint:   os.Getenv("CDN_ENDPOINT"),
		Authorization: os.Getenv("AUTHORIZATION"),
		Production:    os.Getenv("PRODUCTION") != "false",
	}

	if cdnConfig.Authorization == "" {
		log.Fatal("No Authorization token provided, closing...")
	}

	mode := "DEVELOPMENT"
	if cdnConfig.Production {
		mode = "PRODUCTION"
	}
	log.Printf("Starting in %v mode", mode)

	cdnS3Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials(cdnConfig.SpacesConfig.SpacesAccessKey, cdnConfig.SpacesConfig.SpacesSecretKey, ""),
		Endpoint:    aws.String(cdnConfig.SpacesConfig.SpacesEndpoint),
		Region:      aws.String(cdnConfig.SpacesConfig.SpacesRegion),
	}

	setUpFirebase()
	setUpFirebaseFirestore()
	setUpFirebaseAuth()
}

func setUpRoutes() {
	server := fiber.New()

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, User-Agent",
	}))

	if cdnConfig.Production {
		server.Static("/", "./public")
	}

	server.Get("/:file", getFileRoute)
	server.Get("/oembed/:file", getOGEmbedRoute)

	api := server.Group("/api")

	// api.Get("/user", authorize, getUserRoute) // auth
	// api.Post("/user", createUserRoute) // auth
	// api.Get("/ws", authorize, getWebSocket) // auth
	api.Post("/verify", verifyAuthRoute) // auth

	// files
	api.Post("/upload", authorize, uploadFileRoute)      // auth
	api.Get("/files", authorize, getFilesRoute)          // auth
	api.Delete("/files/:id", authorize, deleteFileRoute) // auth

	// folders
	api.Get("/folders", authorize, getFoldersRoute)    // auth
	api.Post("/folders", authorize, createFolderRoute) // auth
	api.Get("/folders/:id", getFolderRoute)
	api.Patch("/folders/:id", authorize, updateFolderRoute)  // auth
	api.Delete("/folders/:id", authorize, deleteFolderRoute) // auth

	log.Fatal(server.Listen(":3000"))
}

func setUpFirebase() {
	options := option.WithCredentialsFile("service-account.json")
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
