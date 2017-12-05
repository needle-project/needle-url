package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/go-redis/redis"
	"strconv"
	"github.com/gorilla/mux"
)

/**
 * MAIN
 */
func main() {
	var configData = getConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     configData.RedisHostname + ":" + strconv.Itoa(configData.RedisPort),
		Password: configData.RedisPassword,
		DB:       configData.RedisDb,
		MaxRetries: 5,
	})

	router := mux.NewRouter()
	// exclude favicon request
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

	// delete item
	router.HandleFunc("/url/{pathReference}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		DeleteHandler(w, r, vars["pathReference"], redisClient)
	}).Methods("DELETE")

	// create item
	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		createHandler(w, r, redisClient)
	}).Methods("POST")

	// handle update item
	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		updateHandler(w, r, redisClient)
	}).Methods("PATCH")

	// get all items
	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fetchHandler(w, r, redisClient)
	}).Methods("GET")

	// redirect zone
	router.HandleFunc("/{pathReference}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		RedirectHandler(w, r, vars["pathReference"], redisClient, configData)
	}).Methods("GET")

	// admin interface
	// serve static - admin interface
	fs := http.FileServer(http.Dir(configData.AdminFilePath))
	http.Handle("/admin/", http.StripPrefix("/admin/", fs))
	http.Handle("/", router)

	fmt.Println("Starting server on port " + strconv.Itoa(configData.Port))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(configData.Port), nil))
}
