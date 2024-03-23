// types.go

package main

import (
	"time"

	"github.com/google/uuid"
)

/* Types */

type Item struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Due    string    `json:"due"`
	Status string    `json:"status"`
}

type Items struct {
	Items []Item `json:"items"`
}

type User struct {
	Username string    `json:"username"` // This is an email
	Hash     string    `json:"hash"`
	List     uuid.UUID `json:"list"`
}

type Users struct {
	Users []User `json:"users"`
}

// Implement sort.Interface for Items based on the Due field.
// This will allow us to use the sort.Sort function to sort the items ascending by due date.
// We need to implement the Len, Less, and Swap methods.

func (items Items) Len() int {
	return len(items.Items)
}

func (items Items) Less(i, j int) bool {
	// Parse the Due field to time.Time
	timeI, err := time.Parse("2006-01-02", items.Items[i].Due)
	if err != nil {
		return false
	}

	timeJ, err := time.Parse("2006-01-02", items.Items[j].Due)
	if err != nil {
		return false
	}

	// Compare the parsed times
	return timeI.Before(timeJ)
}

func (items Items) Swap(i, j int) {
	items.Items[i], items.Items[j] = items.Items[j], items.Items[i]
}
