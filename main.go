package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/daos"
    "github.com/pocketbase/pocketbase/models"
    "github.com/pocketbase/pocketbase/models/schema"
)

var formTemplate = `
<!DOCTYPE html>
<html>
<head>
  <title>Mock Data Entry</title>
</head>
<body>
  <h2>Add Mock User</h2>
  <form method="POST" action="/mockdata">
    <label>Name:</label><br>
    <input type="text" name="name" required><br><br>
    <label>Email:</label><br>
    <input type="email" name="email" required><br><br>
    <input type="submit" value="Add">
  </form>
</body>
</html>
`

func main() {
    app := pocketbase.New()

    app.OnBeforeServe().Add(func(e *pocketbase.ServeEvent) error {
        // Try to find “users” collection
        coll, err := daos.New(app).FindCollectionByNameOrId("users")
        if err != nil {
            // If not exists, create it
            coll = &models.Collection{
                Name: "users",
                Type: models.CollectionTypeBase,
                Schema: schema.NewSchema(
                    &schema.SchemaField{Name: "name", Type: schema.FieldTypeText},
                    &schema.SchemaField{Name: "email", Type: schema.FieldTypeEmail},
                ),
            }
            if err := daos.New(app).SaveCollection(coll); err != nil {
                return fmt.Errorf("failed to create collection: %w", err)
            }
            log.Println("Created collection users")
        }

        // GET endpoint to show the form
        e.Router.GET("/mockdata", func(c *apis.RequestContext) error {
            tmpl := template.Must(template.New("form").Parse(formTemplate))
            return tmpl.Execute(c.Response(), nil)
        })

        // POST endpoint to receive form and save record
        e.Router.POST("/mockdata", func(c *apis.RequestContext) error {
            name := c.Request().FormValue("name")
            email := c.Request().FormValue("email")

            rec := models.NewRecord(coll)
            rec.Set("name", name)
            rec.Set("email", email)

            if err := daos.New(app).SaveRecord(rec); err != nil {
                return c.String(http.StatusInternalServerError, "Error saving record: "+err.Error())
            }

            return c.String(http.StatusOK, "Mock user added! <a href=\"/mockdata\">Add more</a>")
        })

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
