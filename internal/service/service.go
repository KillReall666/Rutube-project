package service

import (
	"context"
	"fmt"

	"github.com/KillReall666/Rutube-project/internal/config"
	"github.com/KillReall666/Rutube-project/internal/logger"
	"github.com/KillReall666/Rutube-project/internal/model"
	"github.com/KillReall666/Rutube-project/internal/storage/postgres"
)

type service struct {
	cfg *config.Config
	log *logger.Logger
	db  *postgres.Database
}

type Service interface {
	UserSetter(ctx context.Context, user, password, id, phoneNumber, dateOfBirthday, email string) error
	CredentialsGetter(ctx context.Context, user string) (string, string, error)
}

func New(cfg *config.Config, log *logger.Logger, db *postgres.Database) *service {
	return &service{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (s *service) Stub() {
	fmt.Println("stub")
}

func (s *service) UserSetter(ctx context.Context, user, password, id, phoneNumber, dateOfBirthDay, email string) error {
	err := s.db.UserSetter(ctx, user, password, id, phoneNumber, dateOfBirthDay, email)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CredentialsGetter(ctx context.Context, user string) (string, string, error) {
	hashPassword, id, err := s.db.CredentialsGetter(ctx, user)
	if err != nil {
		return "", "", err
	}

	return hashPassword, id, err
}

func (s *service) UserInformationGetter(ctx context.Context, email string) (*model.Employee, error) {
	user, err := s.db.UserInformationGetter(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) SubscriptionSetter(ctx context.Context, userID string, user model.Employee) error {
	err := s.db.SubscriptionSetter(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UnSubscribe(ctx context.Context, email string) error {
	err := s.db.UnSubscribe(ctx, email)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UsersWithDataGetter(ctx context.Context) ([]model.Employee, error) {
	user, err := s.db.UsersWithDataGetter(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
