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

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	decoder := json.NewDecoder(r.Body)
	var model UpdateUserInfoModel
	err = decoder.Decode(&model)

	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])
	err = UpdateUserInfo(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
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

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	decoder := json.NewDecoder(r.Body)
	var model UpdateUserPasswordModel
	err = decoder.Decode(&model)

	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])
	err = UpdateUserPassword(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	user, err := GetUserInfo(userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
	}
}

func FollowUserHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var model FollowUserModel
	model.FollowerID = authenticatedUserID
	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])

	err = FollowUser(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var model FollowUserModel
	model.FollowerID = authenticatedUserID
	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])

	err = UnfollowUser(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GetUserFollowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	users, err := GetUserFollows(userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		w.WriteHeader(http.StatusOK)
	}
}

func GetUserFollowersHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	users, err := GetUserFollowers(userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		w.WriteHeader(http.StatusOK)
	}
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

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	posts, err := GetUserPosts(userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
		w.WriteHeader(http.StatusOK)
	}
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	idParams := r.URL.Query()["id"]
	var ids []int
	for i := 0; i < len(idParams); i++ {
		id, _ := strconv.Atoi(idParams[i])
		ids = append(ids, id)
	}

	posts, err := GetPosts(authenticatedUserID, ids)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
		w.WriteHeader(http.StatusOK)
	}
}

func GetPostDetailHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	postID, err := strconv.Atoi(mux.Vars(r)["postId"])
	post, err := GetPostDetail(postID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
		w.WriteHeader(http.StatusOK)
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	postID, err := strconv.Atoi(mux.Vars(r)["postId"])
	err = DeletePost(postID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var model LikePostModel
	model.UserID = authenticatedUserID
	model.PostID, err = strconv.Atoi(mux.Vars(r)["postId"])

	err = LikePost(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func UnlikePostHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var model LikePostModel
	model.UserID = authenticatedUserID
	model.PostID, err = strconv.Atoi(mux.Vars(r)["postId"])

	err = UnlikePost(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
