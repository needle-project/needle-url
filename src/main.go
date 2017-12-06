package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/go-redis/redis"
	"strconv"
	"github.com/gorilla/mux"
	"strings"
	"os"
	"encoding/base64"
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
		user, pass, _ := r.BasicAuth()
		if user != configData.BasicAuth.Username && pass != configData.BasicAuth.Password {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		DeleteHandler(w, r, vars["pathReference"], redisClient)
	}).Methods("DELETE")

	// create item
	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != configData.BasicAuth.Username && pass != configData.BasicAuth.Password {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		createHandler(w, r, redisClient)
	}).Methods("POST")

	// handle update item
	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != configData.BasicAuth.Username && pass != configData.BasicAuth.Password {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		updateHandler(w, r, redisClient)
	}).Methods("PATCH")

	// get all items
	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != configData.BasicAuth.Username && pass != configData.BasicAuth.Password {
			http.Error(w, "Unauthorized.", 401)
			return
		}
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
	http.HandleFunc("/admin-test/", func(w http.ResponseWriter, r *http.Request) {
		authenticate(w, r, configData, true)

		pathSegments := strings.Split(string(r.URL.Path), "/")
		pathSegments = pathSegments[2:]

		// if last element is a directory, search for "index.html"
		if pathSegments[len(pathSegments) - 1] == "" {
			pathSegments[len(pathSegments) - 1] = "index.html"
		}

		filePath := configData.AdminFilePath + strings.Join(pathSegments, string(os.PathSeparator))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Page not found!", 404)
			log.Println("Requested", filePath , "which does not exists!")
			return
		}

		http.ServeFile(w, r, filePath)
	})


	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir(configData.AdminFilePath))))
	http.Handle("/", router)

	fmt.Println("Starting server on port " + strconv.Itoa(configData.Port))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(configData.Port), nil))
}

func authenticate(w http.ResponseWriter, r *http.Request, configData ConfigJson, withRequest bool) {
	if withRequest {
		w.Header().Set("WWW-Authenticate", `Basic realm="Please provide authentication details"`)
	}
	user, pass, _ := r.BasicAuth()
	if user != configData.BasicAuth.Username || pass != configData.BasicAuth.Password {
		log.Println("Could not authenticate with", user, "and", pass, "from user", r.RemoteAddr)
		http.Error(w, "Unauthorized.", 401)
		return
	}
	// if we force the "pop-up" the we are sure that is the interface
	// and we can write a cookie so we can forward the authentication
	// for the API calls
	if withRequest {
		encodedCredentials := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))

		cookie := http.Cookie{}
		cookie.Name = "btoa"
		cookie.Value = encodedCredentials
		cookie.Path = "/"
		cookie.MaxAge = 0
		cookie.Secure = false
		cookie.HttpOnly = false
		http.SetCookie(w, &cookie)
	}
}
