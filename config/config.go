package config

import(
	"log"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/akhilesh-ingle-ge/jwt/models"
)

var DB *gorm.DB
var dsn = "postgres://postgres:12345@localhost:5432/jwt?sslmode=disable"

func SetUpDB(){
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal("Error connecting to the database")
	}

	db.AutoMigrate(&models.User{})
	DB = db
}