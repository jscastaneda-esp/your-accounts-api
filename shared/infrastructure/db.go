// TODO: Pendientes tests

package infrastructure

import (
	"log"
	"os"

	budget "api-your-accounts/budget/infrastructure"
	project "api-your-accounts/project/infrastructure"
	user "api-your-accounts/user/infrastructure"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	log.Println("Init connection to database")

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("Environment variable DATABASE_DSN is mandatory")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&user.User{},
		&project.Project{},
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

	return db
}
