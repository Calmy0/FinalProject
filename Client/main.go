package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type cmdT struct {
	name string
	desr string
	exec func() (Response string, ExitFlag int)
}

var cmds [5](*cmdT)

func main() {
	cmds[0] = &cmdHelp
	cmds[1] = &cmdExit
	cmds[2] = &cmdShow
	cmds[3] = &cmdNew
	cmds[4] = &cmdNewUser

	ExitFlag := 0
	response := ""
	input := ""

	fmt.Print("Hotel bookings.\nType 'help' to see available commands.")

ReadCommand:
	for ExitFlag == 0 {
		print("\n>")
		_, err := fmt.Scanln(&input)
		if err != nil {
			print("Error while reading.\n")
			continue
		}

		for _, cmd := range cmds {
			if input == cmd.name {
				response, ExitFlag = cmd.exec()
				fmt.Println(response)
				continue ReadCommand
			}
		}
		// if come here, then there weren't match with any of registered commands

		fmt.Print("Unknown command. Type 'help' to see available commands.")
	}
}

var cmdExit = cmdT{
	name: "exit",
	desr: "performs exit from this program",
	exec: func() (Response string, ExitFlag int) {
		return "Exit!", 1
	},
}

var cmdHelp = cmdT{
	name: "help",
	desr: "shows short description and usage of each command",
	exec: func() (Response string, ExitFlag int) {
		for i := 0; i < len(cmds); i++ {
			println("cmd '", cmds[i].name, " ' -", cmds[i].desr)
		}
		return "", 0
	},
}

type bookingClient struct{}

type bookingEntry struct {
	RoomNum    int    `bson:"RoomNum"`
	StartDate  int    `bson:"StartDate"`
	EndDate    int    `bson:"EndDate"`
	ClientName string `bson:"ClientName"`
	Author     string `bson:"Author"`
}

var cmdShow = cmdT{
	name: "show",
	desr: "shows bookings for specified room number and date interval (start date .. end date)",
	exec: func() (Response string, ExitFlag int) {
		type QueryParamsT struct {
			name  string
			input string
		}

		var QueryParams = [...]QueryParamsT{
			{"Room number (0 - if any):", ""},
			{"Start date YYYYMMDD (0 - if any):", ""},
			{"End date YYYYMMDD(0 - if any):", ""},
		}
		println("Enter search parameters.")
		for i := range QueryParams {
			print(QueryParams[i].name)
			_, err := fmt.Scanln(&QueryParams[i].input)
			if err != nil {
				print("Error while reading.\n")
				continue
			}
		}

		// assert params

		resp, err := http.Get(fmt.Sprintf("http://localhost:3000/bookings?RoomNum=%s&StartDate=%s&EndDate=%s",
			QueryParams[0].input,
			QueryParams[1].input,
			QueryParams[2].input))
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		BodyByteSlice, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(2)
		}

		fmt.Println(string(BodyByteSlice))

		var bookings []bookingEntry
		err = json.Unmarshal(BodyByteSlice, &bookings)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("Found bookings:")
		fmt.Println("\tRoom \tStart Date \tEnd Date \tClient \tAuthor")
		for _, entry := range bookings {
			fmt.Println(fmt.Sprintf("\t%d \t%d \t%d \t%s \t\t%s",
				entry.RoomNum,
				entry.StartDate,
				entry.EndDate,
				entry.ClientName,
				entry.Author))
		}

		return "", 0
	},
}

var cmdNew = cmdT{
	name: "new",
	desr: "creates new booking with specified room number and date interval (start date .. end date), if it is possible",
	exec: func() (Response string, ExitFlag int) {
		type QueryParamsT struct {
			name  string
			input string
		}

		var QueryParams = [...]QueryParamsT{
			{"Room number (0 - if any):", ""},
			{"Start date YYYYMMDD (0 - if any):", ""},
			{"End date YYYYMMDD(0 - if any):", ""},
			{"Client's name:", ""},
		}
		println("Enter bookings parameters.")
		for i := range QueryParams {
			print(QueryParams[i].name)
			_, err := fmt.Scanln(&QueryParams[i].input)
			if err != nil {
				print("Error while reading.\n")
				continue
			}
		}

		// assert params

		// jsonQuery, err := json.Marshal(payload)
		resp, err := http.Post("http://localhost:3000/newbooking",
			"application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf(
				"{\"RoomNum\": %s,\"StartDate\": %s,\"EndDate\": %s,\"Client\": \"%s\"}",
				QueryParams[0].input,
				QueryParams[1].input,
				QueryParams[2].input,
				QueryParams[3].input))))

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		BodyByteSlice, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(2)
		}

		return string(BodyByteSlice), 0
	},
}

var cmdNewUser = cmdT{
	name: "newUser",
	desr: "creates new user",
	exec: func() (Response string, ExitFlag int) {
		type QueryParamsT struct {
			name  string
			input string
		}

		var QueryParams = [...]QueryParamsT{
			{"User's login:", ""},
			{"Password:", ""},
			{"Repeat password:", ""},
			{"User's name:", ""},
			{"User's type (admin / agent):", ""},
		}
		println("Enter new user's parameters.")
		for i := range QueryParams {
			print(QueryParams[i].name)
			_, err := fmt.Scanln(&QueryParams[i].input)
			if err != nil {
				print("Error while reading.\n")
				continue
			}
		}

		// assert params
		if QueryParams[1].input != QueryParams[2].input {
			return "Passwords don't match!", 0
		}

		h := md5.New()
		io.WriteString(h, QueryParams[1].input)

		// jsonQuery, err := json.Marshal(payload)
		resp, err := http.Post("http://localhost:3000/newuser",
			"application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf(
				"{\"Login\": \"%s\",\"Password\": \"%x\",\"Name\": \"%s\",\"Type\": \"%s\"}",
				QueryParams[0].input,
				h.Sum(nil),
				QueryParams[3].input,
				QueryParams[4].input))))

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		BodyByteSlice, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(2)
		}

		return string(BodyByteSlice), 0
	},
}

// var cmdLogin = cmdT{
// 	name: "login",
// 	desr: "performs login to the system",
// 	exec: func() (Response string, ExitFlag int) {
// 		type QueryParamsT struct {
// 			name  string
// 			input string
// 		}

// 		var QueryParams = [...]QueryParamsT{
// 			{"login:", ""},
// 			{"password:", ""},
// 		}
// 		for i := range QueryParams {
// 			print(QueryParams[i].name)
// 			_, err := fmt.Scanln(&QueryParams[i].input)
// 			if err != nil {
// 				print("Error while reading.\n")
// 				continue
// 			}
// 		}

// 		// assert params

// 		// jsonQuery, err := json.Marshal(payload)
// 		resp, err := http.Post("http://localhost:3000/newbooking",
// 			"application/json",
// 			bytes.NewBuffer(fmt.Sprintf(
// 				"{\"RoomNum\": %s,\"StartDate\": %s,\"EndDate\": %s}",
// 				QueryParams[0].input,
// 				QueryParams[1].input,
// 				QueryParams[2].input)))

// 		if err != nil {
// 			fmt.Println("error:", err)
// 			os.Exit(1)
// 		}
// 		defer resp.Body.Close()

// 		BodyByteSlice, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println("error:", err)
// 			os.Exit(2)
// 		}

// 		fmt.Println(string(BodyByteSlice))

// 		return "", 0
// 	},
// }
