package interrogator

import (
	"context"
	"fmt"
	"time"

	"github.com/KillReall666/Rutube-project/internal/logger"
	"github.com/KillReall666/Rutube-project/internal/model"
	"github.com/KillReall666/Rutube-project/internal/storage/postgres"
)

type Interrogator struct {
	db  *postgres.Database
	log *logger.Logger
}

func NewInterrogator(db *postgres.Database, log *logger.Logger) *Interrogator {
	return &Interrogator{
		db:  db,
		log: log,
	}
}

func (i *Interrogator) BirthDaysFinder() {
	for {
		now := time.Now()
		nextDay := now.AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		duration := nextDay.Sub(now) //для теста надо закоментить, а то уснёт на сутки.
		time.Sleep(duration)         //это тоже.

		formattedDate := now.Format("02.01.2006")
		users, err := i.db.SelectPersonsWithBirthDay(context.Background(), formattedDate)
		if err != nil {
			i.log.LogInfo("err when select persons with birthday", err)
		}

		for _, user := range users {
			usersEmail, err := i.db.EmailGetter(context.Background(), user.UserID)
			if err != nil {
				i.log.LogError("err when get email for sending", err)
			}
			i.CongratulationsSender(usersEmail, user)
		}
	}
}

// CongratulationsSender Заглуша имитирующая отправку оповещения на почту.
func (i *Interrogator) CongratulationsSender(usersEmail string, user model.Employee) {
	fmt.Printf("Reminder have been sent on %v. Data of the user who celebrates his birthday today. Name: %v, Mail: %v, Phone: %v.", usersEmail, user.UserName, user.Email, user.PhoneNumber)
}
