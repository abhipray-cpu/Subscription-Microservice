// Package data defines the data layer of the application, including models and database operations.
package data

import (
	"context" // Used for managing the lifetime of database requests.
	"fmt"
	"log"
	"regexp"
	"strings"
	"time" // Used for handling time-related data.

	"github.com/jackc/pgx/v4" // PostgreSQL driver for Go.
)

// User represents a user entity in the system with various attributes.
type User struct {
	ID                 int64     `json:"id"`                 // Unique identifier for the user.
	UserName           string    `json:"userName"`           // Username of the user.
	GithubName         string    `json:"githubName"`         // GitHub username of the user.
	GithubId           string    `json:"githubId"`           // GitHub ID of the user.
	FirstName          string    `json:"firstName"`          // First name of the user.
	LastName           string    `json:"lastName"`           // Last name of the user.
	AvatarUrl          string    `json:"avatarUrl"`          // URL of the user's avatar.
	AccessToken        string    `json:"accessToken"`        // Access token for authentication.
	Bio                string    `json:"bio"`                // Biography of the user.
	Email              string    `json:"email"`              // Email address of the user.
	Contact            string    `json:"contact"`            // Contact number of the user.
	ExpiresAt          time.Time `json:"expiresAt"`          // Expiration time of the user's session or token.
	Password           string    `json:"password"`           // Password of the user.
	Verified           bool      `json:"verified"`           // Verification status of the user.
	SubscriptionStatus string    `json:"subscriptionStatus"` // Subscription status of the user.
	SubscriptionID     float64   `json:"subscriptionId"`     // Subscription ID of the user.
	SubscriptionType   string    `json:"subscriptionType"`   // Subscription type of the user.
}

// Models wraps all the models in the application for easy access.
type Models struct {
	User User // User model instance.
}

// NewModels initializes a new instance of Models with a database connection.
func NewModels(conn *pgx.Conn) Models {
	ensureTableExists(conn) // Ensure the table exists in the database.
	return Models{
		User: User{}, // Initialize the User model.
	}
}

