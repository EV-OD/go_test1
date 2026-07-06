package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(URL string) (*gorm.DB, error) {

	DB, err := gorm.Open(postgres.Open(URL), &gorm.Config{})

	// test db
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	} else {
		log.Default().Println("Database connected successfully")
	}

	// err = DB.AutoMigrate(&models.EUser{})
	// if err != nil {
	// 	log.Fatalf("Failed to run schema migrations: %v", err)
	// }
	// log.Println("Database schemas migrated successfully!")

	return DB, nil
}
