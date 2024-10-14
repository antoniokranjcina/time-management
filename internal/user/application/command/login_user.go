package command

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"os"
	"time"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
)

const JwtId = "id"
const JwtRole = "role"
const JwtExp = "exp"
const JwtExpirationTime = time.Hour * 24

type LoginUserCommand struct {
	Email    string
	Password string
}

type LoginUserHandler struct {
	Repo domain.UserRepository
}

func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUserCommand) (*string, error) {
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return nil, sharedUtil.NewValidationError(domain.ErrEmailWrongFormat)
	}

	user, err := h.Repo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, domain.ErrInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.Password))
	if err != nil {
		return nil, domain.ErrInvalidEmailOrPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JwtId:   user.Id,
		JwtRole: user.Role,
		JwtExp:  time.Now().Add(JwtExpirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
