package main

import (
	"database/sql"
	"database/sql/driver"
	"github.com/go-sql-drver/mysql"
)

func sessionData(r *http.Request) *User {
	return &User{
		UserName:      "demo",
		CookieSession: "demo",
	}
}

// userLogin() checks entered credentials
// against the db.  If no match is found
// an empty user, and error is returned,
// otherwise a new *User is returned.
func userLogin(db *sql.DB, w http.ResponseWriter, r *http.Request) (string, *User, error) {
	username := r.FormValue("username")
	enteredPassword := []byte(r.FormValue("password"))
	sLog(fmt.Sprintf("auth.go: userLogin(): username: %v", username))

	// Check the database for user and retrieve password.
	err := db.QueryRow(planner.LogGet, username).Scan(&user.Username, &user.Secret)
	if err != nil {
		return "Please try again.", user, fmt.Errorf("auth.go: userLogin(): dbQueryRow(): error: %v", err)
	}
	sLog(fmt.Sprintf("auth.go: userLogin(): db.QueryRow successful username match: %v", username))

	// Compare the stored password to the entered password.
	err = passwordCompare(user.Secret, enteredPassword)
	if err != nil {
		return "Please try again.", user, fmt.Errorf("auth.go: userLogin(): dbQueryRow(): error: %v", err)
	}
	sLog(fmt.Sprintf("auth.go: userLogin(): call to passwordCompare(): successful password match for user: %v", username))

	// Once password and username are matched serve cookie.
	cookie := serveCookie(w, r)
	user.CookieSession = cookie.Value
}

// passwordCompare() returns an error on failure
// and compares the password stored in the db
// with the password entered by user.
// "golang.org/x/crypto/bcrypt" is used for encryption.
func passwordCompare(storedPassword, enteredPassword []byte) error {
	sLog("auth.go: passwordCompare()")
	if err := bcrypt.CompareHashAndPassword(storedPassword, enteredPassword); err != nil {
		return sLog(fmt.Sprintf("auth.go: passwordCompare(): error: %v", err))
	}
	sLog("auth.go: passwordCompare() successful match:")
	return nil
}

// registerUser() ensures that the requested username is unique
// adds the new user to the db, serves a cookie, returns a success message,
// and error on failure.
func registerUser(db *sql.DB, w http.ResponseWriter, r *http.Request) (string, *User, error) {
	sLog("auth.go: registerUser()")
	// Set username and password variables to compare with potential db entries
	// and to set the db entry if registration is successful.
	username := r.FormValue("username")
	password := []byte(r.FormValue("password"))

	// If the response from the db matches the
	// requested username the request is denied
	// and a message asking for a different username
	// is delivered.
	var user *User
	username := r.FormValue("username")
	sLog(fmt.Sprintf("auth.go: registerUser(): username: %v", username))
	db.QueryRow("SELECT username, secret FROM users WHERE username = ?", user.Username).Scan(&user.Username)
	if user.Username == username {
		sLog(fmt.Sprintf("auth.go: registerUser(): user already exists: %v", username))
		return "Username already exists. Please choose a different username.", nil
	} else {
		cookie := serveCookie(w, r)
		storeSession(cookie, user)
	}
	sLog(fmt.Sprintf("auth.go: registerUser(): registration successful: %v", username))
	// Return the struct with new user data.
	return "Registration of new user successful", &User{Username: username, Secret: secret, CookieSession: cookie.Value}, err
}

// storeSession() takes a cookie and stores it to the user database.
func storeSession(db *sql.DB, cookie *http.Cookie) error {
	user.CookieSession = cookie.Value
	// After creating the cookie session data
	// the requested password is encrypted using bcrypt.
	secret, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("auth.go: registerUser(): error encrypting password secret: %v: error: %v", secret, err)
	}
	sLog(fmt.Sprintf("auth.go: registerUser(): encrypt password: %v", secret))
	// Prepare the db insert statement for the new user.
	registerStmt, err := db.Prepare("INSERT users SET username=?,secret=?,cookieSession=?")
	if err != nil {
		return fmt.Errorf("auth.go: registerUser(): registerStmt: %v: error: %v", registerStmt, err)
	}
	// Execute the db insert statement for the new user.
	_, err = registerStmt.Exec(username, secret, user.CookieSession)
	if err != nil {
		return fmt.Errorf("auth.go: registerUser(): error executing registerStmt: %v: error: %v", registerStmt, err)
	}
}

// serveCookie() creates, serves, and returns a cookie.
func serveCookie(w *http.ResponseWriter, r *http.Request) *http.Cookie {
	sLog("auth.go: serveCookie():")
	// r.Cookie tries to retrieve the cookie.
	// a new cookie is created on failure.
	cookie, err := r.Cookie("docker-golang-mysql-event-planner")
	if err != nil {
		cookieID, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:   "docker-golang-mysql-event-planner",
			Value:  cookieID.String(),
			MaxAge: 0,
		}
	}
	// Serve cookie to user.
	http.SetCookie(w, cookie)
	sLog(fmt.Sprintf("auth.go: serveCookie(): success: cookie: %v", cookie))
	return cookie
}

// verifySession() checks for an existing cookie,
// attempts to match the found cookie to the db,
// and if it finds none or the cookie doesn't match
// the user is redirected to the login page.
func verifySession(db *sql.DB) (*User, error) {
	sLog(fmt.Sprintf("auth.go: verifySession(): db: %v", db))
	var user *User
	// Read the cookie from the user.
	// If no cookie is found the user
	// is sent to the login/register page.
	cookie, err := r.Cookie("golang-event-planner")
	if err != nil {
		// Cookie was not found.
		return user, fmt.Errorf("auth.go: verifySession(): error getting cookie: %v", err)
	}
	// If a cookie was found, retrieve the user data and return user.
	user.CookieSession = cookie.Value
	err = db.QueryRow("SELECT username FROM users WHERE cookieSession=?", cookie.Value).Scan(&user.Username)
	if err != nil {
		// Could not find cookie match.
		return user, fmt.Errorf("auth.go: verifySession(): db.QueryRow(): cookie not found in db error: %v", err)
	}
	sLog(fmt.Sprintf("auth.go: verifySession(): success: user: %v", user))
	return user, nil
}