// This functions ensures that a table exists on startup
// If the table does not exist, it creates the table
func ensureTableExists(conn *pgx.Conn) {
	query := `
    DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        user_name VARCHAR(255) NOT NULL CHECK (user_name ~ '^[A-Za-z]+$'),
        github_name VARCHAR(255) UNIQUE NOT NULL,
        github_id VARCHAR(255),
        first_name VARCHAR(255) CHECK (first_name ~ '^[A-Za-z ]+$'),
        last_name VARCHAR(255) CHECK (last_name ~ '^[A-Za-z ]+$'),
        avatar_url TEXT,
        access_token TEXT,
        bio VARCHAR(500),
        email VARCHAR(255) NOT NULL UNIQUE CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
        expires_at TIMESTAMP NOT NULL,
        password VARCHAR(255) NOT NULL,
        contact VARCHAR(255) UNIQUE CHECK (contact ~ '^\+91[0-9]{10}$'),
        verified BOOLEAN DEFAULT FALSE,
        subscription_status VARCHAR(255),
        subscription_id FLOAT UNIQUE,
        subscription_type VARCHAR(255)
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
	if u.Contact != "" {
		if !isValidIndianPhoneNumber(u.Contact) {
			return false, "contact number is invalid"
		}

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

// isValidIndianPhoneNumber checks if the phone number provided matches the Indian phone number format.
// Returns true if the phone number is valid, false otherwise.
func isValidIndianPhoneNumber(phoneNumber string) bool {
	// Regex for validating an Indian phone number
	// Starts with 7, 8, or 9 and followed by 9 more digits
	regex := regexp.MustCompile(`^[789]\d{9}$`)
	return regex.MatchString(phoneNumber)
}

// InsertUser inserts a new user into the database.
// This method is useful for registering new users in the system.
// Parameters:
// - user: The User struct containing the user's information.
// Returns:
// - An error if the query execution fails.
func (u *User) InsertUser(connection *pgx.Conn, user User) error {
	if user.Contact != "" {
		user.Contact = "+91" + user.Contact
	}
	// SQL query to insert a new user, returning the generated ID.
	query := `INSERT INTO users (user_name, github_name, github_id, first_name, last_name, avatar_url, bio, email,contact, expires_at,password,verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,$12) RETURNING id`
	// Execute the query and scan the returned ID into the User struct.
	err := connection.QueryRow(context.Background(), query, user.UserName, user.GithubName, user.GithubId, user.FirstName, user.LastName, user.AvatarUrl, user.Bio, user.Email, user.Contact, user.ExpiresAt, user.Password, user.Verified).Scan(&u.ID)
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
func (u *User) GetUser(connection *pgx.Conn, id int64) error {
	// SQL query to select a user by ID.
	query := `SELECT id, user_name, github_name, github_id, first_name, last_name, avatar_url, bio, email,contact,verified FROM users WHERE id=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.UserName, &u.GithubName, &u.GithubId, &u.FirstName, &u.LastName, &u.AvatarUrl, &u.Bio, &u.Email, &u.Contact, &u.Verified)
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
func (u *User) UpdateUser(connection *pgx.Conn, id int64, updatedUser User) error {
	baseQuery := "UPDATE users SET "
	var args []interface{}
	var updates []string
	argCounter := 1

	// For each field, check if provided and append to query and args as necessary
	if updatedUser.UserName != "" {
		updates = append(updates, fmt.Sprintf("user_name=$%d", argCounter))
		args = append(args, updatedUser.UserName)
		argCounter++
	}
	if updatedUser.GithubName != "" {
		updates = append(updates, fmt.Sprintf("github_name=$%d", argCounter))
		args = append(args, updatedUser.GithubName)
		argCounter++
	}
	if updatedUser.GithubId != "" {
		updates = append(updates, fmt.Sprintf("github_id=$%d", argCounter))
		args = append(args, updatedUser.GithubId)
		argCounter++
	}
	if updatedUser.FirstName != "" {
		updates = append(updates, fmt.Sprintf("first_name=$%d", argCounter))
		args = append(args, updatedUser.FirstName)
		argCounter++
	}
	if updatedUser.LastName != "" {
		updates = append(updates, fmt.Sprintf("last_name=$%d", argCounter))
		args = append(args, updatedUser.LastName)
		argCounter++
	}
	if updatedUser.AvatarUrl != "" {
		updates = append(updates, fmt.Sprintf("avatar_url=$%d", argCounter))
		args = append(args, updatedUser.AvatarUrl)
		argCounter++
	}
	if updatedUser.Bio != "" {
		updates = append(updates, fmt.Sprintf("bio=$%d", argCounter))
		args = append(args, updatedUser.Bio)
		argCounter++
	}
	if updatedUser.Email != "" {
		updates = append(updates, fmt.Sprintf("email=$%d", argCounter))
		args = append(args, updatedUser.Email)
		argCounter++
	}
	if updatedUser.Contact != "" {
		updates = append(updates, fmt.Sprintf("contact=$%d", argCounter))
		args = append(args, updatedUser.Contact)
		argCounter++
	}
	if !updatedUser.ExpiresAt.IsZero() {
		updates = append(updates, fmt.Sprintf("expires_at=$%d", argCounter))
		args = append(args, updatedUser.ExpiresAt)
		argCounter++
	}
	if updatedUser.Password != "" {
		updates = append(updates, fmt.Sprintf("password=$%d", argCounter))
		args = append(args, updatedUser.Password)
		argCounter++
	}
	if updatedUser.Verified {
		updates = append(updates, fmt.Sprintf("verified=$%d", argCounter))
		args = append(args, updatedUser.Verified)
		argCounter++
	}
	if updatedUser.SubscriptionStatus != "" {
		updates = append(updates, fmt.Sprintf("subscription_status=$%d", argCounter))
		args = append(args, updatedUser.SubscriptionStatus)
		argCounter++
	}
	if updatedUser.SubscriptionID != 0 {
		updates = append(updates, fmt.Sprintf("subscription_id=$%d", argCounter))
		args = append(args, updatedUser.SubscriptionID)
		argCounter++
	}
	if updatedUser.SubscriptionType != "" {
		updates = append(updates, fmt.Sprintf("subscription_type=$%d", argCounter))
		args = append(args, updatedUser.SubscriptionType)
		argCounter++
	}

	// Finalize query
	if len(updates) == 0 {
		return nil // No updates to make
	}
	query := baseQuery + strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id=$%d", argCounter)
	args = append(args, id)

	// Execute the query
	cmdTag, err := connection.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	// Check if no rows were affected.
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}
	return nil
}

