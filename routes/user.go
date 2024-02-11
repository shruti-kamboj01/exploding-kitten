package routes

import (
	"context"
	"encoding/json"
	"net/http"
    "fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	// Add more fields as needed
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	// Decode the JSON request body into a struct
	var requestBody struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the username
	if requestBody.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if the username already exists
	exists, err := CheckUsernameExists(requestBody.Username, client)
	if err != nil {
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Create a new user
	newUser := User{
		
		Username: requestBody.Username,
		// Add more fields as needed
	}

	// Save the user to Redis
	err = SaveUser(newUser, client)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully: %s", newUser.Username)
}

// CheckUsernameExists checks if a username already exists in the user model stored in Redis
func CheckUsernameExists(username string, client *redis.Client) (bool, error) {
	// Retrieve all user IDs from Redis
	userIDs, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return false, err
	}

	// Iterate over each user ID
	for _, userID := range userIDs {
		// Retrieve the user data from Redis
		userJSON, err := client.Get(ctx, userID).Bytes()
		if err != nil {
			return false, err
		}

		// Deserialize the JSON data into a User struct
		var user User
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			return false, err
		}

		// Check if the username matches
		if user.Username == username {
			return true, nil
		}
	}

	// Username does not exist
	return false, nil
}

// SaveUser saves a user to Redis
func SaveUser(user User, client *redis.Client) error {
	// Serialize the user struct to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Save the user data to Redis
	err = client.Set(ctx, user.ID, userJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}