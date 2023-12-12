package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/siruspen/logrus"
	"strconv"
	"working_with_users"
)

type CreateUserPostgres struct {
	db *sqlx.DB
}

func NewCreateUserPostgres(db *sqlx.DB) *CreateUserPostgres {
	return &CreateUserPostgres{db: db}
}

func (r *CreateUserPostgres) CreateUser(user working_with_users.User) (int, string, string, string, string, bool,
	working_with_users.ModifyError) {
	var id int
	var firstName, lastName, patronymic, login string
	archive := false

	var check_login_archive []working_with_users.User

	err := r.db.Select(&check_login_archive, fmt.Sprintf("SELECT archive FROM %s WHERE login = $1", usersTable), user.Login)
	if err != nil {
		return 0, "", "", "", "", false,
			working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
	}

	is_same_login := false
	for i := 0; i < len(check_login_archive); i++ {
		if check_login_archive[i].Archive == false {
			is_same_login = true
		}
	}
	if is_same_login == false {
		query := fmt.Sprintf(
			"INSERT INTO %s (firstName, lastName, patronymic, login, archive) values ($1, $2, $3, $4, $5) RETURNING id, firstName, lastName, patronymic, login, archive",
			usersTable)

		row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Patronymic, user.Login, false)
		if err := row.Scan(&id, &firstName, &lastName, &patronymic, &login, &archive); err != nil {
			return 0, firstName, lastName, patronymic, login, archive,
				working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
		}

		return id, firstName, lastName, patronymic, login, archive, working_with_users.ModifyError{}
	} else {
		var errDetail []working_with_users.Detail
		errDetail = append(errDetail, working_with_users.Detail{Key: "login", Value: user.Login})
		return id, firstName, lastName, patronymic, login, archive,
			working_with_users.NewModifyError(
				"userAlreadyExistError", fmt.Sprintf(
					"Пользователь с указанным логином '%s' уже существует", user.Login), errDetail)
	}

}

