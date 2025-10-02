package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
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
	// Start PocketBase in the background
	go func() {
		cmd := exec.Command("pocketbase", "serve")
		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to start PocketBase:", err)
		}
	}()

	http.HandleFunc("/mockdata", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.New("form").Parse(formTemplate))
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			name := r.FormValue("name")
			email := r.FormValue("email")

			// Prepare JSON payload
			payload := map[string]interface{}{
				"name":  name,
				"email": email,
			}
			data, _ := json.Marshal(payload)

			// POST to PocketBase REST API
			resp, err := http.Post("http://127.0.0.1:8090/api/collections/users/records", "application/json", bytes.NewBuffer(data))
			if err != nil {
				http.Error(w, "Failed to insert record: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				http.Error(w, "PocketBase API error", resp.StatusCode)
				return
			}

			fmt.Fprintf(w, "✅ Mock user added! <a href='/mockdata'>Add another</a>")
			return
		}
	})

	fmt.Println("Server running on http://localhost:8080/mockdata")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
