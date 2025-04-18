package integration

import (
	"net/http"
	"testing"
)

// This test performs a complete end-to-end test of the admin API user management functionality.
// It tests the following user operations in sequence:
// 1. Creating a new user with all required fields
// 2. Retrieving user statistics
// 3. Listing users with pagination
// 4. Getting a specific user's details
// 5. Updating a user's information
// 6. Deleting a user and verifying deletion
// The test uses a test tenant "acme" and realm "customers" for all operations.

func TestUserAPI_E2E(t *testing.T) {
	e := SetupIntegrationTest(t, "")

	// Test user data
	testUser := map[string]interface{}{
		"username":     "admin_test_user",
		"display_name": "Admin Test User",
		"given_name":   "Admin",
		"family_name":  "Test",
		"email":        "admin_test@example.com",
		"phone":        "+1234567890",
		"status":       "active",
		"roles":        []string{"admin"},
		"groups":       []string{"test_group"},
		"attributes":   map[string]string{"test": "value"},
	}

	// Test creating a user
	t.Run("Create User", func(t *testing.T) {
		e.POST("/acme/customers/admin/users/"+testUser["username"].(string)).
			WithJSON(testUser).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			HasValue("username", testUser["username"]).
			HasValue("display_name", testUser["display_name"])
	})

	// Test getting user stats
	t.Run("Get User Stats", func(t *testing.T) {
		e.GET("/acme/customers/admin/users/stats").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			HasValue("total_users", 1).
			HasValue("active_users", 1)
	})

	// Test listing users with pagination
	t.Run("List Users", func(t *testing.T) {
		resp := e.GET("/acme/customers/admin/users").
			WithQuery("page", 1).
			WithQuery("page_size", 20).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object()

		// Check pagination metadata
		resp.Value("pagination").
			Object().
			HasValue("page", 1).
			HasValue("page_size", 20).
			HasValue("total_items", 1).
			HasValue("total_pages", 1)

		// Check user data
		resp.Value("data").
			Array().
			Length().
			Equal(1)
	})

	// Test getting a specific user
	t.Run("Get User", func(t *testing.T) {
		e.GET("/acme/customers/admin/users/"+testUser["username"].(string)).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			HasValue("username", testUser["username"]).
			HasValue("display_name", testUser["display_name"])
	})

	// Test updating a user
	t.Run("Update User", func(t *testing.T) {
		updatedUser := testUser
		updatedUser["display_name"] = "Updated Admin Test User"

		e.PUT("/acme/customers/admin/users/"+testUser["username"].(string)).
			WithJSON(updatedUser).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			HasValue("display_name", updatedUser["display_name"])
	})

	// Test deleting a user
	t.Run("Delete User", func(t *testing.T) {
		e.DELETE("/acme/customers/admin/users/" + testUser["username"].(string)).
			Expect().
			Status(http.StatusNoContent)

		// Verify user is deleted
		e.GET("/acme/customers/admin/users/" + testUser["username"].(string)).
			Expect().
			Status(http.StatusNotFound)
	})
}
