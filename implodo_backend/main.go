package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"implodo_backend/auth"
	"implodo_backend/config"
	"implodo_backend/store"
)

func main() {
	cfg := config.Load()
	sessions := store.NewSessionStore()
	authHandler := auth.NewHandler(cfg, sessions)

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
	}

	log.Printf("Server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
