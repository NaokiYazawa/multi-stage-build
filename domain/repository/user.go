package repository

import (
	"multi-stage-build/domain/model"
)

// メソッド名(とコメント)で機能を「定義」だけする。
// データにアクセスするためのインターフェース
// Data Access Interface
// UserRepository Interface of an user repository
// UserRepository is interface for infrastructure
// Clean Architecture の依存性の順番を守るために interface（抽象）のみを定義し、実際の実装は infra 層で行う。
type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	ReadByID(id int) (*model.User, error)
	ReadAll() (*model.Users, error)
	Update(user *model.User) (*model.User, error)
	Delete(user *model.User) error
}
