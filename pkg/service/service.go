package service

import (
	"working_with_users"
	"working_with_users/pkg/repository"
)

type WorkWithUser interface {
	CreateUser(user working_with_users.User) (int, string, string, string, string, bool, working_with_users.ModifyError)
	UpdateUser(user working_with_users.User, updatePatronymic bool) (int, string, string,
		string, string, bool, working_with_users.ModifyError)
	ArchiveUser(user working_with_users.User) (int, working_with_users.ModifyError)
	GetUserById(id int) (int, string, string, string, string, bool, working_with_users.ModifyError)
	FindUsers(userPaging working_with_users.UserPaging) ([]working_with_users.User,
		working_with_users.PagingResult, working_with_users.ModifyError)
}

type Service struct {
	WorkWithUser
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		WorkWithUser: NewCreateService(repos.WorkWithUser),
	}
}
