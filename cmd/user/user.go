package user

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/vickon16/third-backend-tutorial/cmd/middleware"
	"github.com/vickon16/third-backend-tutorial/cmd/sqlc"
	"github.com/vickon16/third-backend-tutorial/cmd/types"
	"github.com/vickon16/third-backend-tutorial/cmd/utils"
)

type Handler struct {
	db *sqlc.Queries
}

func NewHandler(db *sqlc.Queries) *Handler {
	return &Handler{db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.Handle("/register",
		alice.New(middleware.Log).ThenFunc(h.handleRegister)).Methods(http.MethodPost)
	router.Handle("/login",
		alice.New(middleware.Log).ThenFunc(h.handleLogin)).Methods(http.MethodPost)

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJsonAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if user name exists
	user, err := h.db.GetUserByEmail(r.Context(), payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email already exists"))
		return
	}

	// hash the user password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.db.CreateUser(r.Context(), sqlc.CreateUserParams{
		ID:             uuid.New().String(),
		Name:           payload.Name,
		Email:          payload.Email,
		Hashedpassword: hashedPassword,
	}); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, err := utils.CreateJWT(user.ID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"userId":      user.ID,
		"message":     "User Created Successfully",
		"accessToken": token,
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJsonAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if user name exists
	user, err := h.db.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("e: Invalid user email or password"))
		return
	}

	if !utils.ComparePassword(user.Hashedpassword, payload.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("p: Invalid user email or password"))
		return
	}

	token, err := utils.CreateJWT(user.ID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"userId":      user.ID,
		"message":     "Successfully Logged In",
		"accessToken": token,
	})
}
