package unsubscribe

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/KillReall666/Rutube-project/internal/authentication"
	"github.com/KillReall666/Rutube-project/internal/handlers/validation"
	"github.com/KillReall666/Rutube-project/internal/logger"
)

type UnSubscribeHandler struct {
	unsubscribe UnSubscriber
	Log         *logger.Logger
}

type UnSubscriber interface {
	UnSubscribe(ctx context.Context, email string) error
}

func NewUnSubscribeHandler(unSubsc UnSubscriber, log *logger.Logger) *UnSubscribeHandler {
	return &UnSubscribeHandler{
		unsubscribe: unSubsc,
		Log:         log,
	}
}

func (s *UnSubscribeHandler) UnSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read data from request body", http.StatusBadRequest)
		return
	}

	emailOfTheEmployeeForUnSubscribe := string(body)
	err = validation.Email(emailOfTheEmployeeForUnSubscribe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	_, ok := authentication.GetUserIDFromCtx(r.Context())
	if !ok {
		s.Log.LogError("could not get caller from context")
	}

	err = s.unsubscribe.UnSubscribe(context.Background(), emailOfTheEmployeeForUnSubscribe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "You are unsubscribed!")

}
