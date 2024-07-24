package test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"subscription-service/data"
	"testing"

	"github.com/jackc/pgx/v4"
)

type UserTestSuite struct {
	connection *pgx.Conn
	userID     []int64
}

func (suite *UserTestSuite) SetupSuite() {
	url := "postgres://root@localhost:26257/defaultdb?sslmode=disable" // Database connection URL.
	conn, err := pgx.Connect(context.Background(), url)                // Attempt to connect to the database.
	if err != nil {
		log.Panic(err) // Panic if the connection fails.
	}
	fmt.Println("Connected to the database")
	ensureTableExists(conn) // Ensure the table exists in the database.
	suite.connection = conn // Assign the connection to the global variable.
}

func (suite *UserTestSuite) TeardownSuite() {
	suite.connection.Close(context.Background())

}

func ensureTableExists(conn *pgx.Conn) {
	query := `
    DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        user_name VARCHAR(255) NOT NULL CHECK (user_name ~ '^[A-Za-z]+$'),
        github_name VARCHAR(255) UNIQUE NOT NULL,
        github_id VARCHAR(255),
        first_name VARCHAR(255) CHECK (first_name ~ '^[A-Za-z ]+$'),
        last_name VARCHAR(255)  CHECK (last_name ~ '^[A-Za-z ]+$'),
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

type TestUser struct {
	ID          int64  `json:"id"`
	UserName    string `json:"username"`
	GithubName  string `json:"githubname"`
	FirstName   string `json:"firstname"`
	SecondName  string `json:"secondname"`
	Bio         string `json:"bio"`
	Email       string `json:"email"`
	Contact     string `json:"contact"`
	Password    string `json:"password"`
	AccessToken string `json:"accesstoken"`
	Verified    bool   `json:"verified"`
}

func (suite *UserTestSuite) TestInsertUser(t *testing.T) {
	testCases := []struct {
		name    string
		user    TestUser
		wantErr bool
	}{
		{
			name: "ValidUserCompleteInfo",
			user: TestUser{
				ID:          1,
				UserName:    "completeInfo",
				GithubName:  "githubComplete",
				FirstName:   "John",
				SecondName:  "Doe",
				Bio:         "Full Stack Developer",
				Email:       "john.doe@example.com",
				Contact:     "1234567890",
				Password:    "securePassword123",
				AccessToken: "validAccessToken",
				Verified:    true,
			},
			wantErr: false,
		},
		{
			name: "ValidUserTechEnthusiast",
			user: TestUser{
				ID:          2,
				UserName:    "techEnthusiast",
				GithubName:  "githubTech",
				FirstName:   "Alice",
				SecondName:  "Smith",
				Bio:         "Tech enthusiast and blogger",
				Email:       "alice.smith@example.com",
				Contact:     "9876543210",
				Password:    "techLover2023",
				AccessToken: "validTechToken",
				Verified:    true,
			},
			wantErr: false,
		},
		{
			name: "ValidUserDesigner",
			user: TestUser{
				ID:          3,
				UserName:    "creativeDesigner",
				GithubName:  "githubDesign",
				FirstName:   "Bob",
				SecondName:  "Brown",
				Bio:         "Creative designer specializing in UX/UI",
				Email:       "bob.brown@example.com",
				Contact:     "1122334455",
				Password:    "designIsLife",
				AccessToken: "validDesignToken",
				Verified:    true,
			},
			wantErr: false,
		},
		{
			name: "MissingEmail",
			user: TestUser{
				ID:          2,
				UserName:    "missingEmail",
				GithubName:  "githubNoEmail",
				FirstName:   "Jane",
				SecondName:  "Doe",
				Bio:         "Data Analyst",
				Contact:     "0987654321",
				Password:    "password123",
				AccessToken: "accessToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "InvalidContactNumber",
			user: TestUser{
				ID:          3,
				UserName:    "invalidContact",
				GithubName:  "githubInvalidContact",
				FirstName:   "Alice",
				SecondName:  "Smith",
				Bio:         "Cybersecurity Expert",
				Email:       "alice.smith@example.com",
				Contact:     "12345",
				Password:    "aliceSecure123",
				AccessToken: "aliceToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "EmptyUserName",
			user: TestUser{
				ID:          4,
				UserName:    "",
				GithubName:  "githubEmptyUser",
				FirstName:   "Bob",
				SecondName:  "Brown",
				Bio:         "Cloud Architect",
				Email:       "bob.brown@example.com",
				Contact:     "9876543210",
				Password:    "bobPassword123",
				AccessToken: "bobToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "DuplicateGithubUserName",
			user: TestUser{
				ID:          5,
				UserName:    "duplicateUser",
				GithubName:  "githubTech",
				FirstName:   "Charlie",
				SecondName:  "Green",
				Bio:         "Machine Learning Engineer",
				Email:       "charlie.green@example.com",
				Contact:     "4564564567",
				Password:    "charliePass123",
				AccessToken: "charlieToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "InvalidEmailFormat",
			user: TestUser{
				ID:          6,
				UserName:    "invalidEmailFormat",
				GithubName:  "githubInvalidEmail",
				FirstName:   "Diana",
				SecondName:  "White",
				Bio:         "Blockchain Developer",
				Email:       "diana.email.com", // Missing '@'
				Contact:     "1231231234",
				Password:    "dianaPassword123",
				AccessToken: "dianaToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "LongBio",
			user: TestUser{
				ID:          8,
				UserName:    "longBioUser",
				GithubName:  "githubLongBio",
				FirstName:   "Fiona",
				SecondName:  "Grey",
				Bio:         strings.Repeat("A", 1001), // Assuming bio has a max length of 1000 characters
				Email:       "fiona.grey@example.com",
				Contact:     "3213214321",
				Password:    "fionaPassword123",
				AccessToken: "fionaToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "NoFirstName",
			user: TestUser{
				ID:          9,
				UserName:    "noFirstName",
				GithubName:  "githubNoFirstName",
				FirstName:   "",
				SecondName:  "Yellow",
				Bio:         "UI/UX Designer",
				Email:       "no.first.name@example.com",
				Contact:     "7897897890",
				Password:    "noFirstName123",
				AccessToken: "noFirstNameToken123",
				Verified:    true,
			},
			wantErr: true,
		},
		{
			name: "SpecialCharactersInUserName",
			user: TestUser{
				ID:          10,
				UserName:    "special@User#Name",
				GithubName:  "githubSpecialUser",
				FirstName:   "George",
				SecondName:  "Violet",
				Bio:         "Software Tester",
				Email:       "george.violet@example.com",
				Contact:     "9879879870",
				Password:    "georgePassword123",
				AccessToken: "georgeToken123",
				Verified:    true,
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{
				UserName:    tc.user.UserName,
				GithubName:  tc.user.GithubName,
				FirstName:   tc.user.FirstName,
				LastName:    tc.user.SecondName,
				Bio:         tc.user.Bio,
				Email:       tc.user.Email,
				Contact:     tc.user.Contact,
				Password:    tc.user.Password,
				AccessToken: tc.user.AccessToken,
				Verified:    tc.user.Verified,
			}

			err := u.InsertUser(suite.connection, u)
			if err != nil && !tc.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil && tc.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			suite.userID = append(suite.userID, u.ID)
		})
	}

}

func (suite *UserTestSuite) TestGetUser(t *testing.T) {
	testCases := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{
			name:    "Valid User ID",
			userID:  suite.userID[0],
			wantErr: false,
		},
		{
			name:    "User ID Does Not Exist",
			userID:  999999, // Assuming this ID does not exist
			wantErr: true,
		},
		{
			name:    "Negative User ID",
			userID:  -1,
			wantErr: true,
		},
		{
			name:    "Zero User ID",
			userID:  0,
			wantErr: true,
		},
		{
			name:    "Very Large User ID",
			userID:  9223372036854775807, // Max int64 value
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{}

			err := u.GetUser(suite.connection, tc.userID)
			if err != nil && !tc.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil && tc.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			// New assertion: Check if u.ID equals tc.userID when there's no error
			if err == nil && u.ID != tc.userID {
				t.Errorf("GetUser() got user ID %v, want %v", u.ID, tc.userID)
			}

		})
	}
}

func (suite *UserTestSuite) TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name    string
		userID  int64
		update  data.User
		wantErr bool
	}{
		{
			name:   "UpdateEmailValid",
			userID: suite.userID[0],
			update: data.User{
				Email: "john.updated@example.com",
			},
			wantErr: false,
		},
		{
			name:   "UpdateContactInvalid",
			userID: suite.userID[1],
			update: data.User{
				Contact: "123", // Invalid contact number
			},
			wantErr: true,
		},
		{
			name:   "UpdateFirstNameEmpty",
			userID: suite.userID[2],
			update: data.User{
				FirstName: "",
			},
			wantErr: false,
		},
		{
			name:   "UpdateLastNameValid",
			userID: suite.userID[0],
			update: data.User{
				LastName: "SmithUpdated",
			},
			wantErr: false,
		},
		{
			name:   "UpdateBioTooLong",
			userID: suite.userID[0],
			update: data.User{
				Bio: strings.Repeat("B", 1001), // Assuming bio has a max length of 1000 characters
			},
			wantErr: true,
		},
		{
			name:   "UpdateUserNameWithSpecialCharacters",
			userID: suite.userID[0],
			update: data.User{
				UserName: "new@User#Name",
			},
			wantErr: true,
		},
		{
			name:   "UpdateVerifiedStatus",
			userID: suite.userID[2],
			update: data.User{
				Verified: false,
			},
			wantErr: false,
		},
		{
			name:   "UpdatePasswordValid",
			userID: suite.userID[1],
			update: data.User{
				Password: "newSecurePassword123",
			},
			wantErr: false,
		},
		{
			name:   "UpdateGithubNameValid",
			userID: suite.userID[2],
			update: data.User{
				GithubName: "newGithubName2321",
			},
			wantErr: false,
		},
		{
			name:   "UpdateEmailInvalid",
			userID: suite.userID[2],
			update: data.User{
				Email: "invalidEmail", // Invalid email format
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{}
			err := u.UpdateUser(suite.connection, tc.userID, tc.update)
			if err != nil && !tc.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil && tc.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func (suite *UserTestSuite) TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name    string
		userID  int64
		setup   func() // Optional setup function if needed to insert user before deletion
		wantErr bool
	}{
		{
			name:   "ValidUser",
			userID: suite.userID[0],
			setup: func() {
				// Insert user with ID 1 into the database for deletion
			},
			wantErr: false,
		},
		{
			name:    "NonExistentUser",
			userID:  999, // Assuming this ID does not exist in the database
			wantErr: true,
		},
		{
			name:    "NegativeUserID",
			userID:  -1,
			wantErr: true,
		},
		{
			name:    "ZeroUserID",
			userID:  0,
			wantErr: true,
		},
		{
			name:   "UserWithAssociatedData",
			userID: 2,
			setup: func() {
				// Insert user with ID 2 and associated data (e.g., posts, comments) that might prevent deletion
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{}
			err := u.DeleteUser(suite.connection, tc.userID)
			t.Log(err)
			if err != nil && !tc.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil && tc.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil {
				err = u.GetUser(suite.connection, tc.userID)
				if err == nil {
					t.Errorf("DeleteUser() user with ID %v still exists", tc.userID)
				}
			}
		})
	}
}

func (suite *UserTestSuite) UpdateSubscription(t *testing.T) {
	var testCases = []struct {
		name               string
		id                 int64
		subscriptionStatus string
		subscriptionId     float64
		subscriptionType   string
		expectedError      bool
	}{
		{
			name:               "Successful Update",
			id:                 suite.userID[0],
			subscriptionStatus: "active",
			subscriptionId:     123.45,
			subscriptionType:   "premium",
			expectedError:      false,
		},
		{
			name:               "Non-existent User",
			id:                 999,
			subscriptionStatus: "active",
			subscriptionId:     123.45,
			subscriptionType:   "premium",
			expectedError:      true, // Assuming no error for no rows updated
		},
		{
			name:               "Invalid Subscription ID",
			id:                 2,
			subscriptionStatus: "active",
			subscriptionId:     -1, // Invalid ID
			subscriptionType:   "basic",
			expectedError:      true, // Assuming your validation catches this before making a DB call
		},
		{
			name:               "Database Error",
			id:                 3,
			subscriptionStatus: "inactive",
			subscriptionId:     234.56,
			subscriptionType:   "standard",
			expectedError:      true,
		},
		{
			name:               "Update With Same Data",
			id:                 suite.userID[2],
			subscriptionStatus: "active",
			subscriptionId:     123.45,
			subscriptionType:   "premium",
			expectedError:      true, // Assuming no error for no rows updated
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		u := data.User{}
		err := u.UpdateUserSubscription(suite.connection, tc.id, tc.subscriptionStatus, tc.subscriptionId, tc.subscriptionType)

		if err != nil && !tc.expectedError {
			t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tc.expectedError)
			return
		}

		if err == nil && tc.expectedError {
			t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tc.expectedError)
			return
		}
	}
}
func (suite *UserTestSuite) TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid Email",
			email:   "john.doe@example.com", // Assuming this is a valid email in your test database
			wantErr: false,
		},
		{
			name:    "Invalid Email Format",
			email:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "Non-Existing Email",
			email:   "nonexisting@example.com",
			wantErr: true,
		},
		{
			name:    "Empty Email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{}
			err := u.GetByEmail(suite.connection, tc.email)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetUserByEmail() for email %v, error = %v, wantErr %v", tc.email, err, tc.wantErr)
			}

			if err == nil && u.Email != tc.email {
				t.Errorf("GetUserByEmail() for email %v, got email %v, want %v", tc.email, u.Email, tc.email)
			}
		})
	}
}

func (suite *UserTestSuite) TestGetUserByContact(t *testing.T) {
	testCases := []struct {
		name    string
		contact string
		wantErr bool
	}{
		{
			name:    "Valid Contact",
			contact: "+919876543210", // Assuming this is a valid contact in your test database
			wantErr: false,
		},
		{
			name:    "Invalid Contact Format",
			contact: "9876543210", // Missing country code
			wantErr: true,
		},
		{
			name:    "Non-Existing Contact",
			contact: "+911224567890", // Assuming this contact does not exist in your test database
			wantErr: true,
		},
		{
			name:    "Empty Contact",
			contact: "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := data.User{}
			err := u.GetByContact(suite.connection, tc.contact)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetUserByContact() for contact %v, error = %v, wantErr %v", tc.contact, err, tc.wantErr)
			}

			if err == nil && u.Contact != tc.contact {
				t.Errorf("GetUserByContact() for contact %v, got contact %v, want %v", tc.contact, u.Contact, tc.contact)
			}
		})
	}
}

func TestUserSuite(t *testing.T) {
	user_suite := UserTestSuite{}
	user_suite.SetupSuite()
	defer user_suite.TeardownSuite()
	t.Run("TestInsertUser", user_suite.TestInsertUser)
	t.Run("TestGetUser", user_suite.TestGetUser)
	t.Run("UpdateSubscription", user_suite.UpdateSubscription)
	t.Run("TestGetUserByEmail", user_suite.TestGetUserByEmail)
	t.Run("TestGetUserByContact", user_suite.TestGetUserByContact)
	t.Run("TestUpdateUser", user_suite.TestUpdateUser)
	t.Run("TestDeleteUser", user_suite.TestDeleteUser)

}