func (r *CreateUserPostgres) UpdateUser(user working_with_users.User, updatePatronymic bool) (int,
	string, string, string, string, bool,
	working_with_users.ModifyError) {

	var check_user_by_id []working_with_users.User
	err := r.db.Select(&check_user_by_id,
		fmt.Sprintf("SELECT id, firstName, lastName, patronymic, login, archive FROM %s WHERE id = $1",
			usersTable), user.Id)
	if err != nil {
		return 0, "", "", "", "", false,
			working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
	}

	if len(check_user_by_id) == 0 { //не найдено ничего с таким id
		var errDetail []working_with_users.Detail
		errDetail = append(errDetail, working_with_users.Detail{Key: "Id", Value: strconv.Itoa(user.Id)})
		return 0, "", "", "", "", false,
			working_with_users.NewModifyError("userNotFoundError",
				fmt.Sprintf("Пользователь с указанным идентификатором '%d' не найден", user.Id), errDetail)
	}

	if check_user_by_id[0].Archive == true {
		var errDetail []working_with_users.Detail
		errDetail = append(errDetail, working_with_users.Detail{Key: "Id", Value: strconv.Itoa(user.Id)})
		return 0, "", "", "", "", false, working_with_users.NewModifyError("userIsArchiveError",
			fmt.Sprintf("Пользователь с указанным идентификатором '%d' является архивным", user.Id), errDetail)
	}

	var check_user_login []working_with_users.User
	err = r.db.Select(&check_user_login,
		fmt.Sprintf(
			"SELECT id, firstName, lastName, patronymic, login, archive FROM %s WHERE id != $1 AND login = $2 AND archive = $3",
			usersTable), user.Id, user.Login, false)
	if err != nil {
		return 0, "", "", "", "", false, working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
	}

	if len(check_user_login) != 0 {
		//err := errors.New(fmt.Sprintf("Пользователь с указанным логином '%s' уже существует", user.Login))
		//return 0, "", "", "", "", false, err
		var errDetail []working_with_users.Detail
		errDetail = append(errDetail, working_with_users.Detail{Key: "login", Value: strconv.Itoa(check_user_login[0].Id)})
		return 0, "", "", "", "", false, working_with_users.NewModifyError(
			"userAlreadyExistError", fmt.Sprintf(
				"Пользователь с указанным логином '%s' уже существует", user.Login), errDetail)
	}

	setQuery := ""
	if check_user_by_id[0].FirstName != user.FirstName {
		setQuery += "firstName =  " + "'" + user.FirstName + "'" + ", "
	}
	if check_user_by_id[0].LastName != user.LastName {
		setQuery += "lastName =  " + "'" + user.LastName + "'" + ", "
	}
	if check_user_by_id[0].Patronymic != user.Patronymic && updatePatronymic == true {
		setQuery += "patronymic =  " + "'" + user.Patronymic + "'" + ", "
	}
	if check_user_by_id[0].Login != user.Login {
		setQuery += "login =  " + "'" + user.Login + "'" + ", "
	}

	if setQuery == "" {
		return check_user_by_id[0].Id, check_user_by_id[0].FirstName, check_user_by_id[0].LastName,
			check_user_by_id[0].Patronymic, check_user_by_id[0].Login, check_user_by_id[0].Archive, working_with_users.ModifyError{}
	} else {
		setQuery = setQuery[:len(setQuery)-2]
		query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d`, usersTable, setQuery, user.Id)
		_, err := r.db.Exec(query)

		if err != nil {
			logrus.Fatalf(setQuery, err.Error())
			return 0, "", "", "", "", false,
				working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
		}

		var check_update []working_with_users.User
		err = r.db.Select(&check_update, fmt.Sprintf(
			"SELECT id, firstName, lastName, patronymic, login, archive FROM %s WHERE id = $1",
			usersTable), user.Id)
		if err != nil {
			return 0, "", "", "", "", false,
				working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
		}
		return user.Id, check_update[0].FirstName, check_update[0].LastName,
			check_update[0].Patronymic, check_update[0].Login, check_update[0].Archive, working_with_users.ModifyError{}
	}
}

func (r *CreateUserPostgres) ArchiveUser(user working_with_users.User) (int, working_with_users.ModifyError) {

	var find_user []working_with_users.User
	err := r.db.Select(&find_user, fmt.Sprintf("SELECT archive FROM %s WHERE id = $1", usersTable), user.Id)
	if err != nil {
		return 0, working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
	}

	if len(find_user) == 0 {
		var errDetail []working_with_users.Detail
		errDetail = append(errDetail, working_with_users.Detail{Key: "Id", Value: strconv.Itoa(user.Id)})
		return 0,
			working_with_users.NewModifyError("userNotFoundError",
				fmt.Sprintf("Пользователь с указанным идентификатором '%d' не найден", user.Id), errDetail)
	}

	if find_user[0].Archive == true {
		return user.Id, working_with_users.ModifyError{}
	}

	query := fmt.Sprintf(`UPDATE %s SET archive = 'true' WHERE id = %d`, usersTable, user.Id)
	_, err = r.db.Exec(query)

	if err != nil {
		return 0, working_with_users.NewModifyError("databaseError", err.Error(), []working_with_users.Detail{})
	}

	return user.Id, working_with_users.ModifyError{}
}

func (r *CreateUserPostgres) GetUserById(id int) (int, string, string, string, string, bool, working_with_users.ModifyError) {

	var find_user []working_with_users.User
	err := r.db.Select(&find_user, fmt.Sprintf(
		"SELECT  id, firstName, lastName, patronymic, login, archive FROM %s WHERE id = $1",
		usersTable), id)
	if err != nil {
		return 0, "", "", "", "", false, working_with_users.NewModifyError("databaseError", err.Error(),
			[]working_with_users.Detail{})
	}
	if len(find_user) == 0 {
		var errDetail []working_with_users.Detail
		errDetail = append(errDetail, working_with_users.Detail{Key: "Id", Value: strconv.Itoa(id)})
		return 0, "", "", "", "", false,
			working_with_users.NewModifyError("userNotFoundError",
				fmt.Sprintf("Пользователь с указанным идентификатором '%d' не найден", id), errDetail)
	}

	return find_user[0].Id, find_user[0].FirstName, find_user[0].LastName,
		find_user[0].Patronymic, find_user[0].Login, find_user[0].Archive, working_with_users.ModifyError{}
}

func (r *CreateUserPostgres) FindUsers(userFindParam working_with_users.UserPaging) ([]working_with_users.User,
	working_with_users.PagingResult, working_with_users.ModifyError) {

	var emptyUsers []working_with_users.User
	var emptyPagingResult working_with_users.PagingResult
	var emptyErr working_with_users.ModifyError

	var find_users []working_with_users.User
	findParam := "WHERE "

	if userFindParam.FirstName != "" {
		findParam += "upper(firstName) LIKE upper('%" + userFindParam.FirstName + "%')"
	}

	if userFindParam.LastName != "" {
		if findParam != "WHERE " {
			findParam += " AND "
		}
		findParam += "upper(lastName) LIKE upper('%" + userFindParam.LastName + "%')"
	}

	if userFindParam.Patronymic != "" {
		if findParam != "WHERE " && findParam[len(findParam)-4:] != "AND " {
			findParam += " AND "
		}
		findParam += "upper(patronymic) LIKE upper('%" + userFindParam.Patronymic + "%')"
	}

	if userFindParam.Login != "" {
		if findParam != "WHERE " && findParam[len(findParam)-4:] != "AND " {
			findParam += " AND "
		}
		findParam += "upper(login) LIKE upper('%" + userFindParam.Login + "%')"
	}

	//if userFindParam.Archive != "" { невозможно узнать добавлялась ли архивация по умолчанию будет false
	if findParam != "WHERE " && findParam[len(findParam)-4:] != "AND " {
		findParam += " AND "
	}
	if userFindParam.Archive == true {
		findParam += "archive =" + "true"
	} else {
		findParam += "archive =" + "false"
	}

	sortOrder := "DESC"
	if userFindParam.Paging.Descending == false {
		sortOrder = "ASC"
	}
	offset := (userFindParam.Paging.PageNumber) * userFindParam.Paging.PageSize

	err := r.db.Select(&find_users, fmt.Sprintf(
		"SELECT  id, firstName, lastName, patronymic, login, archive FROM %s %s ORDER BY %s %s LIMIT $1 OFFSET $2",
		usersTable, findParam, userFindParam.Paging.SortField, sortOrder), userFindParam.Paging.PageSize, offset)
	if err != nil {
		return emptyUsers, emptyPagingResult, working_with_users.NewModifyError("databaseError", err.Error(),
			[]working_with_users.Detail{})
	}

	/*type count struct {
		count int64 `json:"count"`
	}
	var countUsers []count*/
	var countUsers []int64
	err = r.db.Select(&countUsers, fmt.Sprintf(
		"SELECT count(*) FROM %s %s",
		usersTable, findParam))
	if err != nil {
		return emptyUsers, emptyPagingResult, working_with_users.NewModifyError("databaseError", err.Error(),
			[]working_with_users.Detail{})
	}

	/*
			1.	Система осуществляет поиск записей в таблице user по следующему правилу:
			upper(user.firstname) like ‘%upper(firstName)%’и
			upper(user.lastname) like ‘%upper(lastName)%’и
			upper(user.patronymic) like ‘%upper(patronymic)%’ и
			upper(user.login) like ‘%upper(login)%’ и
			user.archive = archive
		-strings.ToUpper(str)
		- ORDER BY Manufacturer DESC; по убыванию; ASC, по возрастанию

			Дополнительно Система накладывает критерий сортировки выборки
			по переданным значениям pagingOptions.sortField и
			pagingOptions.descending и параметры пагинации
			LIMIT pagingOptions.pageSize
			OFFSET (pagingOptions.pageNumber) * pagingOptions.pageSize

			Система запоминает найденные записи.
			Примечание: в запросе участвуют только не пустые и переданные поля фильтрации.
			Если никакой фильтр не заполнен, ограничений на выборку не накладывается.

			2.   Система запрашивает общее количество записей в
			таблице user по критерию с шага 2 и запоминает количество записей.
			--SELECT count(*) FROM

			3.    Система заполняет ответ по следующему правилу:
			users – коллекция записей о пользователях с шага 2
			pаgingResult.pageNumber = pagingOptions.pageNumber
			pаgingResult.pageSize – количество элементов, полученных на шаге 2
			pаgingResult.pageTotal - количество элементов, полученных на шаге 3, деленное на pagingOptions.pageSize
			pаgingResult.itemsCount – количество элементов, полученных на шаге 3
			pаgingResult.hasNextPage – результат проверки условия pаgingResult.pageTotal – 1 > pagingOptions.pageNumber


	*/
	pageTotal := int64(0)
	if len(find_users) != 0 {
		pageTotal = countUsers[0] / int64(userFindParam.Paging.PageSize)
		if countUsers[0]%int64(userFindParam.Paging.PageSize) != 0 {
			pageTotal += 1
		}
	}
	nextPage := false
	if len(find_users) != 0 && pageTotal-1 > int64(userFindParam.Paging.PageNumber) {
		nextPage = true
	}

	return find_users,
		working_with_users.PagingResult{userFindParam.Paging.PageNumber, len(find_users), pageTotal,
			countUsers[0], nextPage},
		emptyErr
}
