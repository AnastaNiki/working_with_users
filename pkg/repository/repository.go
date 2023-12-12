package repository

import (
	"github.com/jmoiron/sqlx"
	"working_with_users"
)

type WorkWithUser interface {
	CreateUser(user working_with_users.User) (int, string, string, string, string, bool, working_with_users.ModifyError)
	UpdateUser(user working_with_users.User, updatePatronymic bool) (int, string, string, string, string, bool, working_with_users.ModifyError)
	ArchiveUser(user working_with_users.User) (int, working_with_users.ModifyError)
	GetUserById(id int) (int, string, string, string, string, bool, working_with_users.ModifyError)
	FindUsers(userPaging working_with_users.UserPaging) ([]working_with_users.User,
		working_with_users.PagingResult, working_with_users.ModifyError)
}

type Repository struct {
	WorkWithUser
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		WorkWithUser: NewCreateUserPostgres(db),
	}
}
