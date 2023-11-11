package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/golang/mock/mockgen/model"

	"lab1/internal/model"
)

//go:generate mockgen -destination=./mocks/mock_CreateUserUsecase.go -package=mocks lab1/internal/handler CreateUserUsecase
type CreateUserUsecase interface {
	CreateUser(user model.User) error
}

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user model.User
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't decode user"))
		return
	}

	err = h.createUserUsecase.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
