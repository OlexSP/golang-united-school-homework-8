package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type userList []User

// argsParsing gets  flags' arguments from the console
func argsParsing() (args Arguments) {
	id := flag.String("id", "", "user ID")
	item := flag.String("item", "", "item information json")
	operation := flag.String("operation", "", "operation to execute")
	fileName := flag.String("fileName", "", "json file name")
	flag.Parse()
	args["id"] = *id
	args["item"] = *item
	args["operation"] = *operation
	args["fileName"] = *fileName
	return
}

// Perform  gets arguments from the console and os.Stdout stream
func Perform(args Arguments, writer io.Writer) error {

	fileName := args["fileName"]
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "":
		return errors.New("-operation flag has to be specified")
	case "list":
		return ReadUserList(fileName, writer)
	case "add":
		return addToUserList(args, writer)
	case "findById":
		return FindByID(args, writer)
	//case "remove":
	//	return removeUser(args, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}

}

// ReadUserList retrieve list from the users.json file and print it to the io.Writer stream
func ReadUserList(fileName string, writer io.Writer) error {
	//reading file into byteSlice
	byteSlice, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file reading error:%w", err)
	}
	//passing byteSlice to writer io
	_, err = writer.Write(byteSlice)
	if err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}
	return nil
}

// addToUserList adding new item to the array inside users.json file
func addToUserList(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]
	item := args["item"]
	//  empty item check
	if item == "" {
		return errors.New("-item flag has to be specified")
	}
	// unmarshal item into User{}
	user1 := User{}
	err := json.Unmarshal([]byte(item), &user1)
	if err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}
	// simple validator
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matches := pattern.MatchString(user1.Email)
	if user1.Id == "" || !matches || user1.Age < 6 {
		return errors.New("item data is not valid")
	}

	// read user list file & unmarshal it to the userList{}
	userList1 := userList{}
	byteSlice, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file reading error:%w", err)
	}
	err = json.Unmarshal(byteSlice, &userList1)
	if err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}

	// ID duplicate check
	for _, v := range userList1 {
		if v.Id == user1.Id {
			return fmt.Errorf("Item with id %s already exists", user1.Id)
		}
	}

	// adding item to the file
	//userList1 = append(userList1, user1)

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(item)); err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}
	return nil
}

func FindByID(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}
	// read user list file & unmarshal it to the userList{}
	userList1 := userList{}
	byteSlice, err := os.ReadFile(args["fileName"])
	if err != nil {
		return fmt.Errorf("file reading error:%w", err)
	}
	err = json.Unmarshal(byteSlice, &userList1)
	if err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}
	// find user
	var userBiteSlice []byte
	for _, user := range userList1 {
		if user.Id == args["id"] {
			userBiteSlice, err = json.Marshal(user)
			break
		}
	}
	//passing byteSlice to writer io
	_, err = writer.Write(userBiteSlice)
	if err != nil {
		return fmt.Errorf("an error has occurred:%w", err)
	}
	return nil

}

//func removeUser(args, writer) error {
//
//	return nil
//}

func main() {
	a := argsParsing()
	err := Perform(a, os.Stdout)
	if err != nil {
		panic(err)
	}
}
