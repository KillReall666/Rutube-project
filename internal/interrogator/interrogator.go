package interrogator

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
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

func (i *Interrogator) EmailSender() { //usersEmail string, user model.Employee
	from := mail.Address{"", ""} //user.Email
	to := mail.Address{"", ""}
	subj := "Email subject"
	body := "This is test email! \n Hail to the King!"

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	serverName := "smtp.mail.ru:465"

	host, _, _ := net.SplitHostPort(serverName)

	auth := smtp.PlainAuth("", "hicobra@mail.ru", "", host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", serverName, tlsconfig)
	if err != nil {
		i.log.LogError("err when make tls dial", err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		i.log.LogError("err when connect to server", err)
	}

	if err = c.Auth(auth); err != nil {
		i.log.LogError("err when authenticate", err)
	}

	if err = c.Mail(from.Address); err != nil {
		i.log.LogError("err when set from", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		i.log.LogError("err when set to", err)
	}

	w, err := c.Data()
	if err != nil {
		i.log.LogError("err when set data", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		i.log.LogError("err when write data", err)
	}

	err = w.Close()
	if err != nil {
		i.log.LogError("err when close", err)
	}

	i.log.LogInfo("Email send!")

	c.Quit()
}
