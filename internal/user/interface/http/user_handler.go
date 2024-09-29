package http

import (
	"encoding/json"
	"net/http"
	"time"
	"time-management/internal/shared/util"
	userCommand "time-management/internal/user/application/command"
	"time-management/internal/user/domain"
)

type UserHandler struct {
	LoginUserHandler userCommand.LoginUserHandler
}

func NewUserHandler(repository domain.UserRepository) *UserHandler {
	return &UserHandler{
		LoginUserHandler: userCommand.LoginUserHandler{Repo: repository},
	}
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	token, err := h.LoginUserHandler.Handle(userCommand.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return util.WriteJson(w, http.StatusUnauthorized, util.ApiError{Error: err.Error()})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    *token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	return util.WriteJson(w, http.StatusOK, nil)
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	return util.WriteJson(w, http.StatusOK, nil)
}
