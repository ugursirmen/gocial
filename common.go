package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {
	http.ResponseWriter
	err  error
	data *ResponseData
}

type ResponseData struct {
	post  *PostDto
	posts *[]PostDto
	user  *UserDto
	users *[]UserDto
}

func (res *Response) GenerateResponse() {

	if res.err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(res.err.Error()))
	} else {
		res.Header().Set("Content-Type", "application/json")

		if res.data != nil {
			if res.data.post != nil {
				json.NewEncoder(res).Encode(res.data.post)
			} else if res.data.posts != nil {
				json.NewEncoder(res).Encode(res.data.posts)
			} else if res.data.user != nil {
				json.NewEncoder(res).Encode(res.data.user)
			} else if res.data.users != nil {
				json.NewEncoder(res).Encode(res.data.users)
			}
		}

		res.WriteHeader(http.StatusOK)
	}
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

type UserDto struct {
	ID             int        `json:"id"`
	Username       string     `json:"username"`
	FullName       string     `json:"fullName"`
	Followed       bool       `json:"followed"`
	CreatedAt      *time.Time `json:"createdAt"`
	ProfilePicture string     `json:"pp"`
}

type PostDto struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Liked       bool      `json:"liked"`
	CreatedAt   time.Time `json:"createdAt"`
	Owner       *UserDto  `json:"owner,omitempty"`
	Image       string    `json:"image"`
}

type PostEntity struct {
	ID          int
	Description string
	Image       string
	Deleted     bool
	CreatedAt   time.Time
	UserID      int
}

type UserEntity struct {
	ID             int
	Username       string
	FirstName      string
	LastName       string
	Bio            string
	ProfilePicture string
}

type CreateUserModel struct {
	Username string
	Email    string
	Password string
}

type UpdateUserInfoModel struct {
	UserID    int
	Username  string
	FirstName string
	LastName  string
	Bio       string
}

type UpdateUserPPModel struct {
	UserID         int
	ProfilePicture string
}

type UpdateUserPasswordModel struct {
	UserID      int
	OldPassword string
	NewPassword string
}

type FollowUserModel struct {
	UserID     int
	FollowerID int
}

type CreatePostModel struct {
	UserID      int
	Description string
	Image       string
}

type LikePostModel struct {
	UserID int
	PostID int
}

type PostsArbitraryModel struct {
	UserID  int
	PostIDs []int
}
