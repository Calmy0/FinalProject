package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* **********************************************
 *   DataBase Bookings
 ************************************************ */
type BookingEntry struct {
	RoomNum    int    `bson:"RoomNum"`
	StartDate  int    `bson:"StartDate"`
	EndDate    int    `bson:"EndDate"`
	ClientName string `bson:"ClientName"`
	Author     string `bson:"Author"`
}

func (db *DB_T) NewBooking() *BookingEntry {
	return new(BookingEntry)
}

// select all rooms, with bookings in the interval [20170605..20170610]
// - StartDate is less then the end of the interval
// AND
// - EndDate is more then the start of the interval
//{StartDate: {$lte:20170610}, EndDate: {$gte:20170605} }
//{ RoomNum:102, StartDate: {$gte:20170605}, EndDate: {$gte:20170605}  }

func (db *DB_T) ShowBookings(RoomNum, DateFrom, DateTill int) (response []BookingEntry) {
	collection := db.client.Database(db.database).Collection(db.colBookings)

	RoomQuery := bson.E{}
	if RoomNum != 0 {
		RoomQuery = bson.E{"RoomNum", RoomNum}
	}

	DateFromQuery := bson.E{}
	if DateTill != 0 {
		DateFromQuery = bson.E{"StartDate", bson.D{{"$lt", DateTill}}}
	}

	DateTillQuery := bson.E{}
	if DateFrom != 0 {
		DateTillQuery = bson.E{"EndDate", bson.D{{"$gt", DateFrom}}}
	}

	options := options.Find().SetSort(bson.D{{"RoomNum", 1}, {"StartDate", 1}})
	filter := bson.D{
		RoomQuery,
		DateFromQuery,
		DateTillQuery,
	}

	ctx := context.TODO()
	cur, err := collection.Find(ctx, filter, options)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem BookingEntry
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		response = append(response, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(ctx)

	fmt.Printf("result:\n%+v\n", response)
	return response
}

func (db *DB_T) AddBooking(Entry *BookingEntry) (response string) {
	// check for RoomNum
	room := struct {
		Num    int `bson:"Num"`
		Active int `bson:"Active"`
	}{}

	err := db.client.Database(db.database).Collection(db.colRooms).FindOne(
		context.Background(),
		bson.D{
			{"Num", Entry.RoomNum},
		}).Decode(&room)

	if err != nil {
		return fmt.Sprintf("Error. Room %d is not found.", Entry.RoomNum)
	}

	if room.Active != 1 {
		return fmt.Sprintf("Error. Room %d is not active now", Entry.RoomNum)
	}

	// checks for StartDay
	if Entry.StartDate >= Entry.EndDate {
		return "Error. Start day must be greater then end day."
	}

	if Entry.StartDate < (time.Now().Year()*10000 + int(time.Now().Month())*100 + time.Now().Day()) {
		return "Error. Start day must be not less then today."
	}

	// check room availability for specified dates
	check := db.ShowBookings(Entry.RoomNum, Entry.StartDate, Entry.EndDate)
	if len(check) > 0 {
		return fmt.Sprintf("Error. there are bookings for this room:\n %+v", check)
	}

	// id = bson.NewObjectId()
	res, err := db.client.Database(db.database).Collection(db.colBookings).InsertOne(
		context.Background(),
		bson.D{
			{"RoomNum", Entry.RoomNum},
			{"StartDate", Entry.StartDate},
			{"EndDate", Entry.EndDate},
			{"ClientName", Entry.ClientName},
			{"Author", Entry.Author},
		})
	if err != nil {
		return fmt.Sprintf("Database insert error: %s", err.Error)
	}
	return fmt.Sprintf("Entry is added. ID: %s", res.InsertedID)
}
