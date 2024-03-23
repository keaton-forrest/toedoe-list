// routes.go

package main

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/* HTTP Request Handlers */

// GET /login
func loginPage(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", gin.H{})
}

// POST /login
func login(context *gin.Context) {
	// Read incoming username and password
	username := context.PostForm("username")
	password := context.PostForm("password")
	if username == "" || password == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	// Get a session for the user
	session := sessions.Default(context)

	// Load existing users
	users, err := LoadUsers()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	// Find the user in the users list
	for _, user := range users.Users {
		// If we match a username
		if user.Username == username {
			// Check the password hash
			if checkPasswordHash(password, user.Hash) {
				// Set the user's list UUID as a cookie which never expires, and is httpOnly
				session.Set("list", user.List.String())
				session.Save()
				// Send a redirect script to the client to redirect to the home page
				context.Header("Content-Type", "text/html")
				context.String(http.StatusOK, "<script>window.location.href = '/';</script>")
				return
			}
			break
		}
	}

	// If the user was not found, return a 404
	context.Status(http.StatusNotFound)
}

// GET /logout
func logout(context *gin.Context) {
	// Get a session for the user
	session := sessions.Default(context)
	// Clear the user's session
	session.Clear()
	session.Save()
	// Send a redirect script to the client to redirect to the login page
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, "<script>window.location.href = '/login';</script>")
}

// GET /register
func registerPage(context *gin.Context) {
	context.HTML(http.StatusOK, "register.html", gin.H{})
}

