package data

import (
	"github.com/jackc/pgx/v4"
)

type User struct {
	ID          int64  // Unique identifier for the user
	GitHubID    string // GitHub user ID
	GitHubLogin string // GitHub username
	AccessToken string // GitHub access token for making authenticated requests
	Email       string // User's email address
	Name        string // User's names
	AvatarURL   string // User's avatar URL
}

var connection *pgx.Conn

type Models struct {
	User User // Assuming UserModel is a struct that handles database operations for User entities
}

// New initializes a new Models instance with a database connection.
func New(conn *pgx.Conn) Models {
	connection = conn
	return Models{
		User: User{},
	}
}
