package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	dbName   string
	timeout  time.Duration
	client   *mongo.Client
	database *mongo.Database
}

func (c *Client) getCxt(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithTimeout(ctx, c.timeout)
}

func (c *Client) Connect(ctx context.Context) error {
	log.Println("Init connection to mongo database")

	ctx, cancel := c.getCxt(ctx)
	defer cancel()
	if err := c.client.Connect(ctx); err != nil {
		return err
	}

	if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	c.database = c.client.Database(c.dbName)
	return nil
}

func (c *Client) GetDatabase() *mongo.Database {
	return c.database
}

func (c *Client) GetCollection(collection string) *mongo.Collection {
	return c.database.Collection(collection)
}

func (c *Client) Disconnect(ctx context.Context) error {
	log.Println("Close connection to mongo database")

	if err := c.client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func New(dbName string, timeout time.Duration, client *mongo.Client) *Client {
	return &Client{
		dbName:  dbName,
		timeout: timeout,
		client:  client,
	}
}
