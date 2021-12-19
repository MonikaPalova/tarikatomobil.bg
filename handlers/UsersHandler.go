package handlers

import (
	"encoding/json"
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"log"
	"net/http"
)

type UsersHandler struct {
	DB Database
}

func (u UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	users, err := u.DB.GetUsers()
	if err != nil {
		log.Printf("Could not fetch users from DB: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("Could not marshal users: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(bytes)
}

func (u UsersHandler) Post(w http.ResponseWriter, r *http.Request) {
	var userToCreate model.User
	if err := json.NewDecoder(r.Body).Decode(&userToCreate); err != nil {
		log.Printf("Could not parse request body as JSON: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := u.DB.CreateUser(userToCreate); err != nil {
		log.Printf("Could not create user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(&userToCreate); err != nil {
		log.Printf("Could not marshal created user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
