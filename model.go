package main

import "time"

type UserModel struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Followed  bool      `json:"followed"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostModel struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Liked       bool      `json:"liked"`
	Owner       UserModel `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateUserModel struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInfoModel struct {
	UserId    int    `json:"userId"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
}
