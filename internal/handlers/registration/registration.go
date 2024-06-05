package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/KillReall666/Rutube-project/internal/authentication"
	"github.com/KillReall666/Rutube-project/internal/credentials"
	"github.com/KillReall666/Rutube-project/internal/handlers/validation"
	"github.com/KillReall666/Rutube-project/internal/logger"
	"github.com/KillReall666/Rutube-project/internal/storage/redis"
)

type RegisterHandler struct {
	setUser     UserSetter
	redisClient *redis.RedisClient
	log         *logger.Logger
}

type UserSetter interface {
	UserSetter(ctx context.Context, user, password, id, phoneNumber, dateOfBirthday, email string) error
}

func NewRegistrationHandler(us UserSetter, redis *redis.RedisClient, log *logger.Logger) *RegisterHandler {
	return &RegisterHandler{
		setUser:     us,
		redisClient: redis,
		log:         log,
	}
}

func (reg *RegisterHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	var user credentials.User
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validation.Email(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = validation.PhoneNumber(user.PhoneNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = user.SetPassword(user.PasswordHash)
	if err != nil {
		reg.log.LogError("error while hashing password", err)
	}

	var token string
	idNew := uuid.New()
	idString := idNew.String()

	err = reg.setUser.UserSetter(ctx, user.Username, user.PasswordHash, idString, user.PhoneNumber, user.DateOfBirthday, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	token, err = authentication.BuildJWTString(idString)
	if err != nil {
		reg.log.LogError("err when get JWT token while registration:", err)
	}

	err = reg.redisClient.Set(idString, token)
	if err != nil {
		reg.log.LogError("err when set value to redis in auth handler:", err)
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You have successfully registered and authorized")

	reg.log.LogInfo("user", idString, "successfully registered and authorized")
}
