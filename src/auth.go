package src

import (
	"net/http"
	"os"
)

func CheckAdminPass(r *http.Request) bool {
	// Get the password from the environment
	adminPass := os.Getenv("ADMIN_PASS")

	// Get the password from the cookies
	cookie, err := r.Cookie("admin_pass")
	if err != nil {
		return false
	}

	// Compare the passwords
	return cookie.Value == adminPass
}
