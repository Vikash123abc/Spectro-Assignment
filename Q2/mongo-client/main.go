package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DbURL          string
	PoolSize       uint64
	initMongo      sync.Once
	err            error
	ServerShutdown bool
	MongoConn      *conn
)

type conn struct {
	Client *mongo.Client
	client *mongo.Client
}

// Disconnect mongo instance
func (c *conn) Disconnect() error {
	err := c.client.Disconnect(context.TODO())
	if err != nil {
		return errors.New("Failed to disconnect mongo")
	}
	return nil
}

func Connect() error {

	var client *mongo.Client

	monitor := &event.PoolMonitor{
		Event: handlePoolMonitor,
	}

	initMongo.Do(func() {
		clientOpts := options.Client().ApplyURI(DbURL).SetMaxPoolSize(PoolSize).
			SetPoolMonitor(monitor).SetHeartbeatInterval(time.Duration(10 * time.Second)).
			SetMonitor(apmmongo.CommandMonitor())
		client, err = mongo.Connect(context.TODO(), clientOpts)
	})

	if err != nil {
		return errors.New("Failed to connect to mongo")
	}

	if err = Ping(client); err != nil {
		log.Println("Error while pinging the DB client", err.Error())
		return err
	}
	MongoConn = &conn{client: client, Client: client}

	return nil
}
func handlePoolMonitor(evt *event.PoolEvent) {
	if !ServerShutdown {
		switch evt.Type {
		case event.PoolClosedEvent:
			log.Println("DB Pool closed...Reconnecting...")
			if err = Connect(); err != nil {
				log.Fatalln("Failed to connect to mongo")
				break
			}
		}
	}
}

// Ping :
func Ping(c *mongo.Client) error {
	err = c.Ping(context.TODO(), nil)
	return err
}

func main() {
	fmt.Println("Hello world")
	Connect()
	// Since the Connectin string is not given, it won't connect for now
	fmt.Printf("DB connected succesfully")
}
