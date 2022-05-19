package session

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var sessions sync.Map

func Set(username, sessionToken string) {

	// check for existing the same user session
	sessions.Range(func(key, value interface{}) bool {
		log.Println("sessions: ", key, value)
		if username == value {
			sessions.Delete(key)
		}
		return true
	})

	sessions.Store(sessionToken, username)
}

// IsSession to define user is on or off
// func IsSession(r *http.Request) bool {
// 	cookie, err := r.Cookie("session_token")
// 	if err != nil {
// 		return false
// 	}

// 	if _, ok := sessions.Load(cookie.Value); ok {
// 		return true
// 	}

// 	return false
// }

// Get user name from syncmap via token
func Get(r *http.Request) (string, bool) {
	var username string
	cookie, err := r.Cookie("session_token")
	if err == nil {
		if value, ok := sessions.Load(cookie.Value); ok {
			username = fmt.Sprint(value)
			return username, true
		}
	}
	return username, false
}

// Clear session
func Clear(w http.ResponseWriter, r *http.Request) {
	if cookieValue, err := getValue(r); err == nil {
		sessions.Delete(cookieValue)
	}
	cookie := &http.Cookie{
		Name:   "your-name",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func getValue(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if cookie == nil || err != nil {
		return "", err
	}

	cookieValue := cookie.Value

	return cookieValue, err
}
