package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/shruti-kamboj01/exploding_kitten/routes"
	"github.com/shruti-kamboj01/exploding_kitten/routes/user"
)

func main() {
    ctx := context.Background()
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis server address
        Password: "",                // no password set
        DB:       0,                 // use default DB
    })

    // Ping the Redis server to check the connection
    pong, err := client.Ping(ctx).Result()
    fmt.Println(pong, err)
}

func SetupRoutes() *mux.Router {
    router := mux.NewRouter()
    
    // Define your routes here
   
    router.HandleFunc("/", routes.CreateUserHandler).Methods("POST")
    // Add more routes as needed
    
    return router
}