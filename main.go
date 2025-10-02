package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var formTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Mock Data Generator</title>
</head>
<body>
	<h2>Add Mock User</h2>
	<form method="POST" action="/mockdata">
		<label>Name:</label><br>
		<input type="text" name="name" required><br><br>
		<label>Email:</label><br>
		<input type="email" name="email" required><br><br>
		<input type="submit" value="Add Mock User">
	</form>
</body>
</html>
`

func main() {
	app := pocketbase.New()

	// Initialize the app
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	// Serve the form
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/mockdata", func(c *core.RequestContext) error {
			tmpl := template.Must(template.New("form").Parse(formTemplate))
			return tmpl.Execute(c.Response(), nil)
		})

		// Handle form submission
		e.Router.POST("/mockdata", func(c *core.RequestContext) error {
			name := c.Request().FormValue("name")
			email := c.Request().FormValue("email")

			// Create a new record
			record := &core.Record{
				Collection: "users",
				Fields: map[string]any{
					"name":  name,
					"email": email,
				},
			}

			// Save the record
			if err := app.Dao().SaveRecord(record); err != nil {
				return fmt.Errorf("failed to save record: %w", err)
			}

			return c.String(http.StatusOK, "Mock user added! <a href='/mockdata'>Add another</a>")
		})

		return nil
	})

	// Start the app
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
