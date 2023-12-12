#### Проект: "Работа с пользователями".  Реализация серверной части - бэкенд
##### Подробное описание сценариев и методов в ТЗ


Для реализации серверной части проекта выбраны следующие инструменты:
- Go
- Веб-фреймворк Gin
- База данных PostgreSQL
- Библиотеки:
    - "github.com/joho/godotenv" для получение пароля из переменной окружения
    - "github.com/spf13/viper" для считвания файлов кнфигурации
    - "github.com/jmoiron/sqlx" для работы с БД


Путь к файлу конфигурации БД задаётся след. образом:
```
go run ./cmd/web -name="your_config" -path="folder/your_path
```

Если не указать, будет по умолчанию path=configs, name=config
```
go run ./cmd/
```

Cтруктура файла:

```
port: "8000"

db:
    username: "postgres"
    host: "localhost"
    port: "5432"
    dbname: "working_with_users_db"
    sslmode: "disable"
```

Для настройки пароля от БД отдельный файл .env, который считывается библиотекой github.com/joho/godotenv
(Необходимо в корне проекта создать файл .env с переменной DB_PASSWORD=...)

Для создания структуры БД есть файлы миграции, 
они подойдут даже, если БД будет в контейнере docker.
Возможно пригодится установка migrate:
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
Подходит для:
go version go1.17.4 windows/amd64

Пример запуска миграции:
migrate -path ./schema -database 'postgres://postgres:password@localhost:5432/working_with_users_db?sslmode=disable' up

Удалить миграцию:
migrate -path ./schema -database 'postgres://postgres:password@localhost:5432/working_with_users_db?sslmode=disable' down


Серверная часть проекта реализована с использованием веб-фреймворка Gin и схемы MVC. Логически приложение разделено на три уровня: 
1)  Handler. Отвечает за принятие и обработку http запросов, а так же за генерацию ответов и отправку ошибок.
2)  Repository. Отвечает за работу с базой данных, запрос и изменение данных
3)  Service. Отображает структуру бизнес логики, отвтечает за связь компонентов.
handler -> service -> repository(postgres)
Такая схема позволяет модифицировать компоненты независомо. Например, есть возможность поменять базу данных, не затрагивая логику обработки http запросов.

MVC (Model-View-Controller) - схема разделения данных приложения, пользовательского интерфейса и управляющей логики на три отдельных компонента: модель, представление и контроллер – таким образом, что модификация каждого компонента может осуществляться независимо
