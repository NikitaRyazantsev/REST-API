package user

// file for routing and handle

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"project/internal/middleware"
	"project/pkg/logging"
)

// constants for standart path
const (
	usersURL = "/users"
	userURL  = "/users/:id"
)

// Handler - structure for handlers with logger
type Handler struct {
	Logger      *logging.Logger
	UserService Service
}

// Register - func for init routs
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/create", middleware.PanicRecovery(middleware.Logging(h.CreateUser)))
	router.HandlerFunc(http.MethodGet, "/friends/:id", middleware.PanicRecovery(middleware.Logging(h.GetUserFriends)))
	router.HandlerFunc(http.MethodPut, userURL, middleware.PanicRecovery(middleware.Logging(h.UpdateUserAge)))
	router.HandlerFunc(http.MethodPost, "/make_friends", middleware.PanicRecovery(middleware.Logging(h.MakeFriends)))
	router.HandlerFunc(http.MethodDelete, usersURL, middleware.PanicRecovery(middleware.Logging(h.DeleteUser)))
}

// CreateUser - creating user by http-request
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Create user")
	w.Header().Set("Content-Type", "application/json")

	// get data from http`s body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()
	var u User

	// unmarshalling json to user struct
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	// call user-service for create user in database
	u.ID, err = h.UserService.Create(r.Context(), u)
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("From 1 instance: User with id:%v was created", u.ID)))
	return nil

}

// GetUserFriends - getting friends from one user by http-request
func (h *Handler) GetUserFriends(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Find User's Friends")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get userID from context")

	// getting id from url params
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	// call user-service for getting user`s friends from database
	friends, err := h.UserService.GetUserFriends(r.Context(), userID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)

	// creating message for answer
	message := fmt.Sprintf("From 1 instance: Друзья Пользователя c id:%s, это:%s", userID, friends)
	w.Write([]byte(message))
	return nil
}

// UpdateUserAge - func for update user`s age in database by http-request
func (h *Handler) UpdateUserAge(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Update user age")
	w.Header().Set("Content-Type", "application/json")

	// getting id from url params
	h.Logger.Debug("get userID from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	// getting new age of user from http message body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	// struct for unmarshalling data from json
	type body struct {
		Age string `json:"age"`
	}
	var age body

	// unmarshalling data from json message
	if err := json.Unmarshal(content, &age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	// call user-service for change age of user in database
	err = h.UserService.UpdateAge(r.Context(), userID, age.Age)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)

	// create message for answer
	message := fmt.Sprintf("From 1 instance: Пользователь c id:%s, обновлен", userID)
	w.Write([]byte(message))
	return nil
}

// MakeFriends - func that making friends two users, data getting from http-request
func (h *Handler) MakeFriends(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Make friends")
	w.Header().Set("Content-Type", "application/json")

	// getting data from http`s body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	// creating struct for unmarshalling data
	type body struct {
		SourceID string `json:"source_id"`
		TargetID string `json:"target_id"`
	}
	var message body

	// unmarshalling data from json
	if err := json.Unmarshal(content, &message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	var firstUser, secondUser User

	// calling user-service to make friends from to users in database
	firstUser, secondUser, err = h.UserService.MakeFriends(r.Context(), message.SourceID, message.TargetID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)

	// creating message for answer
	answer := fmt.Sprintf("From 1 instance: Пользователи %s и %s теперь друзья", firstUser.Username, secondUser.Username)
	w.Write([]byte(answer))
	return nil
}

// DeleteUser - func that delete user from database and delete from friends arrays of all users
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE USER")
	w.Header().Set("Content-Type", "application/json")

	// getting data from http-request body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	// create struct for unmarshalling
	type body struct {
		TargetID string `json:"target_id"`
	}
	var message body

	// unmarshall data from json
	if err := json.Unmarshal(content, &message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	// call user-service to delete user from database
	err = h.UserService.Delete(r.Context(), message.TargetID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("From 1 instance: Пользователь c id:%s - удален", message.TargetID)
	w.Write([]byte(text))

	return nil
}
