package db

import (
	"log"
	"os"
	"time"

	budget "your-accounts-api/budget/infrastructure/entity"
	project "your-accounts-api/project/infrastructure/entity"
	user "your-accounts-api/user/infrastructure/entity"

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
			log.New(os.Stdout, "\r\n", log.Llongfile+log.LstdFlags),
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

		if err = declareEnums(DB); err != nil {
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
	}
}

func declareEnums(db *gorm.DB) error {
	if err := db.Exec(`
	DO $$ BEGIN
		CREATE TYPE project_type AS ENUM ('budget');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`	
	DO $$ BEGIN
		CREATE TYPE budget_bill_category AS ENUM (
			'house',
			'entertainment',
			'personal',
			'vehicle_transportation',
			'education',
			'services',
			'saving',
			'others'
		);
	EXCEPTION
		WHEN duplicate_object THEN null;

		ALTER TYPE budget_bill_category ADD VALUE IF NOT EXISTS 'saving' BEFORE 'others';
	END $$;
	`).Error; err != nil {
		return err
	}

	return nil
}
