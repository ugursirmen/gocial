package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"CreateUser",
		"POST",
		"/users",
		CreateUserHandler,
	},
	Route{
		"UpdateUserInfo",
		"PATCH",
		"/users/{userId}",
		UpdateUserInfoHandler,
	},
	// Route{
	// 	"UpdateUserPP",
	// 	"PATCH",
	// 	"/users/{userId}/pp",
	// 	UpdateUserPPHandler,
	// },
	// Route{
	// 	"UpdateUserPassword",
	// 	"PATCH",
	// 	"/users/{userId}/password",
	// 	UpdateUserPasswordHandler,
	// },
	// Route{
	// 	"GetUserInfo",
	// 	"GET",
	// 	"/users/{userId}",
	// 	GetUserInfoHandler,
	// },
	// Route{
	// 	"GetUserFollows",
	// 	"GET",
	// 	"/users/{userId}/follows",
	// 	GetUserFollowsHandler,
	// },
	// Route{
	// 	"GetUserFollowers",
	// 	"GET",
	// 	"/users/{userId}/followers",
	// 	GetUserFollowersHandler,
	// },
	// Route{
	// 	"GetUserPosts",
	// 	"GET",
	// 	"/users/{userId}/posts",
	// 	GetUserPostsHandler,
	// },
	// Route{
	// 	"CreatePost",
	// 	"POST",
	// 	"/posts",
	// 	CreatePostHandler,
	// },
	// Route{
	// 	"GetPosts",
	// 	"GET",
	// 	"/posts",
	// 	GetPostsHandler,
	// },
	// Route{
	// 	"GetPostDetail",
	// 	"GET",
	// 	"/posts/{postId}",
	// 	GetPostDetailHandler,
	// }
	// Route{
	// 	"DeletePost",
	// 	"DELETE",
	// 	"/posts/{postId}",
	// 	DeletePostHandler,
	// },
	// Route{
	// 	"LikePost",
	// 	"POST",
	// 	"/posts/{postId}/like",
	// 	LikePostHandler,
	// },
	// Route{
	// 	"UnLikePost",
	// 	"POST",
	// 	"/posts/{postId}/unlike",
	// 	UnLikePostHandler,
	// }
}
