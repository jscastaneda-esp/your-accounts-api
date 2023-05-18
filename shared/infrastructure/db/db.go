package db

import (
	"log"
	"os"
	"time"

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
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.Llongfile+log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		)
		if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			log.Fatal(err)
		}

		if err = DB.AutoMigrate(
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
		); err != nil {
			log.Fatal(err)
		}
	}
}
