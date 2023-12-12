package handler

import (
	//"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"working_with_users"
)

func (h *Handler) createUser(c *gin.Context) {
	var input working_with_users.User

	//для использования с ошибками типа ErrorResponse
	/*if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}*/

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) //ошибка возвращается в ответе
		return
	}

	var errDetail []working_with_users.Detail
	if input.Login == "" {
		errDetail = append(errDetail, working_with_users.Detail{Key: "login", Value: "null"})
	}
	if input.FirstName == "" {
		errDetail = append(errDetail, working_with_users.Detail{Key: "firstName", Value: "null"})
	}
	if input.LastName == "" {
		errDetail = append(errDetail, working_with_users.Detail{Key: "lastName", Value: "null"})
	}

	if len(errDetail) != 0 {
		message := "В запросе отсутствуют обязательные параметры"
		working_with_users.RunModifyError(c, http.StatusBadRequest, working_with_users.NewModifyError(
			"validationError", message, errDetail))
		return
	}

	id, lastName, firstName, patronymic, login, archive, modErr := h.services.WorkWithUser.CreateUser(input)
	if modErr.Code != "" {
		//newErrorResponse(c, http.StatusInternalServerError, err.Error())
		working_with_users.RunModifyError(c, http.StatusInternalServerError, modErr)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":         id,
		"lastName":   lastName,
		"firstName":  firstName,
		"patronymic": patronymic,
		"login":      login,
		"archive":    archive,
	})
}

func (h *Handler) updateUser(c *gin.Context) {
	var input working_with_users.User
	updatePatronymic := true
	input.Id = -14687984213598                                               //для проверки введен ли id в запросе
	input.Patronymic = "non-existent patronymic for checking an empty value" //для проверки введен ли patronymic в запросе

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) //ошибка возвращается в ответе
		return
	}

	var errDetail []working_with_users.Detail
	if input.Id == -14687984213598 {
		errDetail = append(errDetail, working_with_users.Detail{Key: "id", Value: "null"})
	}

	if len(errDetail) != 0 {
		message := "В запросе отсутствуют обязательные параметры"
		working_with_users.RunModifyError(c, http.StatusBadRequest,
			working_with_users.NewModifyError("validationError", message, errDetail))
		return
	}

	if input.Patronymic == "non-existent patronymic for checking an empty value" {
		updatePatronymic = false
		input.Patronymic = ""
	}

	id, lastName, firstName, patronymic, login, archive, modErr := h.services.WorkWithUser.UpdateUser(input, updatePatronymic)
	if modErr.Code != "" {
		if modErr.Code == "userNotFoundError" { //404
			working_with_users.RunModifyError(c, http.StatusNotFound, modErr)
			return
		}
		if modErr.Code == "userIsArchiveError" || modErr.Code == "userAlreadyExistError" { //400
			working_with_users.RunModifyError(c, http.StatusBadRequest, modErr)
			return
		}

		working_with_users.RunModifyError(c, http.StatusInternalServerError, modErr) //500
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":         id,
		"lastName":   lastName,
		"firstName":  firstName,
		"patronymic": patronymic,
		"login":      login,
		"archive":    archive,
	})

}

func (h *Handler) archiveUser(c *gin.Context) {
	var input working_with_users.User
	input.Id = -14687984213598

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) //ошибка возвращается в ответе
		return
	}

	var errDetail []working_with_users.Detail
	if input.Id == -14687984213598 {
		errDetail = append(errDetail, working_with_users.Detail{Key: "id", Value: "null"})
	}

	if len(errDetail) != 0 {
		message := "В запросе отсутствуют обязательные параметры"
		working_with_users.RunModifyError(c, http.StatusBadRequest,
			working_with_users.NewModifyError("validationError", message, errDetail))
		return
	}

	id, modErr := h.services.WorkWithUser.ArchiveUser(input)

	if modErr.Code != "" {
		if modErr.Code == "userNotFoundError" { //404
			working_with_users.RunModifyError(c, http.StatusNotFound, modErr)
			return
		}

		working_with_users.RunModifyError(c, http.StatusInternalServerError, modErr) //500
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) findUsers(c *gin.Context) {

	emptyPaging := working_with_users.PagingOptions{-123213343224, -123213343224, "null", false}
	var input working_with_users.UserPaging
	input.Paging = emptyPaging

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) //ошибка возвращается в ответе
		return
	}

	if input.Paging != emptyPaging {
		var errDetail []working_with_users.Detail
		if input.Paging.PageNumber == -123213343224 {
			errDetail = append(errDetail, working_with_users.Detail{Key: "pageNumber", Value: "null"})
		}
		if input.Paging.PageSize == -123213343224 {
			errDetail = append(errDetail, working_with_users.Detail{Key: "pageSize", Value: "null"})
		}
		if input.Paging.SortField == "null" {
			errDetail = append(errDetail, working_with_users.Detail{Key: "sortField", Value: "null"})
		}
		//проверку поля input.Paging.Descending не установить

		if len(errDetail) != 0 {
			message := "В запросе отсутствуют обязательные параметры"
			working_with_users.RunModifyError(c, http.StatusBadRequest, working_with_users.NewModifyError(
				"validationError", message, errDetail))
			return
		}

		if input.Paging.PageNumber < 0 {
			errDetail = append(errDetail, working_with_users.Detail{Key: "pageNumber",
				Value: strconv.Itoa(input.Paging.PageNumber)})
		}

		if input.Paging.PageSize < 1 {
			errDetail = append(errDetail, working_with_users.Detail{Key: "pageSize",
				Value: strconv.Itoa(input.Paging.PageSize)})
		}

		if input.Paging.SortField != "lastName" && input.Paging.SortField != "firstName" &&
			input.Paging.SortField != "patronymic" && input.Paging.SortField != "login" &&
			input.Paging.SortField != "archive" && input.Paging.SortField != "id" {
			errDetail = append(errDetail, working_with_users.Detail{Key: "sortField", Value: input.Paging.SortField})
		}

		if len(errDetail) != 0 {
			message := "Переданы некорректные параметры запроса"
			working_with_users.RunModifyError(c, http.StatusBadRequest, working_with_users.NewModifyError(
				"validationError", message, errDetail))
			return
		}

	} else {
		input.Paging.PageNumber = 0
		input.Paging.PageSize = 20
		input.Paging.SortField = "id"
		input.Paging.Descending = false
	}

	users, pegingResult, modErr := h.services.WorkWithUser.FindUsers(input)

	if modErr.Code != "" {
		working_with_users.RunModifyError(c, http.StatusInternalServerError, modErr) //500
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"users":        users,
		"pagingResult": pegingResult,
	})
	return
}

func (h *Handler) getUserById(c *gin.Context) {

	input_id, err := strconv.Atoi(c.Param("id"))
	if err != nil || input_id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id param"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) //ошибка возвращается в ответе
		return
	}

	id, lastName, firstName, patronymic, login, archive, modErr := h.services.WorkWithUser.GetUserById(input_id)
	if modErr.Code != "" {
		if modErr.Code == "userNotFoundError" { //404
			working_with_users.RunModifyError(c, http.StatusNotFound, modErr)
			return
		}

		working_with_users.RunModifyError(c, http.StatusInternalServerError, modErr) //500
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":         id,
		"lastName":   lastName,
		"firstName":  firstName,
		"patronymic": patronymic,
		"login":      login,
		"archive":    archive,
	})

}
