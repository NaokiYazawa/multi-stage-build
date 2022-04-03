// Entities
package model

import (
	"errors"
)

// User Struct of an user
// Userを表すドメインモデルであるから、データベース関連のことは記述しない
type User struct {
	ID   int
	Name string
}

// User Struct of users
type Users []User

// NewUser Constructor of an user
// ビジネスのルールが分かる。
func NewUser(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("名前は必須です。")
	}

	user := &User{
		Name: name,
	}

	return user, nil
}

// SetUser Setter of an User
// This function used to update the user.
func (user *User) SetUser(name string) error {
	if name == "" {
		return errors.New("名前は必須です。")
	}

	user.Name = name

	return nil
}
