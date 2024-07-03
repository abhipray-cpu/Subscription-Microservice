// Package data defines the data layer of the application, including models and database operations.
package data

import (
	"context" // Used for managing the lifetime of database requests.
	"log"
	"regexp"
	"strings"
	"time" // Used for handling time-related data.

	"github.com/jackc/pgx/v4" // PostgreSQL driver for Go.
)

// User represents a user entity in the system with various attributes.
type User struct {
	ID          int64     `json:"id"`          // Unique identifier for the user.
	UserName    string    `json:"userName"`    // Username of the user.
	GithubName  string    `json:"githubName"`  // GitHub username of the user.
	GithubId    string    `json:"githubId"`    // GitHub ID of the user.
	FirstName   string    `json:"firstName"`   // First name of the user.
	LastName    string    `json:"lastName"`    // Last name of the user.
	AvatarUrl   string    `json:"avatarUrl"`   // URL of the user's avatar.
	AccessToken string    `json:"accessToken"` // Access token for authentication.
	Bio         string    `json:"bio"`         // Biography of the user.
	Email       string    `json:"email"`       // Email address of the user.
	ExpiresAt   time.Time `json:"expiresAt"`   // Expiration time of the user's session or token.
	Password    string    `json:"password"`    // Password of the user.
}

// connection holds a global database connection, shared across instances of Models.
var connection *pgx.Conn

// Models wraps all the models in the application for easy access.
type Models struct {
	User User // User model instance.
}

// NewModels initializes a new instance of Models with a database connection.
func NewModels(conn *pgx.Conn) Models {
	connection = conn       // Set the global connection.
	ensureTableExists(conn) // Ensure the table exists in the database.
	return Models{
		User: User{}, // Initialize the User model.
	}
}

// This functions ensures that a table exists on startup
// If the table does not exist, it creates the table
func ensureTableExists(conn *pgx.Conn) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        user_name VARCHAR(255) NOT NULL,
        github_name VARCHAR(255) UNIQUE NOT NULL,
        github_id VARCHAR(255) UNIQUE,
        first_name VARCHAR(255),
        last_name VARCHAR(255),
        avatar_url TEXT,
        access_token TEXT,
        bio TEXT,
        email VARCHAR(255) NOT NULL UNIQUE,
        expires_at TIMESTAMP NOT NULL,
		password VARCHAR(255) NOT NULL
    );`

	if _, err := conn.Exec(context.Background(), query); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// validateUser checks the user's fields for correctness.
// Returns true if validation passes along with an empty string,
// or false and an error message if validation fails.
func (u *User) ValidateUser() (bool, string) {
	// Check for empty fields
	if strings.TrimSpace(u.UserName) == "" {
		return false, "user_name cannot be empty"
	}
	if strings.TrimSpace(u.GithubName) == "" {
		return false, "github_name cannot be empty"
	}

	// Email validation
	if !isValidEmail(u.Email) {
		return false, "email format is invalid"
	}

	// All validations passed
	return true, ""
}

// isValidEmail checks if the email provided passes the regex validation.
// Returns true if the email is valid, false otherwise.
func isValidEmail(email string) bool {
	// Simple regex for checking email format
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regex.MatchString(email)
}

// InsertUser inserts a new user into the database.
// This method is useful for registering new users in the system.
// Parameters:
// - user: The User struct containing the user's information.
// Returns:
// - An error if the query execution fails.
func (u *User) InsertUser(user User) error {
	// SQL query to insert a new user, returning the generated ID.
	query := `INSERT INTO users (user_name, github_name, github_id, first_name, last_name, avatar_url, bio, email, expires_at,password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	// Execute the query and scan the returned ID into the User struct.
	err := connection.QueryRow(context.Background(), query, user.UserName, user.GithubName, user.GithubId, user.FirstName, user.LastName, user.AvatarUrl, user.Bio, user.Email, user.ExpiresAt, user.Password).Scan(&u.ID)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}

// GetUser retrieves a user by their ID from the database.
// This method is useful for fetching user details based on their unique identifier.
// Parameters:
// - id: The ID of the user to retrieve.
// Returns:
// - nil if the user is successfully found and the User struct is populated.
// - An error if the query execution or scan fails.
func (u *User) GetUser(id int64) error {
	// SQL query to select a user by ID.
	query := `SELECT id, user_name, github_name, github_id, first_name, last_name, avatar_url, bio, email, expires_at FROM users WHERE id=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.UserName, &u.GithubName, &u.GithubId, &u.FirstName, &u.LastName, &u.AvatarUrl, &u.Bio, &u.Email, &u.ExpiresAt)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}

// UpdateUser updates an existing user's information in the database.
// This method is useful for updating user details, such as their name, email, or avatar.
// Parameters:
// - id: The ID of the user to update.
// - updatedUser: The updated User struct containing the new information.
// Returns:
// - An error if the query execution fails.
func (u *User) UpdateUser(id int64, updatedUser User) error {
	// SQL query to update a user's information by ID.
	query := `UPDATE users SET user_name=$1, github_name=$2, github_id=$3, first_name=$4, last_name=$5, avatar_url=$6, access_token=$7, bio=$8, email=$9, expires_at=$10 WHERE id=$11`
	// Execute the query without returning any result.
	_, err := connection.Exec(context.Background(), query, updatedUser.UserName, updatedUser.GithubName, updatedUser.GithubId, updatedUser.FirstName, updatedUser.LastName, updatedUser.AvatarUrl, updatedUser.AccessToken, updatedUser.Bio, updatedUser.Email, updatedUser.ExpiresAt, id)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}

// DeleteUser removes a user from the database by their ID.
// This method is useful for removing a user from the system.
// Parameters:
// - id: The ID of the user to delete.
// Returns:
// - An error if the query execution fails.s
func (u *User) DeleteUser(id int64) error {
	// SQL query to delete a user by ID.
	query := `DELETE FROM users WHERE id=$1`
	// Execute the query without returning any result.
	_, err := connection.Exec(context.Background(), query, id)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}

// GetByGitId retrieves a user by their GitHub ID from the database.
// This method is useful for integrating GitHub authentication,
// allowing the application to fetch user details based on GitHub account information.
//
// Parameters:
// - githubId: The GitHub ID of the user to retrieve.
//
// Returns:
// - nil if the user is successfully found and the User struct is populated.
// - An error if the query execution or scan fails.
func (u *User) GetByGitId(githubId string) error {
	// SQL query to select a user by GitHub ID.
	query := `SELECT id, user_name, github_name, github_id, first_name, last_name, avatar_url, access_token, bio, email, expires_at FROM users WHERE github_id=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, githubId).Scan(&u.ID, &u.UserName, &u.GithubName, &u.GithubId, &u.FirstName, &u.LastName, &u.AvatarUrl, &u.AccessToken, &u.Bio, &u.Email, &u.ExpiresAt)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}

// GetByEmail retrieves a user by their email address from the database.
// This method is useful for authenticating users based on their email address,
// allowing the application to fetch user details based on their email.
// Parameters:
// - email: The email address of the user to retrieve.
// Returns:
// - nil if the user is successfully found and the User struct is populated.
// - An error if the query execution or scan fails.
func (u *User) GetByEmail(email string) error {
	// SQL query to select a user by GitHub ID.
	query := `SELECT id, user_name,password,email FROM users WHERE email=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, email).Scan(&u.ID, &u.UserName, &u.Password, &u.Email)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}
