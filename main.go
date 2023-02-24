package main

import (
	budget "api-your-accounts/budget/infrastructure/gorm/model"
	project "api-your-accounts/project/infrastructure/gorm/model"
	"api-your-accounts/shared/infrastructure"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.Llongfile + log.LUTC + log.LstdFlags)
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Load .env file")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}

	dsn := "host=192.168.1.14 user=postgres password=test dbname=gorm port=5432 sslmode=disable TimeZone=America/Bogota"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&project.Project{}, &budget.Budget{}, &budget.BudgetAvailableBalance{}, &budget.CategoryBill{}, &budget.BudgetBill{}, &budget.BudgetBillTransaction{}, &budget.BudgetBillShared{})

	infrastructure.NewServer()
}
