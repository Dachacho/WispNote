# WispNotes

A tiny note-taking app built with **Go**, **SQLite**, **htmx**, and **Pico.css**.  
This is just a fun project to get familiar with Go web servers and htmx for instant UI updates.

## Stack

- **Backend:** Go (net/http), SQLite (via mattn/go-sqlite3)
- **Frontend:** htmx (for AJAXy stuff), Pico.css (for quick clean styles)
- **Templates:** Go's html/template

## What it does

- Lets you add notes (they show up instantly, no reload)
- Lets you delete notes (they disappear instantly)
- Stores everything in a local SQLite file

## How to use

1. Clone this repo.
2. Run `go mod tidy` if you want.
3. Start it up:
   ```sh
   go run main.go
   ```
4. Open [http://localhost:8080](http://localhost:8080) in your browser.
5. Type a note, hit "Add Note", and watch it appear.
6. Click "Delete" to remove a note.

## Why?

Just messing around with Go and htmx.  
If you want to learn how to build a simple, modern-feeling web app with Go and no JS (we hate js here), this is a good place to start.

---

MIT License.
