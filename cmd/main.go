package main

import (
	"fmt"
	"worker/db"
	"worker/internal/controllers"

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

	controller := controllers.QueueController{db}

	res, err := controller.CheckQueue()

	if err != nil {
		logrus.Fatalf("Something go wrong")
	}

	if res == false {
		fmt.Println("Queue is empty")
	}

}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
