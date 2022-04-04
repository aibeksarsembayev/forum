package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/signin", app.signIn)
	mux.HandleFunc("/signup", app.signUp)
	mux.HandleFunc("/signout", app.signOut)
	mux.HandleFunc("/post", app.showPost)
	mux.HandleFunc("/post/create", app.createPost)
	mux.HandleFunc("/comment/create", app.createComment)
	mux.HandleFunc("/post/like", app.likePost)
	mux.HandleFunc("/post/dislike", app.dislikePost)
	mux.HandleFunc("/post/comment/like", app.likeComment)
	mux.HandleFunc("/post/comment/dislike", app.dislikeComment)
	mux.HandleFunc("/profile", app.showProfile)
	mux.HandleFunc("/profile/liked", app.showLikedPosts)
	mux.HandleFunc("/category", app.showPostByCategory)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
