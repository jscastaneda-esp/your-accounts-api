package db

import (
	"log"
	"os"

	budget "api-your-accounts/budget/infrastructure/entity"
	project "api-your-accounts/project/infrastructure/entity"
	user "api-your-accounts/user/infrastructure/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Postgres DB
var DB *gorm.DB

func NewDB() {
	if DB == nil {
		log.Println("Init connection to database")

		dsn := os.Getenv("DATABASE_DSN")
		if dsn == "" {
			log.Fatal("Environment variable DATABASE_DSN is mandatory")
		}

		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err != nil {
			log.Fatal(err)
		}

		err = DB.AutoMigrate(
			&user.User{},
			&user.UserToken{},
			&project.Project{},
			&project.ProjectLog{},
			&budget.Budget{},
			&budget.BudgetAvailableBalance{},
			&budget.CategoryBill{},
			&budget.BudgetBill{},
			&budget.BudgetBillTransaction{},
			&budget.BudgetBillShared{},
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
