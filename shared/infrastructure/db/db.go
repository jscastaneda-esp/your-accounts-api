// TODO: Pendientes tests

package db

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	budget "api-your-accounts/budget/infrastructure/entity"
	project "api-your-accounts/project/infrastructure/entity"
	selfMongo "api-your-accounts/shared/infrastructure/db/mongo"
	user "api-your-accounts/user/infrastructure/entity"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	}
}

// Mongo DB
const (
	defaultTimeout = 30 * time.Second
)

var MongoClient *selfMongo.Client

func NewMongoClient() {
	if MongoClient == nil {
		log.Println("Init client to mongo database")

		uri := os.Getenv("MONGO_URL")
		if uri == "" {
			log.Fatal("Environment variable MONGODB_URI is mandatory")
		}

		dbName := os.Getenv("MONGODB")
		if dbName == "" {
			log.Fatal("Environment variable MONGODB is mandatory")
		}

		timeout := defaultTimeout
		if timeoutEnv := os.Getenv("MONGOTIMEOUT"); timeoutEnv != "" {
			val, err := strconv.Atoi(timeoutEnv)
			if err != nil {
				log.Fatal("Environment variable MONGOTIMEOUT invalid value")
			}

			timeout = time.Duration(val) * time.Second
		}

		serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
		client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions))
		if err != nil {
			log.Fatal("Create the Pool failed:", err)
		}

		MongoClient = selfMongo.New(dbName, timeout, client)
		if err := MongoClient.Connect(context.TODO()); err != nil {
			log.Fatal("Error connecting in database:", err)
		}
		if err := MongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal("Error disconnect in database:", err)
		}
	}
}
