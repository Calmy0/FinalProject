package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	. "Booking/data"
)

func main() {
	Database.Connect()
	defer Database.Disconnect()

	r := mux.NewRouter()

	r.HandleFunc("/bookings", GetBookings).Methods("GET")
	r.HandleFunc("/newbooking", PostBooking).Methods("POST")

	r.HandleFunc("/newuser", PostUser).Methods("POST")
	http.ListenAndServe(":3000", r)

}

/* **********************************************
 *   HTTP Handlers
 ************************************************ */

func GetBookings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	RoomNum, err := strconv.ParseInt(q["RoomNum"][0], 10, 32)
	if err != nil {
		fmt.Printf("Error parsing room number: %s\n", q["RoomNum"][0])
	}

	DateFrom, err := strconv.ParseInt(q["StartDate"][0], 10, 32)
	if err != nil {
		fmt.Printf("Error parsing room number: %s\n", q["StartDate"][0])
	}

	DateTill, err := strconv.ParseInt(q["EndDate"][0], 10, 32)
	if err != nil {
		fmt.Printf("Error parsing room number: %s\n", q["EndDate"][0])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Database.ShowBookings(int(RoomNum), int(DateFrom), int(DateTill)))
}

func PostBooking(w http.ResponseWriter, r *http.Request) {
	value := Database.NewBooking()

	err := json.NewDecoder(r.Body).Decode(value)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("new bookings entry: %+v", value)

	fmt.Fprint(w, Database.AddBooking(value))
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	value := Database.NewUser()

	err := json.NewDecoder(r.Body).Decode(value)
	if err != nil {
		fmt.Println("error parsing NewUser's body:", err)
		defer r.Body.Close()

		BodyByteSlice, err1 := ioutil.ReadAll(r.Body)
		if err1 != nil {
			fmt.Println("error read r.Body:", err1)
			os.Exit(2)
		}
		fmt.Println(string(BodyByteSlice))
	}

	fmt.Println("new users entry: %+v", value)

	fmt.Fprint(w, Database.AddUser(value))
}

func BookingsPut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PUT Method response\n")
}

func BookingsDel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DELETE Method response\n")
}
