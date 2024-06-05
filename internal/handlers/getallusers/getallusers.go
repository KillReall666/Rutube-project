package getallusers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/KillReall666/Rutube-project/internal/authentication"
	"github.com/KillReall666/Rutube-project/internal/logger"
	"github.com/KillReall666/Rutube-project/internal/model"
)

type GetAllUsersHandler struct {
	ug  UsersGetter
	Log *logger.Logger
}

type UsersGetter interface {
	UsersWithDataGetter(ctx context.Context) ([]model.Employee, error)
}

func NewGetAllUsersHandler(uG UsersGetter, log *logger.Logger) *GetAllUsersHandler {
	return &GetAllUsersHandler{
		ug:  uG,
		Log: log,
	}
}

func (g *GetAllUsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	_, ok := authentication.GetUserIDFromCtx(r.Context())
	if !ok {
		g.Log.LogError("could not get caller from context")
	}

	users, err := g.ug.UsersWithDataGetter(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
