// Package activity contains the business logic for user activities.
package activity

// Import the data package from the subscription-service project.
// This package likely contains definitions for interacting with the database.
import "subscription-service/data"

// UserResponse is a struct used to represent the response structure for user-related queries.
// It holds information about the user's ID, contact details, and subscription status.
type UserResponse struct {
	ID                 int64   // ID is the unique identifier for the user.
	Contact            string  // Contact represents the user's contact information, such as an email address.
	SubscriptionStatus float64 // SubscriptionStatus holds the subscription status as a float64. This might represent a subscription ID or status code.
}

// GetUser is a method on ActivitiesImpl (not shown here but presumably a struct related to user activities)
// that retrieves a user's information based on their email address.
func (ac *ActivitiesImpl) GetUser(email string) (UserResponse, error) {
	// Initialize an empty User struct from the data package.
	// This struct likely contains methods for database operations related to a User.
	user := data.User{}

	// Call the GetByEmail method on the user object, passing in the email address.
	// This method is expected to populate the user object with data from the database.
	err := user.GetByEmail(email)

	// If there was an error retrieving the user, return an empty UserResponse and the error.
	if err != nil {
		return UserResponse{}, err
	}

	// If the user was successfully retrieved, populate a UserResponse struct with the user's data
	// and return it along with a nil error.
	return UserResponse{
		ID:                 user.ID,             // User's unique identifier.
		Contact:            user.Contact,        // User's contact information.
		SubscriptionStatus: user.SubscriptionID, // User's subscription status or ID.
	}, nil
}

// UpdateSubscription is a method on ActivitiesImpl that updates a user's subscription information.
func (ac *ActivitiesImpl) UpdateSubscription(id int64, subscriptionStatus string, subscriptionId float64, subscriptionType string) error {
	// Initialize an empty User struct from the data package.
	user := data.User{}

	// Call the UpdateUserSubscription method on the user object, passing in the user's ID,
	// new subscription status, subscription ID, and subscription type.
	// This method is expected to update the user's subscription information in the database.
	err := user.UpdateUserSubscription(id, subscriptionStatus, subscriptionId, subscriptionType)

	// If there was an error updating the user's subscription, return the error.
	if err != nil {
		return err
	}

	// If the update was successful, return a nil error.
	return nil
}
