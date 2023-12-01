package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrorPasswordEmailInvalid = errors.New("email ou mot de passe invalid")

// type UserDB struct {
// 	id             string
// 	username       string
// 	passwordHashed string
// 	email          string
// 	verified       bool
// 	createdAt      time.Time
// 	updatedAt      sql.NullTime
// 	lastLogin      sql.NullTime
// }
type UserCreate struct {
	Username       string
	Email          string
	Password       string
	RepeatPassword string
}

type User struct {
	ID        string
	Username  string
	Email     string
	Verified  bool
	LastLogin time.Time
}
type UserManager struct {
	DB *sql.DB
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate password :%w", err)
	}
	return string(b), nil
}

func ComparedHashesPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (um *UserManager) Create(ctx context.Context, username, email, password string) (*User, error) {
	h, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("something went wrong on hashing")
	}
	newUser := UserCreate{
		Username: username,
		Email:    email,
		Password: h,
	}
	var user User
	err = um.DB.QueryRowContext(
		ctx, `insert into users (username, email, password_hashed) values ($1, $2, $3) returning id, username, email, verified, last_login;`,
		newUser.Username,
		newUser.Email,
		newUser.Password,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Verified, &user.LastLogin)
	if err != nil {
		return nil, fmt.Errorf("insert user : %w", err)
	}
	return &user, nil
}

// Search the email in table users if found the users already exists return true
func (um *UserManager) SearchUserByEmail(ctx context.Context, email string) bool {
	var exists bool
	err := um.DB.QueryRowContext(ctx, "select exists(select 1 from users where email = $1);", email).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}
	return exists
}
func (um *UserManager) Login(ctx context.Context, email, password string) (string, error) {
	var id, h string
	err := um.DB.QueryRowContext(ctx, `select id,password_hashed from users where email = $1`, email).Scan(&id, &h)
	if err != nil {
		return "", err
	}
	if err = ComparedHashesPassword(password, h); err != nil {
		return "", ErrorPasswordEmailInvalid
	}
	return id, nil
}

func (um *UserManager) FindUserByCookie(ctx context.Context, token string) (*User, error) {
	um.DB.QueryRowContext(ctx, `select id, username,email, verified from users where  `)
	return nil, nil
}

