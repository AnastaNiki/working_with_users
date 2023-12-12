package service

import (
	"working_with_users"
	"working_with_users/pkg/repository"
)

type CreateService struct {
	repo repository.WorkWithUser
}

func NewCreateService(repo repository.WorkWithUser) *CreateService {
	return &CreateService{repo: repo}
}

//имплимитируем логику создания пользователя
func (s *CreateService) CreateUser(user working_with_users.User) (int, string, string, string,
	string, bool, working_with_users.ModifyError) {
	return s.repo.CreateUser(user)
}

func (s *CreateService) UpdateUser(user working_with_users.User, updatePatronymic bool) (int, string, string,
	string, string, bool, working_with_users.ModifyError) {
	return s.repo.UpdateUser(user, updatePatronymic)
}

func (s *CreateService) ArchiveUser(user working_with_users.User) (int, working_with_users.ModifyError) {
	return s.repo.ArchiveUser(user)
}

func (s *CreateService) FindUsers(userPaging working_with_users.UserPaging) ([]working_with_users.User,
	working_with_users.PagingResult, working_with_users.ModifyError) {
	return s.repo.FindUsers(userPaging)
}

func (s *CreateService) GetUserById(id int) (int, string, string, string, string,
	bool, working_with_users.ModifyError) {
	return s.repo.GetUserById(id)
}
