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

	decoder := json.NewDecoder(r.Body)
	var model FollowUserModel
	err = decoder.Decode(&model)

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

	decoder := json.NewDecoder(r.Body)
	var model FollowUserModel
	err = decoder.Decode(&model)

	model.UserID, err = strconv.Atoi(mux.Vars(r)["userId"])
	err = UnfollowUser(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
