package mgo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	user   string
	pswd   string
	host   string
	port   int
	dbname string

	database *mongo.Database
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

func NewDb(user, password, host string, port int, dbname string) *MongoDB {
	return &MongoDB{
		user:   user,
		pswd:   password,
		host:   host,
		port:   port,
		dbname: dbname,
	}
}

// https://www.digitalocean.com/community/tutorials/how-to-use-go-with-mongodb-using-the-mongodb-go-driver
// https://www.mongodb.com/golang
func (m *MongoDB) Connect() error {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d", m.user, m.pswd, m.host, m.port)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	m.database = client.Database(m.dbname)
	return nil
}
