package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	decoder := json.NewDecoder(r.Body)
	var model CreateUserModel
	err = decoder.Decode(&model)

	err = CreateUser(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func UpdateUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	decoder := json.NewDecoder(r.Body)
	var model UpdateUserInfoModel
	err = decoder.Decode(&model)

	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])
	err = UpdateUserInfo(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func UpdateUserPPHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var model UpdateUserPPModel
	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])

	pp, _, _ := r.FormFile("pp")

	if pp != nil {
		ppBytes, err := ioutil.ReadAll(pp)
		if err != nil {
			panic(err)
		}

		defer pp.Close()

		b64String := base64.StdEncoding.EncodeToString(ppBytes)
		model.ProfilePicture = b64String
	}

	err = UpdateUserPP(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	decoder := json.NewDecoder(r.Body)
	var model UpdateUserPasswordModel
	err = decoder.Decode(&model)

	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])
	err = UpdateUserPassword(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	user, err := GetUserInfo(userID)

	rw := &Response{w, err, &ResponseData{nil, nil, &user, nil}}
	rw.GenerateResponse()
}

func FollowUserHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var model FollowUserModel
	model.FollowerID = authenticatedUserID
	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])

	err = FollowUser(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var model FollowUserModel
	model.FollowerID = authenticatedUserID
	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])

	err = UnfollowUser(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func GetUserFollowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	users, err := GetUserFollows(userID)

	rw := &Response{w, err, &ResponseData{nil, nil, nil, &users}}
	rw.GenerateResponse()
}

func GetUserFollowersHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	users, err := GetUserFollowers(userID)

	rw := &Response{w, err, &ResponseData{nil, nil, nil, &users}}
	rw.GenerateResponse()
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var model CreatePostModel

	image, _, _ := r.FormFile("image")

	if image != nil {
		imageBytes, err := ioutil.ReadAll(image)
		if err != nil {
			panic(err)
		}

		defer image.Close()

		b64String := base64.StdEncoding.EncodeToString(imageBytes)
		model.Image = b64String
	}

	model.UserID, err = strconv.Atoi(r.FormValue("userId"))
	model.Description = r.FormValue("description")

	err = CreatePost(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	posts, err := GetUserPosts(userID)

	rw := &Response{w, err, &ResponseData{nil, &posts, nil, nil}}
	rw.GenerateResponse()
}

func GetPostsArbitraryHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	decoder := json.NewDecoder(r.Body)
	var model PostsArbitraryModel
	err = decoder.Decode(&model)

	posts, err := GetPostsArbitrary(model)

	rw := &Response{w, err, &ResponseData{nil, &posts, nil, nil}}
	rw.GenerateResponse()
}

func GetPostsArbitraryNonRelationalHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	decoder := json.NewDecoder(r.Body)
	var model PostsArbitraryModel
	err = decoder.Decode(&model)

	userEntity, err := GetUserInfo(model.UserID)

	postEntities, err := GetPostsByIds(model.PostIDs)

	var posts []PostDto

	if len(postEntities) > 0 {
		for _, postEntity := range postEntities {

			liked, _ := IsLiked(userEntity.ID, postEntity.ID)

			post := PostDto{
				ID:          postEntity.ID,
				Description: postEntity.Description,
				Image:       postEntity.Image,
				CreatedAt:   postEntity.CreatedAt,
				Liked:       liked,
			}

			followed, _ := IsFollowed(postEntity.UserID, model.UserID)

			post.Owner = &UserDto{
				ID:             userEntity.ID,
				Username:       userEntity.Username,
				FullName:       userEntity.FullName,
				ProfilePicture: userEntity.ProfilePicture,
				Followed:       followed,
			}

			posts = append(posts, post)
		}
	}

	rw := &Response{w, err, &ResponseData{nil, &posts, nil, nil}}
	rw.GenerateResponse()
}

func NewsfeedHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	posts, err := Newsfeed(authenticatedUserID)

	rw := &Response{w, err, &ResponseData{nil, &posts, nil, nil}}
	rw.GenerateResponse()
}

func GetPostDetailHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	postID, err := strconv.Atoi(mux.Vars(r)["postId"])
	post, err := GetPostDetail(postID)

	rw := &Response{w, err, &ResponseData{post, nil, nil, nil}}
	rw.GenerateResponse()
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	postID, err := strconv.Atoi(mux.Vars(r)["postId"])
	err = DeletePost(postID)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var model LikePostModel
	model.UserID = authenticatedUserID
	model.PostID, err = strconv.Atoi(mux.Vars(r)["postId"])

	err = LikePost(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}

func UnlikePostHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var model LikePostModel
	model.UserID = authenticatedUserID
	model.PostID, err = strconv.Atoi(mux.Vars(r)["postId"])

	err = UnlikePost(model)

	rw := &Response{w, err, nil}
	rw.GenerateResponse()
}
