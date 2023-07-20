package db

import (
	"log"
	"os"
	"time"

	budget "your-accounts-api/budget/infrastructure/entity"
	project "your-accounts-api/project/infrastructure/entity"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/config"
	infra "your-accounts-api/shared/infrastructure/db/persistent"
	user "your-accounts-api/user/infrastructure/entity"

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
			&user.User{},
			&user.UserToken{},
			&project.Project{},
			&project.ProjectLog{},
			&budget.Budget{},
			&budget.BudgetAvailableBalance{},
			&budget.BudgetBill{},
			&budget.BudgetBillTransaction{},
			&budget.BudgetBillShared{},
		); err != nil {
			log.Fatal(err)
		}

		Tm = infra.NewTransactionManager(DB)
	}
}
