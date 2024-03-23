// file_handlers.go

package main

import (
	"encoding/json"
	"log"
	"os"
)

/* File Handlers */

// LoadItems loads the items from the JSON file associated with the specified list UUID
func LoadItems(listID string) (Items, error) {
	var items Items
	filename := listID + ".json" // Construct the filename using the list UUID
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return items, err
	}
	err = json.Unmarshal(data, &items)
	return items, err
}

// SaveItems saves the items to the JSON file associated with the specified list UUID
func SaveItems(listID string, items Items) error {
	filename := listID + ".json" // Construct the filename using the list UUID
	data, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		log.Println(err)
		return err
	}
	// Write the data to the file named {listID}.json inside the current directory with 0644 permissions
	return os.WriteFile(filename, data, 0644)
}

// LoadUsers loads the users from the JSON file
func LoadUsers() (Users, error) {
	var users Users
	data, err := os.ReadFile("users.json")
	if err != nil {
		log.Println(err)
		return users, err
	}
	err = json.Unmarshal(data, &users)
	return users, err
}

func SaveUsers(users Users) error {
	data, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		log.Println(err)
		return err
	}
	// Write the data to the file named users.json inside the current directory with 0644 permissions
	return os.WriteFile("users.json", data, 0644)
}

// Create a new list file with the specified list UUID if it doesn't exist
func CreateListFileIfNotExist(listID string) error {
	filename := listID + ".json"

	// Write {} to the file to represent an empty json file
	err := os.WriteFile(filename, []byte("{}"), 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// If the file already exists or has been created successfully, return nil
	return nil
}
