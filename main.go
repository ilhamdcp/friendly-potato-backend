package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	. "github.com/ilhamdcp/friendly-potato/internal/data/firebase"
	. "github.com/ilhamdcp/friendly-potato/internal/delivery/http"
	"github.com/ilhamdcp/friendly-potato/internal/service"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userRepository := NewUserRepository(client)
	userService := service.NewUserServiceImpl(userRepository, os.Getenv("HASH_SECRET"))
	contactRepository := NewContactRepository(client)
	contactService := service.NewContactServiceImpl(contactRepository)
	handler := NewHandler(userService, contactService)
	http.HandleFunc("/users", handler.CreateUser)
	http.HandleFunc("/hello-world", handler.HelloWorld)
	http.HandleFunc("/users/{id}", handler.GetUser)
	http.HandleFunc("/users/{id}/update", handler.UpdateUser)
	http.HandleFunc("/users/auth", handler.AuthenticateUser)
	http.HandleFunc("/users/sign-in", handler.SignInUser)
	http.HandleFunc("/users/sign-out", handler.SignOutUser)
	http.HandleFunc("/contacts/list", handler.GetContacts)
	http.HandleFunc("/contacts/add", handler.AddContact)
	http.HandleFunc("/contacts/remove", handler.RemoveContact)
	fmt.Println("server running at 8080")
	http.ListenAndServe(":8080", nil)
}
