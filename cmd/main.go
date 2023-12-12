package main

import (
	"flag"
	"github.com/siruspen/logrus"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"working_with_users"
	"working_with_users/pkg/handler"
	"working_with_users/pkg/repository"
	"working_with_users/pkg/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s, err.Error()")
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initializing db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(working_with_users.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error ocured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	//go run ./cmd/web -name="your_config" -path="folder/your_path"
	path := flag.String("path", "configs", "Путь к файлу конфигурации")
	name := flag.String("name", "config", "Имя файла конфигурации")
	flag.Parse()
	viper.AddConfigPath(*path)
	viper.SetConfigName(*name)
	return viper.ReadInConfig()
}
