package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* **********************************************
 *   DataBase
 ************************************************ */
type DB_T struct {
	server      string
	database    string
	colBookings string
	colRooms    string
	colUsers    string
	client      *mongo.Client
}

var Database = DB_T{
	server:      "mongodb://localhost:27017",
	database:    "Hotel",
	colBookings: "Booking",
	colRooms:    "Rooms",
	colUsers:    "Users",
}

func (db *DB_T) Connect() {
	// create client
	var err error
	db.client, err = mongo.NewClient(options.Client().ApplyURI(db.server))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = db.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to MongoDB!")
}

func (db *DB_T) Disconnect() {
	err := db.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection closed")
}
