package model

import (
	"fmt"
	"log"
	"streelity/v1/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/event"
)

// type Configuration struct {
// 	Server   string
// 	Database string `json:"dbname"`
// 	Username string
// 	Password string
// }

var Db *gorm.DB

// var Config Configuration
var OnDisconnect *event.Event = event.NewEvent()
var OnConnected *event.Event = event.NewEvent()

// func loadConfig(path string) {
// 	file, fileErr := os.Open(path)
// 	if fileErr != nil {

// 		log.Panic(fileErr)
// 	}

// 	defer file.Close()
// 	decoder := json.NewDecoder(file)
// 	Config = Configuration{}

// 	err := decoder.Decode(&Config)

// 	if err != nil {
// 		log.Panic(err)
// 	}

// }

func connect() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.Config.Username, config.Config.Password, config.Config.Server, config.Config.Database)
	db, err := gorm.Open("mysql", connectionString)
	Db = db

	if err != nil {
		OnDisconnect.Invoke()
		log.Println(err.Error())
	} else {
		OnConnected.Invoke()
		log.Println("[Database] connect success")
	}

}

func reconnect() {
	timer := time.NewTimer(10 * time.Second)
	<-timer.C

	connect()
}

func init() {
	// _, b, _, _ := runtime.Caller(0)
	// basepath := filepath.Dir(b)
	// configPath := filepath.Join(filepath.Dir(basepath), "config", "config.json")

	// loadConfig(configPath)
	OnDisconnect.Subscribe(reconnect)
	// OnConnected.Subscribe(LoadService)
	go connect()

	log.Println("Hi mom init db")
}
