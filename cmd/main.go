package main

import (
	"context"
	"fmt"
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
	handler := NewHandler(userService)
	http.HandleFunc("/users", handler.HelloWorld)
	http.HandleFunc("/hello-world", handler.CreateUser)
	http.HandleFunc("/users/{id}", handler.GetUser)
	http.HandleFunc("/users/{id}/update", handler.UpdateUser)
	fmt.Println("server running at 8080")
	http.ListenAndServe(":8080", nil)
}
