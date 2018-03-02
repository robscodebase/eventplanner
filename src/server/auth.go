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

func userLogin(w http.ResponseWriter, r *http.Request) (string, *User, error) {
	sLog("auth.go: userLogin()")
	// Check the database for user and retrieve password.
	username := r.FormValue("username")
	enteredPassword := []byte(r.FormValue("password"))
	sLog(fmt.Sprintf("auth.go: userLogin(): username: %v", username))
	err := db.QueryRow(planner.LogGet, username).Scan(&user.Username, &user.Secret)
	if err != nil {
		return "Please try again.", user, fmt.Errorf("auth.go: userLogin(): dbQueryRow(): error: %v", err)
	}
	err = passwordCompare(user.Secret, enteredPassword)
	if err != nil {
		return "Please try again.", user, fmt.Errorf("auth.go: userLogin(): dbQueryRow(): error: %v", err)
	}
}

func passwordCompare(storedPassword, enteredPassword []byte) error {
	sLog("auth.go: passwordCompare()")
	if err := bcrypt.CompareHashAndPassword(storedPassword, enteredPassword); err != nil {
		return sLog(fmt.Sprintf("auth.go: passwordCompare(): error: %v", err))
	}
	return nil
}

func registerUser(w http.ResponseWriter, r *http.Request) (string, *User, error) {
	sLog("auth.go: registerUser()")
	// Set username and password variables to compare with potential db entries
	// and to set the db entry if registration is successful.
	username := r.FormValue("username")
	password := r.FormValue("password")

	// user is used to get a response from the db.
	// If the response from the db matches the
	// requested username the request is rejected.
	var user *User
	username := r.FormValue("username")
	sLog(fmt.Sprintf("auth.go: registerUser(): username: %v", username))
	db.QueryRow("SELECT username, secret FROM users WHERE username = ?", user.Username).Scan(&user.Username)
	if user.Username == username {
		sLog(fmt.Sprintf("auth.go: registerUser(): user already exists: %v", username))
		return "Username already exists. Please choose a different username.", nil
	} else {
		cookie := plantCookie(w, r)
		user.CookieSession = cookie.Value
		// After creating the cookie session data
		// the requested password is encrypted using bcrypt.
		sLog(fmt.Sprintf("auth.go: registerUser(): set cookie: %v", cookie))
		secret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Sprintf("auth.go: registerUser(): error encrypting password secret: %v: error: %v", secret, err), user, err
		}
		sLog(fmt.Sprintf("auth.go: registerUser(): encrypt password: %v", secret))
		// Prepare the db insert statement for the new user.
		registerStmt, err := db.Prepare("INSERT users SET username=?,secret=?,cookieSession=?")
		if err != nil {
			return fmt.Sprintf("auth.go: registerUser(): registerStmt: %v: error: %v", registerStmt, err), user, err
		}
		// Execute the db insert statement for the new user.
		_, err = registerStmt.Exec(username, secret, user.CookieSession)
		if err != nil {
			return fmt.Sprintf("auth.go: registerUser(): error executing registerStmt: %v: error: %v", registerStmt, err), user, err
		}
	}
	sLog(fmt.Sprintf("auth.go: registerUser(): registration successful: %v", username))
	// Return the struct with new user data.
	return "Registration of new user successful", &User{Username: username, Secret: secret, CookieSession: cookie.Value}, err
}

func plantCookie(w *http.ResponseWriter, r *http.Request) *http.Cookie {
	sLog("auth.go: plantCookie():")
	// If there is no cookie data create new.
	cookie, err = r.Cookie("docker-golang-mysql-event-planner")
	if err != nil {
		cookieID, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:   "docker-golang-mysql-event-planner",
			Value:  cookieID.String(),
			MaxAge: 0,
		}
	}
	http.SetCookie(w, cookie)
	sLog(fmt.Sprintf("auth.go: plantCookie(): success: cookie: %v", cookie))
	return cookie
}

func verifySession(db *sql.DB) (*User, error) {
	sLog(fmt.Sprintf("auth.go: verifySession(): db: %v", db))
	var user *User
	cookie, err := r.Cookie("golang-event-planner")
	if err != nil {
		return user, fmt.Errorf("auth.go: verifySession(): error getting cookie: %v", err)
	} else {
		user.CookieSession = cookie.Value
		err = db.QueryRow("SELECT username FROM users WHERE cookieSession=?", cookie.Value).Scan(&user.Username)
		if err != nil {
			return user, fmt.Errorf("auth.go: verifySession(): db.QueryRow(): error: %v", err)
		}
	}
	sLog(fmt.Sprintf("auth.go: verifySession(): success: user: %v", user))
	return user, nil
}
