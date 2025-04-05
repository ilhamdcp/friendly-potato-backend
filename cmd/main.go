package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	. "github.com/ilhamdcp/friendly-potato/internal/data/firebase"
	. "github.com/ilhamdcp/friendly-potato/internal/delivery/http"
	"github.com/ilhamdcp/friendly-potato/internal/service"
)

func generateClient(ctx context.Context) *firestore.Client {
	projectID := "friendly-potato"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func main() {
	client := generateClient(context.Background())
	defer client.Close()
	userRepository := NewUserRepository(client)
	userService := service.NewUserServiceImpl(userRepository)
	userHandler := NewUserHandler(userService)
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/users/{id}", userHandler.GetUser)
	http.HandleFunc("/users/{id}/update", userHandler.UpdateUser)
	http.HandleFunc("/users/signin", userHandler.SignInWithGoogle)
	http.ListenAndServe(":8080", nil)
}
