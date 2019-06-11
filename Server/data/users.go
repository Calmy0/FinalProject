package data

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

/* **********************************************
 *   DataBase Users
 ************************************************ */
type UsersEntry struct {
	Login    string `bson:"Login"`
	Password string `bson:"Password"`
	Name     string `bson:"Name"`
	Type     string `bson:"Type"`
}

func (db *DB_T) NewUser() *UsersEntry {
	return new(UsersEntry)
}

func (db *DB_T) AddUser(Entry *UsersEntry) (response string) {
	value := db.NewUser()
	// check whether login is busy
	err := db.client.Database(db.database).Collection(db.colUsers).FindOne(
		context.Background(),
		bson.D{
			{"Login", Entry.Login},
		}).Decode(&value)

	if err == nil {
		return fmt.Sprintf("Error. Login %s is already exist.", Entry.Login)
	}

	if (strings.Compare(Entry.Type, "admin") != 0) &&
		(strings.Compare(Entry.Type, "agent") != 0) {
		return fmt.Sprintf("Error. Type '%s' is unknown.", Entry.Type)
	}

	res, err2 := db.client.Database(db.database).Collection(db.colUsers).InsertOne(
		context.Background(),
		bson.D{
			{"Login", Entry.Login},
			{"Password", Entry.Password},
			{"Name", Entry.Name},
			{"Type", Entry.Type},
		})
	if err2 != nil {
		return fmt.Sprintf("Database insert error: %s", err2.Error)
	}
	return fmt.Sprintf("Entry is added. ID: %s", res.InsertedID)
}
