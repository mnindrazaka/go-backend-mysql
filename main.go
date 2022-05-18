package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

type AlbumRequest struct {
	Title  string
	Artist string
	Price  float32
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func handleAlbums(w http.ResponseWriter, r *http.Request) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album")

	if err != nil {
		w.Write([]byte("Failed to get data"))
	}

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			w.Write([]byte("Failed to get data"))
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		w.Write([]byte("Failed to get data"))
	}

	json.NewEncoder(w).Encode(albums)
}

func handleAlbumDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var album Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			w.Write([]byte("No album found"))
		} else {
			w.Write([]byte("Failed to get data"))
		}
	}

	json.NewEncoder(w).Encode(album)
}

func handleAlbumCreate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var albumRequest AlbumRequest
	json.Unmarshal(reqBody, &albumRequest)

	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", albumRequest.Title, albumRequest.Artist, albumRequest.Price)

	if err != nil {
		w.Write([]byte("Failed to create data"))
	}

	w.Write([]byte("Success to create data"))
}

func handleAlbumUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var albumRequest AlbumRequest
	json.Unmarshal(reqBody, &albumRequest)

	_, err := db.Exec("Update album SET title=?, artist=?, price=? WHERE id=?", albumRequest.Title, albumRequest.Artist, albumRequest.Price, id)

	if err != nil {
		w.Write([]byte("Failed to update data"))
	}

	w.Write([]byte("Success to update data"))
}

func handleAlbumDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM album WHERE id=?", id)

	if err != nil {
		w.Write([]byte("Failed to delete data"))
	}

	w.Write([]byte("Success to delete data"))
}

var db *sql.DB

func main() {
	mysqlConfig := mysql.Config{
		User:   "root",
		Passwd: "roottoor",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", mysqlConfig.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handleHome)
	router.HandleFunc("/albums", handleAlbumCreate).Methods("POST")
	router.HandleFunc("/albums", handleAlbums)
	router.HandleFunc("/albums/{id}", handleAlbumUpdate).Methods("PUT")
	router.HandleFunc("/albums/{id}", handleAlbumDelete).Methods("DELETE")
	router.HandleFunc("/albums/{id}", handleAlbumDetail)

	http.ListenAndServe(":3000", router)
}
