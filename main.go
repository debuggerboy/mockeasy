package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
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

	app.OnBeforeServe().Add(func(e *pocketbase.ServeEvent) error {
		// ensure "users" collection exists
		collection, err := app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			collection = &models.Collection{
				Name: "users",
				Type: models.CollectionTypeBase,
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name: "name",
						Type: schema.FieldTypeText,
					},
					&schema.SchemaField{
						Name: "email",
						Type: schema.FieldTypeEmail,
					},
				),
			}
			if err := app.Dao().SaveCollection(collection); err != nil {
				return err
			}
			log.Println("Created 'users' collection")
		}

		// custom GET/POST route
		e.Router.GET("/mockdata", func(c *apis.RequestContext) error {
			tmpl := template.Must(template.New("form").Parse(formTemplate))
			return tmpl.Execute(c.Response(), nil)
		})

		e.Router.POST("/mockdata", func(c *apis.RequestContext) error {
			name := c.Request().FormValue("name")
			email := c.Request().FormValue("email")

			record := models.NewRecord(collection)
			record.Set("name", name)
			record.Set("email", email)

			if err := app.Dao().SaveRecord(record); err != nil {
				return c.String(http.StatusInternalServerError, "❌ Failed to save record")
			}

			return c.String(http.StatusOK, "✅ Mock user added! <a href='/mockdata'>Add another</a>")
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
