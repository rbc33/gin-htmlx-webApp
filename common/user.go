package common

import (
	"html"
	"strings"

	"github.com/rbc33/gocms/utils/token"
	"golang.org/x/crypto/bcrypt"
)

// Interfaz para evitar dependencia circular
type UserRepository interface {
	CreateUser(user User) (int, error)
	GetUserByUsername(username string) (User, error)
}

type User struct {
	Id       uint   `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	// pointer to allow NULL values
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db UserRepository) (string, error) {

	var err error

	user, err := db.GetUserByUsername(username)

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, user.Password)

	if err != nil {
		return "", err
	}

	token, err := token.GenerateToken(user.Id)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) SaveUser(repo UserRepository) (*User, error) {
	err := u.BeforeSave()
	if err != nil {
		return nil, err
	}

	// Pasar el valor, no el puntero
	_, err = repo.CreateUser(*u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}