// POST /register
func register(context *gin.Context) {
	// Read incoming username and password
	username := context.PostForm("username")
	password := context.PostForm("password")
	if username == "" || password == "" {
		log.Println("No username or password provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Load existing users
	users, err := LoadUsers()
	if err != nil {
		log.Println("Error loading users")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Check if the username is already taken
	for _, user := range users.Users {
		if user.Username == username {
			log.Println("Username already taken")
			context.Status(http.StatusConflict)
			return
		}
	}

	// Hash the password
	hash, err := hashPassword(password)
	if err != nil {
		log.Println("Error hashing password")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Generate a new list UUID
	newListUUID := uuid.New()

	// Add the new user to the users list
	newUser := User{
		Username: username,
		Hash:     hash,
		List:     newListUUID,
	}
	users.Users = append(users.Users, newUser)

	// Save the updated users back to the JSON file
	if err := SaveUsers(users); err != nil {
		log.Println("Error saving users")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Create a new list file for the user
	if err := CreateListFileIfNotExist(newListUUID.String()); err != nil {
		log.Println("Error creating list file")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Send a redirect script to the client to redirect to the login page
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, "<script>window.location.href = '/login';</script>")
}

// GET /
func indexPage(context *gin.Context) {
	// Check if we have a session and if not redirect to the login page
	session := sessions.Default(context)
	if session.Get("list") == nil {
		context.Redirect(http.StatusFound, "/login")
		return
	}

	// Send the home page
	context.HTML(http.StatusOK, "index.html", gin.H{})
}

// GET /items
func items(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if the user's list UUID is in the session
	listUUID, err := uuid.Parse(session.Get("list").(string))
	if err != nil {
		// If the user's list UUID is not set in the cookie, 400 Bad Request
		context.Status(http.StatusBadRequest)
		return
	}

	items, err := LoadItems(listUUID.String())
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	// Sort the items by due date
	sort.Sort(items)

	htmlResponse, err := generateTodoListItems(items)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, htmlResponse)
}

// POST /items/create
func createItem(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if the user's list UUID is in the session
	listUUID, err := uuid.Parse(session.Get("list").(string))
	if err != nil {
		// If the user's list UUID is not set in the cookie, 400 Bad Request
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item name
	itemName := context.PostForm("item")
	if itemName == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item due date
	itemDue := context.PostForm("date")
	var formattedDueDate string // Variable to hold the formatted due date
	if itemDue == "" {
		// If no date is provided, default to 3 days from now, formatted as YYYY-MM-DD
		formattedDueDate = time.Now().Add(72 * time.Hour).Format("2006-01-02")
	} else {
		// Parse the provided date assuming it's in YYYY-MM-DD format
		dueDate, err := time.Parse("2006-01-02", itemDue)
		if err != nil {
			context.Status(http.StatusBadRequest)
			return
		}
		// Format the date to YYYY-MM-DD format to ensure consistency
		formattedDueDate = dueDate.Format("2006-01-02")
	}

	// Load existing items
	items, err := LoadItems(listUUID.String())
	if err != nil {
		log.Println("Error loading items")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Add the new item to the items list
	newID := uuid.New()
	newItem := Item{
		ID:     newID,
		Name:   itemName,
		Due:    formattedDueDate,
		Status: "visible",
	}
	items.Items = append(items.Items, newItem)

	// Sort the items by due date
	sort.Sort(items)

	// Save the updated items back to the JSON file
	if err := SaveItems(listUUID.String(), items); err != nil {
		log.Println("Error saving items")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Generate and send the HTML response
	htmlResponse, err := generateTodoListItems(items)
	if err != nil {
		log.Println("Error generating HTML response")
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, htmlResponse)
}

// DELETE /items/delete/:id
func deleteItem(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if the user's list UUID is in the session
	listUUID, err := uuid.Parse(session.Get("list").(string))
	if err != nil {
		// If the user's list UUID is not set in the cookie, 400 Bad Request
		log.Println("No list UUID in session")
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item ID from the URL
	itemID := context.Param("id")
	if itemID == "" {
		log.Println("No item ID provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Load existing items
	items, err := LoadItems(listUUID.String())
	if err != nil {
		log.Println("Error loading items")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Find and remove the item from the items list
	var found bool
	for i, item := range items.Items {
		if item.ID.String() == itemID {
			found = true
			// Remove the item from the list
			items.Items = append(items.Items[:i], items.Items[i+1:]...)
			break
		}
	}

	// If the item was not found, return a 404
	if !found {
		context.Status(http.StatusNotFound)
		return
	} else {
		// Save the updated items back to the JSON file
		if err := SaveItems(listUUID.String(), items); err != nil {
			log.Println("Error saving items")
			context.Status(http.StatusInternalServerError)
			return
		}
	}

	htmlResponse, err := generateTodoListItems(items)
	if err != nil {
		log.Println("Error generating HTML response")
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, htmlResponse)
}

// GET /items/edit/:id
func editItem(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if the user's list UUID is in the session
	listUUID, err := uuid.Parse(session.Get("list").(string))
	if err != nil {
		// If the user's list UUID is not set in the cookie, 400 Bad Request
		log.Println("No list UUID in session")
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item ID from the URL
	itemID := context.Param("id")
	if itemID == "" {
		log.Println("No item ID provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Load existing items
	items, err := LoadItems(listUUID.String())
	if err != nil {
		log.Println("Error loading items")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Find the item in the items list
	var found bool
	var itemToEdit Item
	for _, item := range items.Items {
		if item.ID.String() == itemID {
			found = true
			// Mark the item as "editing"
			item.Status = "editing"
			// Get a reference to the item
			itemToEdit = item
			break
		}
	}

	// If the item was not found, return a 404
	if !found {
		log.Println("Item not found")
		context.Status(http.StatusNotFound)
		return
	}

	// Send the item back to the client as a set of inputs
	htmlResponse, err := generateTodoItemForm(itemToEdit)
	if err != nil {
		log.Println("Error generating HTML response")
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, htmlResponse)
}

// POST /items/update/:id
func updateItem(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if the user's list UUID is in the session
	listUUID, err := uuid.Parse(session.Get("list").(string))
	if err != nil {
		// If the user's list UUID is not set in the cookie, 400 Bad Request
		log.Println("No list UUID in session")
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item ID from the URL
	itemID := context.Param("id")
	if itemID == "" {
		log.Println("No item ID provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item name
	itemName := context.PostForm("item")
	if itemName == "" {
		log.Println("No item name provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Read incoming item due date
	itemDue := context.PostForm("date")
	if itemDue == "" {
		log.Println("No item due date provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Load existing items
	items, err := LoadItems(listUUID.String())
	if err != nil {
		log.Println("Error loading items")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Find the item in the items list
	var found bool
	var itemToUpdate Item
	for i, item := range items.Items {
		if item.ID.String() == itemID {
			found = true
			// Update the item's name and due date
			items.Items[i].Name = itemName
			items.Items[i].Due = itemDue
			// Mark the item as "visible"
			items.Items[i].Status = "visible"
			// Get a reference to the updated item
			itemToUpdate = items.Items[i]
			break
		}
	}

	// If the item was not found, return a 404
	if !found {
		log.Println("Item not found")
		context.Status(http.StatusNotFound)
		return
	}

	// Save the updated items back to the JSON file
	if err := SaveItems(listUUID.String(), items); err != nil {
		log.Println("Error saving items")
		context.Status(http.StatusInternalServerError)
		return
	}

	htmlResponse, err := generateTodoItem(itemToUpdate)
	if err != nil {
		log.Println("Error generating HTML response")
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, htmlResponse)
}

// GET /ping
func ping(context *gin.Context) {
	context.String(http.StatusOK, "pong")
}
