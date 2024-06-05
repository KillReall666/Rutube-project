package credentials

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username       string `json:"login"`
	PasswordHash   string `json:"password"`
	PhoneNumber    string `json:"phone_number"`
	DateOfBirthday string `json:"date_of_birthday"`
	Email          string `json:"email"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) ComparePassword(hashedPasswordFromDB, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(password))
	return err == nil
}
