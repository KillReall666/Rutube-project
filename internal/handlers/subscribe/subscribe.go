package subscribe

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/KillReall666/Rutube-project/internal/authentication"
	"github.com/KillReall666/Rutube-project/internal/handlers/validation"
	"github.com/KillReall666/Rutube-project/internal/logger"
	"github.com/KillReall666/Rutube-project/internal/model"
)

type SubscribeHandler struct {
	subscribe Subscriber
	Log       *logger.Logger
}

type Subscriber interface {
	UserInformationGetter(ctx context.Context, emailFromReq string) (*model.Employee, error)
	SubscriptionSetter(ctx context.Context, userID string, user model.Employee) error
}

func NewSubscribeHandler(subsc Subscriber, log *logger.Logger) *SubscribeHandler {
	return &SubscribeHandler{
		subscribe: subsc,
		Log:       log,
	}
}

func (s *SubscribeHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read data from request body", http.StatusBadRequest)
		return
	}

	emailOfTheEmployeeForSubscribe := string(body)
	err = validation.Email(emailOfTheEmployeeForSubscribe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userID, ok := authentication.GetUserIDFromCtx(r.Context())
	if !ok {
		s.Log.LogError("could not get caller from context")
	}

	user, err := s.subscribe.UserInformationGetter(context.Background(), emailOfTheEmployeeForSubscribe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = s.subscribe.SubscriptionSetter(context.Background(), userID, *user)
	if err != nil {

	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "You are subscribed!")
}
