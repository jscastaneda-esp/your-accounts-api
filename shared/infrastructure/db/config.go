package db

import (
	golog "log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"

	budgets "your-accounts-api/budgets/infrastructure/db/entity"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/config"
	shared "your-accounts-api/shared/infrastructure/db/entity"
	users "your-accounts-api/users/infrastructure/db/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Postgres DB
var DB *gorm.DB
var Tm persistent.TransactionManager

func NewDB() {
	if DB == nil {
		var err error
		newLogger := logger.New(
			golog.New(os.Stdout, "\r\n", golog.Llongfile+golog.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		)

		log.Info("Init connection to database")
		if DB, err = gorm.Open(postgres.Open(config.DATABASE_DSN), &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			log.Fatal(err)
		}

		if err = DB.AutoMigrate(
			new(users.User),
			new(users.UserToken),
			new(budgets.Budget),
			new(budgets.BudgetAvailable),
			new(budgets.BudgetBill),
			new(shared.Log),
		); err != nil {
			log.Fatal(err)
		}

		Tm = NewTransactionManager(DB)
	}
}
