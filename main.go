package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
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

	// Create collection on app start
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		_, err := app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			coll := &models.Collection{
				Name: "users",
				Type: models.CollectionTypeBase,
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name:     "name",
						Type:     schema.FieldTypeText,
						Required: true,
					},
					&schema.SchemaField{
						Name:     "email",
						Type:     schema.FieldTypeEmail,
						Required: true,
						Unique:   true,
					},
				),
			}

			if err := app.Dao().SaveCollection(coll); err != nil {
				return fmt.Errorf("failed to create collection: %w", err)
			}
			log.Println("Created collection: users")
		}
		return nil
	})

	// Serve the form and handle submission
	http.HandleFunc("/mockdata", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.New("form").Parse(formTemplate))
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			name := r.FormValue("name")
			email := r.FormValue("email")

			coll, err := app.Dao().FindCollectionByNameOrId("users")
			if err != nil {
				http.Error(w, "Collection not found", http.StatusInternalServerError)
				return
			}

			rec := app.Dao().NewRecord(coll)
			rec.Set("name", name)
			rec.Set("email", email)

			if err := app.Dao().SaveRecord(rec); err != nil {
				http.Error(w, "Failed to save record: "+err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "✅ Mock user added! <a href='/mockdata'>Add another</a>")
			return
		}
	})

	// Start PocketBase
	go func() {
		log.Println("PocketBase starting...")
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Web server running on http://localhost:8080/mockdata")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
