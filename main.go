package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Ensure "users" collection exists
		collection, err := app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			collection = core.NewBaseCollection("users")
			collection.Schema = core.Schema{
				&core.SchemaField{Name: "name", Type: "text"},
				&core.SchemaField{Name: "email", Type: "email"},
			}
			if err := app.Dao().SaveCollection(collection); err != nil {
				return err
			}
			log.Println("Created 'users' collection")
		}

		// Custom route: GET = show form, POST = save mock data
		e.Router.AddRoute("/mockdata", func(c http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				tmpl := template.Must(template.New("form").Parse(formTemplate))
				tmpl.Execute(c, nil)

			case http.MethodPost:
				name := r.FormValue("name")
				email := r.FormValue("email")

				record := core.NewRecord(collection)
				record.Set("name", name)
				record.Set("email", email)

				if err := app.Dao().SaveRecord(record); err != nil {
					http.Error(c, "Failed to save record", http.StatusInternalServerError)
					return
				}

				fmt.Fprintf(c, "✅ Mock user added!<br><a href=\"/mockdata\">Add another</a>")
			}
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