// UpdateUserSubscription updates an existing user's subscription status and name in the database.
// Parameters:
// - id: The ID of the user whose subscription is to be updated.
// - status: The new subscription status.
// - name: The new subscription name.
// Returns:
// - An error if the query execution fails.
func (u *User) UpdateUserSubscription(connection *pgx.Conn, id int64, subscriptionStatus string, subscriptionId float64, subscriptionType string) error {
	// SQL query to update a user's subscription status and name by ID.
	query := `UPDATE users SET subscription_status=$1, subscription_id=$2, subscription_type=$3 WHERE id=$4`
	// Execute the query without returning any result.
	cmdTag, err := connection.Exec(context.Background(), query, subscriptionStatus, subscriptionId, subscriptionType, id)
	if err != nil {
		return err // Return any errors encountered.
	}
	// Check if no rows were affected.
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}
	return nil // Return nil on success.
}

// DeleteUser removes a user from the database by their ID.
// This method is useful for removing a user from the system.
// Parameters:
// - id: The ID of the user to delete.
// Returns:
// - An error if the query execution fails.s
func (u *User) DeleteUser(connection *pgx.Conn, id int64) error {
	// SQL query to delete a user by ID.
	query := `DELETE FROM users WHERE id=$1`
	// Execute the query.
	cmdTag, err := connection.Exec(context.Background(), query, id)
	if err != nil {
		return err // Return any errors encountered.
	}
	// Check if no rows were affected.
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", id)
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
func (u *User) GetByGitId(connection *pgx.Conn, githubId string) error {
	// SQL query to select a user by GitHub ID.
	query := `SELECT id, user_name, github_name, github_id, first_name, last_name, avatar_url, bio, email,contact,verified FROM users WHERE github_id=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, githubId).Scan(&u.ID, &u.UserName, &u.GithubName, &u.GithubId, &u.FirstName, &u.LastName, &u.AvatarUrl, &u.Bio, &u.Email, &u.Contact, &u.Verified)
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
func (u *User) GetByEmail(connection *pgx.Conn, email string) error {
	// SQL query to select a user by GitHub ID.
	query := `SELECT id, user_name,password,email,contact FROM users WHERE email=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, email).Scan(&u.ID, &u.UserName, &u.Password, &u.Email, &u.Contact)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}

// GetByContact retrieves a user by their contact number from the database.
// This method is useful for authenticating users based on their email address,
// allowing the application to fetch user details based on their email.
// Parameters:
// - email: The email address of the user to retrieve.
// Returns:
// - nil if the user is successfully found and the User struct is populated.
// - An error if the query execution or scan fails.
func (u *User) GetByContact(connection *pgx.Conn, contact string) error {
	// SQL query to select a user by GitHub ID.
	query := `SELECT id, user_name,password,email,contact FROM users WHERE contact=$1`
	// Execute the query and scan the result into the User struct.
	err := connection.QueryRow(context.Background(), query, contact).Scan(&u.ID, &u.UserName, &u.Password, &u.Email, &u.Contact)
	if err != nil {
		return err // Return any errors encountered.
	}
	return nil // Return nil on success.
}
