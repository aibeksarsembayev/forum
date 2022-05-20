package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"git.01.alem.school/quazar/forum/pkg/models"
	"git.01.alem.school/quazar/forum/pkg/session"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	p, err := app.posts.LatestPost()
	if err != nil {
		app.serverError(w, err)
		return
	}

	c, err := app.posts.AllCategory()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Count votes for post by range
	for _, post := range *p {
		likes, err := app.posts.CountVotes(post.PostID, true)
		if err != nil {
			app.serverError(w, err)
			return
		}
		dislikes, err := app.posts.CountVotes(post.PostID, false)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post.Votes.Likes = uint(likes)
		post.Votes.Dislikes = uint(dislikes)
	}

	username, IsSession := session.Get(r)

	// render page
	app.render(w, r, "home.page.html", &templateData{
		Posts:      *p,
		User:       models.User{Username: username},
		IsSession:  IsSession,
		Categories: *c,
	})
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.renderNoDB(w, r, "signin.page.html")
	} else if r.Method == "POST" {
		user := models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Confirm:  r.FormValue("confirm"),
			Created:  time.Now(),
		}

		// Get the existing entry present in the database for the given username
		u, err := app.posts.GetUser(user.Username)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if errors.Is(err, models.ErrNoRecord) {
				app.unauthorized(w)
			} else {
				// If the error is of any other type, send a 500 status
				app.serverError(w, err)
			}
			return
		}

		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			app.unauthorized(w)
			return
		}

		// Create a new random session token
		sessionToken := uuid.NewV4().String()

		// Set the token in the map, along with the user whom it represents
		session.Set(user.Username, sessionToken)

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds, the same as the cache
		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    url.QueryEscape(sessionToken),
			Expires:  time.Now().Add(120 * time.Second),
			HttpOnly: true, // Cookies provided only for HTTP(HTTPS) requests only
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusFound)

		// app.render(w, r, "home.page.html", &templateData{
		// 	User: u,
		// })
	}
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.renderNoDB(w, r, "signup.page.html")
	} else if r.Method == "POST" {
		user := models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Confirm:  r.FormValue("confirm"),
			Created:  time.Now(),
		}

		err := userSignUpForm(user)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		// fmt.Println(user)

		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing,
		// which we arbitrarily set as 8 (this value can be more or less,
		// depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

		// Convert time to string
		dateFormat := "2006-01-02 15:04:05 -0700 MST"
		strCreated := user.Created.Format(dateFormat)
		// fmt.Println(strCreated)

		// Next, insert the username, along with the hashed password into the database
		userID, err := app.posts.InsertUser(user.Username, string(hashedPassword), user.Email, strCreated)
		if err != nil {
			app.serverError(w, err)
			return
		}

		fmt.Printf("User ID: %d\nHashed Password: %s\n", userID, hashedPassword)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		app.methodNotAllowed(w)
		return
	}
}

