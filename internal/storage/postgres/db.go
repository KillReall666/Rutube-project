package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/KillReall666/Rutube-project/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

const createUsersTableQuery = `
      CREATE TABLE IF NOT EXISTS users (
    UserId VARCHAR(255) PRIMARY KEY,
    UserName VARCHAR(255) UNIQUE,
    Password VARCHAR(255),
	PhoneNumber VARCHAR(255) UNIQUE,
    DateOfBirthday VARCHAR(255),
    Email VARCHAR(255) UNIQUE,
    CONSTRAINT unique_person UNIQUE (Username)
);`

const createUserSubscriptionsTableQuery = `
	CREATE TABLE IF NOT EXISTS user_subscriptions (
	UserId VARCHAR(255),
	UserName VARCHAR(255),
	PhoneNumber VARCHAR(255),
	DateOfBirthDay VARCHAR(255),
	Email VARCHAR(255),
	FOREIGN KEY (UserId) REFERENCES users(UserId)
	    );
`

func New(connString string) (*Database, error) {
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUsersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating user table: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUserSubscriptionsTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating user subscriptions table: %v", err)
	}

	return &Database{db: conn}, nil
}

func (d *Database) UserSetter(ctx context.Context, user, password, id, phoneNumber, dateOfBirthday, email string) error {

	insertQuery := `
                INSERT INTO users (Username, Password, UserID, PhoneNumber, DateOfBirthday, Email)
				VALUES ($1, $2, $3, $4, $5, $6)
			
            `
	_, err := d.db.Exec(ctx, insertQuery, user, password, id, phoneNumber, dateOfBirthday, email)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(err.Error(), "unique_person") {
					return errors.New("this username already exists")
				} else if strings.Contains(err.Error(), "users_phonenumber_key") {
					return errors.New("this phone number already exists")
				} else if strings.Contains(err.Error(), "users_email_key") {
					return errors.New("this email already exists")
				}
			}
			return fmt.Errorf("error when trying to add user to database: %v", err)
		}
	}

	return nil
}

func (d *Database) CredentialsGetter(ctx context.Context, user string) (string, string, error) {
	var password, id string
	err := d.db.QueryRow(ctx, "SELECT password, userid FROM users WHERE username = $1", user).Scan(&password, &id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", errors.New("user not found")
		}

		return "", "", fmt.Errorf("error when getting hash password from database: %v", err)
	}

	return password, id, nil
}

func (d *Database) EmailGetter(ctx context.Context, userID string) (string, error) {
	var email string
	err := d.db.QueryRow(ctx, "SELECT email FROM users WHERE UserId = $1", userID).Scan(&email)
	if err != nil {
		return "", errors.New("user not found")
	}
	return email, nil
}

func (d *Database) UsersWithDataGetter(ctx context.Context) ([]model.Employee, error) {
	selectQuery := `
        SELECT  userName, phoneNumber, dateOfBirthday, email FROM users
    `

	var employees []model.Employee
	rows, err := d.db.Query(ctx, selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee model.Employee
		err = rows.Scan(&employee.UserName, &employee.PhoneNumber, &employee.DateOfBirth, &employee.Email)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil

}

func (d *Database) UserInformationGetter(ctx context.Context, emailFromReq string) (*model.Employee, error) {
	var user model.Employee
	err := d.db.QueryRow(ctx, "SELECT userName, phoneNumber, dateOfBirthday, email FROM users WHERE email = $1", emailFromReq).Scan(&user.UserName, &user.PhoneNumber, &user.DateOfBirth, &user.Email)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}

		return nil, fmt.Errorf("error when getting user information from database: %v", err)
	}
	return &user, nil
}

func (d *Database) SubscriptionSetter(ctx context.Context, userID string, user model.Employee) error {
	insertQuery := `
                INSERT INTO user_subscriptions (UserID, UserName, PhoneNumber, DateOfBirthday, Email)
				VALUES ($1, $2, $3, $4, $5)
			
            `

	_, err := d.db.Exec(ctx, insertQuery, userID, user.UserName, user.PhoneNumber, user.DateOfBirth, user.Email)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("this subscription already exists for you")
			} else {
				return fmt.Errorf("error when trying to add subscription to database: %v", err)
			}
		}
	}

	return nil
}

func (d *Database) UnSubscribe(ctx context.Context, email string) error {
	deleteQuery := `
        DELETE FROM user_subscriptions WHERE Email = $1
    `

	result, err := d.db.Exec(ctx, deleteQuery, email)
	if err != nil {
		return fmt.Errorf("err when unsubscribe: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("email not found")
	}

	return nil
}

func (d *Database) SelectPersonsWithBirthDay(ctx context.Context, dateOfBirthDay string) ([]model.Employee, error) {
	selectQuery := `
        SELECT userID, userName, phoneNumber, email FROM user_subscriptions WHERE DateOfBirthday = $1
    `

	var employees []model.Employee
	rows, err := d.db.Query(ctx, selectQuery, dateOfBirthDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee model.Employee
		err = rows.Scan(&employee.UserID, &employee.UserName, &employee.PhoneNumber, &employee.Email)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
