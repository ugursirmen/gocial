package main

import (
	"encoding/json"
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
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	decoder := json.NewDecoder(r.Body)
	var model UpdateUserInfoModel
	err = decoder.Decode(&model)

	model.UserId, err = strconv.Atoi(mux.Vars(r)["userId"])
	err = UpdateUserInfo(model)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
