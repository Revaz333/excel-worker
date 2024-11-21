package main

import (
	"fmt"
	"worker/db"
	"worker/internal/controllers"
	"worker/internal/helpers"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("config init error:%s", err.Error())
	}

	db, err := db.NewDb(db.Config{
		Host:       viper.GetString("db.host"),
		Port:       viper.GetString("db.port"),
		Dbname:     viper.GetString("db.name"),
		Dbuser:     viper.GetString("db.user"),
		Dbpassword: viper.GetString("db.password"),
	})

	if err != nil {
		logrus.Fatalf("Something go wrong while connect to Db(")
	}

	dbHelpers := helpers.Db{db}
	arrHelpers := helpers.Arrays{}
	xlHelpers := helpers.Exel{}
	controller := controllers.ExportController{dbHelpers, arrHelpers, xlHelpers}

	res, err := controller.CheckQueues()

	if err != nil {
		logrus.Fatalf("Something go wrong")
	}

	if res == false {
		fmt.Println("Queue is empty")
	}

	fmt.Println("export success")
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
