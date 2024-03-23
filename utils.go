// utils.go

package main

import (
	"bytes"
	"html/template"
)

/* Utility Functions */

// Outputs an HTML li element for each item in the list
func generateTodoListItems(items Items) (string, error) {
	var htmlResponse bytes.Buffer

	for _, item := range items.Items {
		itemHTML, err := generateTodoItem(item)
		if err != nil {
			return "", err // Returning on the first error encountered
		}
		htmlResponse.WriteString(itemHTML)
	}

	return htmlResponse.String(), nil
}

// Outputs an HTML li element for an item
func generateTodoItem(item Item) (string, error) {
	tmpl := template.Must(template.New("todoItem").Parse(`
		<li hx-get='/items/edit/{{.ID}}' hx-swap='outerHTML' hx-indicator="#loading">
			<div class='item'>
				<div class='item-name'>{{.Name}}</div>
				<div class='item-due'>Due: {{.Due}}</div>
			</div>
			<div class='controls'>
				<button class='control-delete' hx-delete='/items/delete/{{.ID}}' hx-target='#items' hx-swap='innerHTML' hx-trigger='click consume'>-</button>
			</div>
		</li>
	`))

	var htmlResponse bytes.Buffer
	if err := tmpl.Execute(&htmlResponse, item); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

// Outputs a set of HTML input elements for an item
func generateTodoItemForm(item Item) (string, error) {
	tmpl := template.Must(template.New("todoItemForm").Parse(`
		<li>
			<form class='item-form' hx-post='/items/update/{{.ID}}' hx-swap='outerHTML' hx-target='closest li' hx-indicator="#loading">
				<div class='item'>
					<input type='text' name='item' value='{{.Name}}' required>
					<input type='date' name='date' value='{{.Due}}' required>
				</div>
				<div class='controls'>
					<button class='control-update' type='submit'>✔</button>
				</div>
			</form>
		</li>
	`))

	var htmlResponse bytes.Buffer
	if err := tmpl.Execute(&htmlResponse, item); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}
