package main

import (
	"context"
	"log"

	"implodo_backend/auth"
	"implodo_backend/config"
	"implodo_backend/imagegen"
	"implodo_backend/imagestorage"
	"implodo_backend/store"
	"implodo_backend/upload"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("local.env"); err != nil {
		log.Println("No local.env file found, using environment variables")
	}
	cfg := config.Load()

	ctx := context.Background()

	genClient, err := imagegen.NewGeminiClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	storageClient, err := imagestorage.NewGCSClient(ctx, cfg.BucketName)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}

	sessions := store.NewSessionStore()
	authHandler := auth.NewHandler(cfg, sessions)
	uploadHandler := upload.NewHandler(genClient, storageClient)

	r := gin.Default()

	// --- Auth flow ---
	// Step 1: App calls this to get the Google auth URL, passing state + code_challenge.
	r.POST("/auth/start", authHandler.Start)

	// Step 2: Google redirects here after login. Backend verifies the user, issues JWT,
	// then redirects the browser to the app's deep link (e.g. implodo://auth/complete?state=...).
	r.GET("/auth/google/callback", authHandler.GoogleCallback)

	// Step 3: App calls this with state + code_verifier to claim the JWT.
	// Returns 202 if the browser flow isn't done yet — app should retry.
	r.POST("/auth/token", authHandler.ExchangeToken)

	// --- Protected routes ---
	api := r.Group("/api")
	api.Use(auth.JWTMiddleware(cfg))
	{
		api.GET("/me", authHandler.Me)
		api.GET("/profile", authHandler.Profile)
		api.POST("/upload", uploadHandler.UploadImage)
		api.GET("/images", uploadHandler.ListImages)
	}

	log.Printf("Server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
