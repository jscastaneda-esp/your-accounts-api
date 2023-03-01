// TODO: Pendientes tests

package infrastructure

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	budget "api-your-accounts/budget/infrastructure"
	project "api-your-accounts/project/infrastructure"
	"api-your-accounts/shared/infrastructure/mongodb"
	user "api-your-accounts/user/infrastructure"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Postgres DB
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

// Mongo DB
const (
	defaultTimeout = 30 * time.Second
)

func NewMongoClient() *mongodb.Client {
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

	mongoClient := mongodb.New(dbName, timeout, client)
	if err := mongoClient.Connect(context.TODO()); err != nil {
		log.Fatal("Error connecting in database:", err)
	}
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatal("Error disconnect in database:", err)
	}

	return mongoClient
}
