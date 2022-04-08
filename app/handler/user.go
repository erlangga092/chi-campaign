package handler

import (
	"chi-app/app/auth"
	"chi-app/app/helper"
	"chi-app/app/key"
	"chi-app/app/user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		data := "Content Type must be application/json"
		response := helper.APIResponse("Failed register user", http.StatusBadRequest, "error", data)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	// https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
	// reference validate struct fields
	v := validator.New()
	input := user.RegisterUserInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Failed register user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Failed register user", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Failed register user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Failed register user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := user.FormatUser(newUser, token)
	response := helper.APIResponse("Account has been created", http.StatusCreated, "success", formatter)
	helper.JSON(w, response, http.StatusCreated)
}

func (h *userHandler) CheckEmailAvailable(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		data := "Content Type must be application/json"
		response := helper.APIResponse("Failed check email", http.StatusBadRequest, "error", data)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	v := validator.New()
	input := user.CheckEmailAvailableInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Failed check email", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Failed check email", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	isAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		response := helper.APIResponse("Failed check email", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	available := false
	if isAvailable {
		available = true
	}

	data := map[string]interface{}{
		"is_available": available,
	}

	response := helper.APIResponse("Success check available email", http.StatusOK, "success", data)
	helper.JSON(w, response, http.StatusOK)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		data := "Content Type must be application/json"
		response := helper.APIResponse("Failed login user", http.StatusBadRequest, "error", data)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	v := validator.New()
	input := user.LoginUserInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Failed login user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Failed login user", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	loggedInUser, err := h.userService.LoginUser(input)
	if err != nil {
		response := helper.APIResponse("Failed login user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.APIResponse("Failed login user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Login Successfully", http.StatusCreated, "success", formatter)
	helper.JSON(w, response, http.StatusCreated)
}

func (h *userHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := r.FormValue("alias")
	uploadedFile, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get user data from middleware
	user := r.Context().Value(key.CtxKeyAuth{}).(user.User)
	filename := fmt.Sprintf("%d-%s", user.ID, handler.Filename)

	if alias != "" {
		filename = fmt.Sprintf("%d-%s%s", user.ID, alias, filepath.Ext(handler.Filename))
	}

	fileLocation := filepath.Join(dir, "images", filename)

	// update avatar image to database
	// if error when update to database, cancel upload to local directory
	_, err = h.userService.UploadAvatar(user.ID, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		data := map[string]interface{}{
			"is_uploaded": false,
		}

		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	defer targetFile.Close()

	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Success upload to local directory!")

	data := map[string]interface{}{
		"is_uploaded": true,
	}

	response := helper.APIResponse("Avatar successfully uploaded!", http.StatusCreated, "success", data)
	helper.JSON(w, response, http.StatusCreated)
}
