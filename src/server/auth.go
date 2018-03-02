package main

func sessionData(r *http.Request) *User {
	return &User{
		UserName:      "demo",
		CookieSession: "demo",
	}
}

func registerUser(w http.ResponseWriter, r *http.Request) (string, *User, error) {
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
		// If the user does not already exist set the cookie value for the session.
		cookie, err = r.Cookie("docker-golang-mysql-event-planner")
		if err != nil {
			cookeID, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:   "golang-event-planner.appspot.com",
				Value:  cookieID.String(),
				MaxAge: 0,
			}
		}
		http.SetCookie(w, cookie)
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
	// Return the struct with new user data.
	return "Registration of new user successful", &User{Username: username, Secret: secret, CookieSession: cookie.Value}, err
}
