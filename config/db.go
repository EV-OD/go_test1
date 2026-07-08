package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	var dsn string
	dsn = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", cfg.Database.HOST, cfg.Database.USER, cfg.Database.PASSWORD, cfg.Database.DB, cfg.Database.PORT)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

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
