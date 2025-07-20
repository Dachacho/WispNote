package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
func initDB(){
	var err error

	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	_, err = db.Exec(`
	        CREATE TABLE IF NOT EXISTS notes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL
        );
	`)

	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}
}

type Note struct{
	ID int
	Content string
}

func getNotes() ([]Note, error){
	rows, err := db.Query("SELECT id, content from notes ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []Note
	for rows.Next(){
		var n Note
		if err := rows.Scan(&n.ID, &n.Content); err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}
	return notes, nil
}

func notesHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodGet:
		notes, err := getNotes()
		if err != nil {
			http.Error(w, "Couldn't load notes: ", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil{
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		content := r.FormValue("content")
		if content == ""{
			http.Error(w, "Missing note content", http.StatusBadRequest)
			return
		}

		id, err := addNote(content)
		if err != nil{
			http.Error(w, "Failed to add note", http.StatusInternalServerError)
			return 
		}

		tmpl := template.Must(template.ParseFiles("templates/note.html"))
		tmpl.Execute(w, Note{ID: id, Content: content})
	
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing note", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("DELETE FROM notes where id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete a note", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}
}

func addNote(content string) (int, error) {
	res, err := db.Exec("INSERT INTO notes (content) VALUES (?)", content)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	notes, err := getNotes()
	if err != nil {
		http.Error(w, "Failed to lead notes", http.StatusInternalServerError)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct{
		Notes []Note
	}{
		Notes : notes,
	}

	tmpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/", indexHandler)

	initDB()

	http.HandleFunc("/notes", notesHandler) 

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}