func (app *application) signOut(w http.ResponseWriter, r *http.Request) {
	session.Clear(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	p, err := app.posts.GetPost(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Count votes for post
	likes, err := app.posts.CountVotes(p.PostID, true)
	if err != nil {
		app.serverError(w, err)
		return
	}
	dislikes, err := app.posts.CountVotes(p.PostID, false)
	if err != nil {
		app.serverError(w, err)
		return
	}

	v := models.VoteCount{
		Likes:    uint(likes),
		Dislikes: uint(dislikes),
	}
	p.Votes = &v

	// get comments
	c, err := app.posts.GetComments(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// to get comments number
	commentNumber := strconv.Itoa(len(*c))

	username, IsSession := session.Get(r)

	app.render(w, r, "show.page.html", &templateData{
		Post:          *p,
		User:          models.User{Username: username},
		IsSession:     IsSession,
		Comments:      *c,
		CommentNumber: commentNumber,
	})
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {

	username, IsSession := session.Get(r)
	if !IsSession {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	// get user
	user, err := app.posts.GetUser(username)
	if err != nil {
		app.badRequest(w)
		return
	}

	c, err := app.posts.AllCategory()
	if err != nil {
		app.badRequest(w)
		return
	}

	if r.Method == "GET" {
		app.render(w, r, "createpost.page.html", &templateData{
			User:       *user,
			Categories: *c,
			IsSession:  IsSession,
		})
	} else if r.Method == "POST" {
		// check for empty entry for title and/od content
		if (strings.TrimSpace(r.FormValue("title")) == "") || (strings.TrimSpace(r.FormValue("content")) == "") {
			app.badRequest(w)
			return
		}

		post := models.Post{
			Title:        r.FormValue("title"),
			Content:      r.FormValue("content"),
			UserID:       user.UserID,
			CategoryName: r.FormValue("category"),
			Created:      time.Now(),
		}

		// Convert time to string
		dateFormat := "2006-01-02 15:04:05 -0700 MST"
		strCreated := post.Created.Format(dateFormat)
		fmt.Println(strCreated)

		postID, err := app.posts.CreatePost(post.Title, post.Content, post.UserID, post.CategoryName, strCreated)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
	}
}

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {

	username, IsSession := session.Get(r)
	if !IsSession {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// get user model
	user, err := app.posts.GetUser(username)
	if err != nil {
		app.badRequest(w)
		return
	}

	if r.Method == "GET" {
		// ...
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	} else if r.Method == "POST" {

		cmnt := models.Comment{
			CommentBody: r.FormValue("comment"),
			PostID:      id,
			UserID:      user.UserID,
			Created:     time.Now(),
		}

		// Convert time to string
		dateFormat := "2006-01-02 15:04:05 -0700 MST"
		strCreated := cmnt.Created.Format(dateFormat)
		fmt.Println(strCreated)

		_, err := app.posts.CreateComment(cmnt.PostID, cmnt.UserID, cmnt.CommentBody, strCreated)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", cmnt.PostID), http.StatusSeeOther)
	}
}

func (app *application) likePost(w http.ResponseWriter, r *http.Request) {
	username, IsSession := session.Get(r)
	if !IsSession {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// get user model
	user, err := app.posts.GetUser(username)
	if err != nil {
		app.badRequest(w)
		return
	}

	if r.Method != "POST" {
		// ...
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	} else {
		vote := models.Vote{
			UserID: user.UserID,
			PostID: id,
			Value:  true,
		}

		// check if vote exists, then create or update
		_, err := app.posts.CreateVote(vote.PostID, vote.UserID, vote.Value)

		if err != nil {
			app.serverError(w, err)
			return
		}

		path := r.FormValue("path")

		http.Redirect(w, r, path, http.StatusFound)
	}
}

// need to combine with likePost func to optimize
func (app *application) dislikePost(w http.ResponseWriter, r *http.Request) {
	username, IsSession := session.Get(r)
	if !IsSession {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// get user model
	user, err := app.posts.GetUser(username)
	if err != nil {
		app.badRequest(w)
		return
	}

	if r.Method != "POST" {
		// ...
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	} else {
		vote := models.Vote{
			UserID: user.UserID,
			PostID: id,
			Value:  false,
		}

		_, err := app.posts.CreateVote(vote.PostID, vote.UserID, vote.Value)

		if err != nil {
			app.serverError(w, err)
			return
		}

		path := r.FormValue("path")

		http.Redirect(w, r, path, http.StatusFound)
	}
}

func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	username, IsSession := session.Get(r)
	if !IsSession {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	comment_id, err := strconv.Atoi(r.URL.Query().Get("comment"))
	if err != nil || comment_id < 1 {
		app.notFound(w)
		return
	}

	//get user model
	user, err := app.posts.GetUser(username)
	if err != nil {
		app.badRequest(w)
		return
	}

	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	} else {
		voteComment := models.VoteComment{
			UserID:    user.UserID,
			PostID:    id,
			CommentID: comment_id,
			Value:     true,
		}

		_, err := app.posts.CreateVoteComment(voteComment.PostID, voteComment.UserID, voteComment.CommentID, voteComment.Value)

		if err != nil {
			app.serverError(w, err)
			return
		}
		path := r.FormValue("path")

		http.Redirect(w, r, path, http.StatusFound)
	}

}

func (app *application) dislikeComment(w http.ResponseWriter, r *http.Request) {
	username, IsSession := session.Get(r)
	if !IsSession {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	comment_id, err := strconv.Atoi(r.URL.Query().Get("comment"))
	if err != nil || comment_id < 1 {
		app.notFound(w)
		return
	}

	//get user model
	user, err := app.posts.GetUser(username)
	if err != nil {
		app.badRequest(w)
		return
	}

	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	} else {
		voteComment := models.VoteComment{
			UserID:    user.UserID,
			PostID:    id,
			CommentID: comment_id,
			Value:     false,
		}

		_, err := app.posts.CreateVoteComment(voteComment.PostID, voteComment.UserID, voteComment.CommentID, voteComment.Value)

		if err != nil {
			app.serverError(w, err)
			return
		}
		path := r.FormValue("path")

		http.Redirect(w, r, path, http.StatusFound)
	}

}

func (app *application) showProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")

	u, err := app.posts.GetUser(username)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	id := u.UserID

	p, err := app.posts.GetPostsByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Count votes for post by range
	for _, post := range *p {
		likes, err := app.posts.CountVotes(post.PostID, true)
		if err != nil {
			app.serverError(w, err)
			return
		}
		dislikes, err := app.posts.CountVotes(post.PostID, false)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post.Votes.Likes = uint(likes)
		post.Votes.Dislikes = uint(dislikes)
	}

	_, IsSession := session.Get(r)

	app.render(w, r, "profile.page.html", &templateData{
		Posts:     *p,
		User:      *u,
		IsSession: IsSession,
	})
}

func (app *application) showPostByCategory(w http.ResponseWriter, r *http.Request) {

	category := r.URL.Query().Get("category")

	p, err := app.posts.GetPostbyCat(category)
	if err != nil {
		app.serverError(w, err)
		return
	}

	c, err := app.posts.AllCategory()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Count votes for post by range
	for _, post := range *p {
		likes, err := app.posts.CountVotes(post.PostID, true)
		if err != nil {
			app.serverError(w, err)
			return
		}
		dislikes, err := app.posts.CountVotes(post.PostID, false)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post.Votes.Likes = uint(likes)
		post.Votes.Dislikes = uint(dislikes)
	}

	username, IsSession := session.Get(r)

	app.render(w, r, "category.page.html", &templateData{
		Posts:      *p,
		User:       models.User{Username: username},
		IsSession:  IsSession,
		Categories: *c,
		Category:   models.Category{CategoryName: category},
	})
}

func (app *application) showLikedPosts(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")

	u, err := app.posts.GetUser(username)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	id := u.UserID

	p, err := app.posts.GetPostbyUserVote(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Count votes for post by range
	for _, post := range *p {
		likes, err := app.posts.CountVotes(post.PostID, true)
		if err != nil {
			app.serverError(w, err)
			return
		}
		dislikes, err := app.posts.CountVotes(post.PostID, false)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post.Votes.Likes = uint(likes)
		post.Votes.Dislikes = uint(dislikes)
	}

	_, IsSession := session.Get(r)

	app.render(w, r, "profile.page.html", &templateData{
		Posts:     *p,
		User:      *u,
		IsSession: IsSession,
	})
}
