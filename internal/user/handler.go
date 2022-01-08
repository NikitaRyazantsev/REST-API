package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"project/internal/middleware"
	"project/pkg/logging"
)

const (
	usersURL = "/users"
	userURL  = "/users/:id"
)

type Handler struct {
	Logger      *logging.Logger
	UserService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/create", middleware.PanicRecovery(middleware.Logging(h.CreateUser)))
	router.HandlerFunc(http.MethodGet, "/friends/:id", middleware.PanicRecovery(middleware.Logging(h.GetUserFriends)))
	router.HandlerFunc(http.MethodPut, userURL, middleware.PanicRecovery(middleware.Logging(h.UpdateUserAge)))
	router.HandlerFunc(http.MethodPost, "/make_friends", middleware.PanicRecovery(middleware.Logging(h.MakeFriends)))
	router.HandlerFunc(http.MethodDelete, usersURL, middleware.PanicRecovery(middleware.Logging(h.DeleteUser)))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Create user")
	w.Header().Set("Content-Type", "application/json")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()
	var u User
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	u.ID, err = h.UserService.Create(r.Context(), u)
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User with id:%v was created", u.ID)))
	return nil

}

func (h *Handler) GetUserFriends(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Find User's Friends")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get userID from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	friends, err := h.UserService.GetUserFriends(r.Context(), userID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("Друзья Пользователя c id:%s, это:%s", userID, friends)
	w.Write([]byte(message))
	return nil
}

func (h *Handler) UpdateUserAge(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Update user age")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get userID from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()
	type body struct {
		Age string `json:"age"`
	}
	var age body
	if err := json.Unmarshal(content, &age); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	err = h.UserService.UpdateAge(r.Context(), userID, age.Age)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("Пользователь c id:%s, обновлен", userID)
	w.Write([]byte(message))
	return nil
}

func (h *Handler) MakeFriends(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Make friends")
	w.Header().Set("Content-Type", "application/json")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()
	type body struct {
		SourceID string `json:"source_id"`
		TargetID string `json:"target_id"`
	}
	var message body
	if err := json.Unmarshal(content, &message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	var firstUser, secondUser User
	firstUser, secondUser, err = h.UserService.MakeFriends(r.Context(), message.SourceID, message.TargetID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	answer := fmt.Sprintf("Пользователи %s и %s теперь друзья", firstUser.Username, secondUser.Username)
	w.Write([]byte(answer))
	return nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE USER")
	w.Header().Set("Content-Type", "application/json")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()
	type body struct {
		TargetID string `json:"target_id"`
	}
	var message body
	if err := json.Unmarshal(content, &message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	err = h.UserService.Delete(r.Context(), message.TargetID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("Пользователь c id:%s - удален", message.TargetID)
	w.Write([]byte(text))

	return nil
}
