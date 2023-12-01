package models

import (
	"context"
	"errors"
)

type UserCtxKey string

const (
	Userkey UserCtxKey = "user"
)

func SetUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, Userkey, user)
}
func GetUserByContext(ctx context.Context) (*User, error) {
	v := ctx.Value(Userkey)
	u, ok := v.(*User)
	if !ok {
		return nil, errors.New("casting user went wrong")
	}
	return u, nil
}
