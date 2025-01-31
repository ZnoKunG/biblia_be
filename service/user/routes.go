package user

import (
	"fmt"
	"net/http"

	"github.com/ZnoKunG/biblia_be.git/service/auth"
	"github.com/ZnoKunG/biblia_be.git/types"
	"github.com/ZnoKunG/biblia_be.git/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload from the request
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload &v", errors))
		return
	}

	// check if user already exists
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashed_password, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create new user if not exists
	err = h.store.CreateUser(types.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashed_password,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
