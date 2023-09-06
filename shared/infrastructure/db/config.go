package db

import (
	"log"
	"os"
	"time"

	budgets "your-accounts-api/budgets/infrastructure/db/entity"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/config"
	shared "your-accounts-api/shared/infrastructure/db/entity"
	users "your-accounts-api/users/infrastructure/db/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Postgres DB
var DB *gorm.DB
var Tm persistent.TransactionManager

func NewDB() {
	if DB == nil {
		log.Println("Init connection to database")

		var err error
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.Llongfile+log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		)
		if DB, err = gorm.Open(mysql.Open(config.DATABASE_DSN), &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			log.Fatal(err)
		}

		if err = DB.AutoMigrate(
			new(users.User),
			new(users.UserToken),
			new(budgets.Budget),
			new(budgets.BudgetAvailableBalance),
			new(budgets.BudgetBill),
			new(shared.Log),
		); err != nil {
			log.Fatal(err)
		}

		Tm = NewTransactionManager(DB)
	}
}
