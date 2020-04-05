package controllers

import (
	"api/auth"
	"api/database"
	"api/models"
	"api/repository"
	"api/repository/crud"
	"api/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersDB(db)

	func(usersRepository repository.UserRepository) {
		users, err := usersRepository.FindAll()
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusCreated, users)
	}(repo)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	// Get variables values
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersDB(db)

	func(usersRepository repository.UserRepository) {
		users, err := usersRepository.FindByID(uint32(uid))
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusCreated, users)
	}(repo)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// request body validation
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// model data validation
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersDB(db)

	func(usersRepository repository.UserRepository) {
		user, err = usersRepository.Save(user)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.RequestURI, user.ID))
		responses.JSON(w, http.StatusCreated, user)
	}(repo)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get variables values
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// request body validation
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// model data validation
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token := auth.ExtractToken(r)
	userId, err := auth.ExtractTokenID(token)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if userId != user.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersDB(db)

	func(usersRepository repository.UserRepository) {
		rows, err := usersRepository.Update(uint32(uid), user)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusOK, rows)
	}(repo)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	token := auth.ExtractToken(r)
	userId, err := auth.ExtractTokenID(token)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if userId != uint32(uuid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersDB(db)

	func(usersRepository repository.UserRepository) {
		rows, err := usersRepository.Delete(uuid)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Entity", fmt.Sprintf("%d%", uuid))
		responses.JSON(w, http.StatusNoContent, rows)
	}(repo)
}
