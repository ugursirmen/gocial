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
	Route{
		"UpdateUserPP",
		"PATCH",
		"/users/{userId}/pp",
		UpdateUserPPHandler,
	},
	Route{
		"UpdateUserPassword",
		"PATCH",
		"/users/{userId}/password",
		UpdateUserPasswordHandler,
	},
	Route{
		"GetUserInfo",
		"GET",
		"/users/{userId}",
		GetUserInfoHandler,
	},
	Route{
		"FollowUser",
		"POST",
		"/users/{userId}/follow",
		FollowUserHandler,
	},
	Route{
		"UnfollowUser",
		"POST",
		"/users/{userId}/unfollow",
		UnfollowUserHandler,
	},
	Route{
		"GetUserFollows",
		"GET",
		"/users/{userId}/follows",
		GetUserFollowsHandler,
	},
	Route{
		"GetUserFollowers",
		"GET",
		"/users/{userId}/followers",
		GetUserFollowersHandler,
	},
	Route{
		"CreatePost",
		"POST",
		"/posts",
		CreatePostHandler,
	},
	Route{
		"GetUserPosts",
		"GET",
		"/users/{userId}/posts",
		GetUserPostsHandler,
	},
	Route{
		"Newsfeed",
		"GET",
		"/newsfeed",
		NewsfeedHandler,
	},
	Route{
		"GetPostsArbitrary",
		"POST",
		"/posts/arbitrary",
		GetPostsArbitraryHandler,
	},
	Route{
		"GetPostsArbitrary",
		"POST",
		"/posts/arbitrary/non-relational",
		GetPostsArbitraryNonRelationalHandler,
	},
	Route{
		"GetPostDetail",
		"GET",
		"/posts/{postId}",
		GetPostDetailHandler,
	},
	Route{
		"DeletePost",
		"DELETE",
		"/posts/{postId}",
		DeletePostHandler,
	},
	Route{
		"LikePost",
		"POST",
		"/posts/{postId}/like",
		LikePostHandler,
	},
	Route{
		"UnLikePost",
		"POST",
		"/posts/{postId}/unlike",
		UnlikePostHandler,
	},
}
