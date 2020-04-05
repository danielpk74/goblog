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

func CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// model validation
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.Prepare()
	err = post.Validate("")
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

	if userId != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryPostsDB(db)

	func(postsRepository repository.PostRepository) {
		post, err = postsRepository.Save(post)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		// w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.RequestURI, post.ID))
		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, post.ID))
		responses.JSON(w, http.StatusCreated, post)
	}(repo)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
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

	repo := crud.NewRepositoryPostsDB(db)

	func(postsRepository repository.PostRepository) {
		users, err := postsRepository.FindByID(uint32(uid))
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusCreated, users)
	}(repo)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryPostsDB(db)

	func(postsRepository repository.PostRepository) {
		posts, err := postsRepository.FindAll()
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusCreated, posts)
	}(repo)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.Prepare()
	err = post.Validate("")
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

	if userId != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryPostsDB(db)
	func(postsRepository repository.PostRepository) {
		rows, err := postsRepository.Update(uint32(uid), post)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusOK, rows)
	}(repo)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId, err := strconv.ParseInt(vars["id"], 10, 32)
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

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryPostsDB(db)

	func(postsRepository repository.PostRepository) {
		rows, err := postsRepository.Delete(postId, userId)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Entity", fmt.Sprintf("%d%", postId))
		responses.JSON(w, http.StatusNoContent, rows)
	}(repo)
}
