package main

import (
	"database/sql"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// userLogin() checks entered credentials
// against the db.  If no match is found
// an empty user, and error is returned,
// otherwise a new *User is returned.
func userLogin(db *sql.DB, w http.ResponseWriter, r *http.Request) (string, *User, error) {
	var user *User
	username := r.FormValue("username")
	enteredPassword := []byte(r.FormValue("password"))
	var storedPassword []byte
	sLog(fmt.Sprintf("auth.go: userLogin(): username: %v", username))

	// Check the database for user and retrieve password.
	err := db.QueryRow("SELECT secret FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		return "Wrong username or password", user, fmt.Errorf("auth.go: userLogin(): dbQueryRow(): error: %v", err)
	}
	sLog(fmt.Sprintf("auth.go: userLogin(): db.QueryRow successful username match: %v", username))

	// Compare the stored password to the entered password.
	err = passwordCompare(storedPassword, enteredPassword)
	if err != nil {
		return "Wrong username or password.", user, fmt.Errorf("auth.go: userLogin(): dbQueryRow(): error: %v", err)
	}
	sLog(fmt.Sprintf("auth.go: userLogin(): call to passwordCompare(): successful password match for user: %v", username))

	// Once password and username are matched serve cookie.
	user = &User{Username: username, Secret: enteredPassword, CookieSession: serveCookie(w, r)}
	err = storeSession(db, user)
	if err != nil {
		sLog(fmt.Sprintf("auth.go: userLogin(): call to storeSession error: %v", err))
		return "userExists", nil, nil
	}
	sLog(fmt.Sprintf("auth.go: userLogin(): call to storeSession success: %v", user))
	return "User login success.", user, nil
}

// passwordCompare() returns an error on failure
// and compares the password stored in the db
// with the password entered by user.
// "golang.org/x/crypto/bcrypt" is used for encryption.
func passwordCompare(storedPassword, enteredPassword []byte) error {
	sLog("auth.go: passwordCompare()")
	if err := bcrypt.CompareHashAndPassword(storedPassword, enteredPassword); err != nil {
		return fmt.Errorf("auth.go: passwordCompare(): error: %v", err)
	}
	sLog("auth.go: passwordCompare() successful match:")
	return nil
}

// registerUser() ensures that the requested username is unique
// adds the new user to the db, serves a cookie, returns a success message,
// and error on failure.
func registerUser(db *sql.DB, w http.ResponseWriter, r *http.Request) (string, *User, error) {
	var user *User
	sLog(fmt.Sprintf("auth.go: registerUser(): r: %v", r))
	// Set username and password variables to compare with potential db entries
	// and to set the db entry if registration is successful.
	username := r.FormValue("username")
	user = &User{Secret: []byte(r.FormValue("password"))}
	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&username)
	if err != sql.ErrNoRows {
		return "userExists", nil, nil
	}

	// If the response from the db matches the
	// requested username the request is denied
	// and a message asking for a different username
	// is delivered.
	sLog(fmt.Sprintf("auth.go: registerUser(): username: %v", username))
	user.CookieSession = serveCookie(w, r)
	err = storeSession(db, user)
	if err != nil {
		sLog(fmt.Sprintf("auth.go: registerUser(): call to storeSession error: %v", err))
		return "userExists", nil, nil
	}
	sLog(fmt.Sprintf("auth.go: registerUser(): registration successful: %v", username))
	// Return the struct with new user data.
	return "Registration of new user successful", user, err
}

// storeSession() takes a cookie value and stores it to the user database.
func storeSession(db *sql.DB, user *User) error {
	// After creating the cookie session data
	// the requested password is encrypted using bcrypt.
	user.Secret, err = bcrypt.GenerateFromPassword(user.Secret, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("auth.go: registerUser(): error encrypting password secret: %v: error: %v", user.Secret, err)
	}
	sLog(fmt.Sprintf("auth.go: registerUser(): encrypt password: %v", user.Secret))
	// Prepare the db insert statement for the new user.
	registerStmt, err := db.Prepare("INSERT users SET username=?,secret=?,cookieSession=?")
	if err != nil {
		return fmt.Errorf("auth.go: registerUser(): registerStmt: %v: error: %v", registerStmt, err)
	}
	// Execute the db insert statement for the new user.
	_, err = registerStmt.Exec(user.Username, user.Secret, user.CookieSession)
	if err != nil {
		return fmt.Errorf("auth.go: registerUser(): storeSession(): error executing registerStmt: %v: error: %v", registerStmt, err)
	}
	return nil
}

// serveCookie() creates, serves, and returns a cookie.
func serveCookie(w http.ResponseWriter, r *http.Request) string {
	sLog("auth.go: serveCookie():")
	// r.Cookie tries to retrieve the cookie.
	// a new cookie is created on failure.
	cookie, err := r.Cookie("golang-event-planner")
	if err != nil {
		cookieID, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:   "golang-event-planner",
			Value:  cookieID.String(),
			MaxAge: 0,
		}
	}
	// Serve cookie to user.
	http.SetCookie(w, cookie)
	sLog(fmt.Sprintf("auth.go: serveCookie(): success: cookie: %v", cookie))
	return cookie.Value
}

// verifySession() checks for an existing cookie,
// attempts to match the found cookie to the db,
// and if it finds none or the cookie doesn't match
// the user is redirected to the login page.
func verifySession(db *sql.DB, r *http.Request) (*User, error) {
	sLog(fmt.Sprintf("auth.go: verifySession(): db: %v", db))
	var username string
	var user *User
	// Read the cookie from the user.
	// If no cookie is found the user
	// is sent to the login/register page.
	cookie, err := r.Cookie("golang-event-planner")
	if err != nil {
		// Cookie was not found.
		return nil, fmt.Errorf("auth.go: verifySession(): error getting cookie: %v", err)
	}
	sLog(fmt.Sprintf("auth.go: verifySession(): found cookie: %v", cookie))
	// If a cookie was found, retrieve the user data and return user.
	err = db.QueryRow("SELECT username FROM users WHERE cookieSession=?", cookie.Value).Scan(&username)
	if err != nil {
		// Could not find cookie match.
		sLog(fmt.Sprintf("auth.go: verifySession(): db.QueryRow(): fail: %v", err))
		return nil, fmt.Errorf("auth.go: verifySession(): db.QueryRow(): cookie not found in db error: %v", err)
	}
	user = &User{Username: username, CookieSession: cookie.Value}
	sLog(fmt.Sprintf("auth.go: verifySession(): success: user: %v", user))
	return user, nil
}
