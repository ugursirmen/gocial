package main

import "time"

type UserDto struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"fullName"`
	Followed       bool      `json:"followed"`
	CreatedAt      time.Time `json:"createdAt"`
	ProfilePicture string    `json:"pp"`
}

type PostDto struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Liked       bool      `json:"liked"`
	Owner       UserDto   `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
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
