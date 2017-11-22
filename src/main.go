package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"github.com/go-redis/redis"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"html"
)

type ConfigJson struct {
	Port 			int 	`json:"port"`
	RedisHostname 	string 	`json:"redis_hostname"`
	RedisPort 		int 	`json:"redis_port"`
	RedisDb 		int 	`json:"redis_db"`
	DefaultRedirect	string 	`json:"default_redirect_path"`
	AdminFilePath	string 	`json:"admin_path"`
}

/**
 * Get all config data
 */
func getConfig() ConfigJson {
	rawFile, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var configJson ConfigJson
	json.Unmarshal(rawFile, &configJson)
	return configJson
}

/**
 * Create redis configuration
 */
func createRedisClient(configData ConfigJson) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     configData.RedisHostname + ":" + strconv.Itoa(configData.RedisPort),
		Password: "",
		DB:       configData.RedisDb,
		MaxRetries: 5,
	})
}

/**
 * MAIN
 */
func main() {
	var configData ConfigJson = getConfig()
	redisClient := createRedisClient(configData)

	// exclude favicon
	http.HandleFunc("/favicon.ico", func(response http.ResponseWriter, r *http.Request) {})
	// redirect Block
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		var path string = strings.Trim(html.EscapeString(request.URL.Path), "/")
		var redirectUrl string

		val, err := redisClient.Get(path).Result()
		if err != nil && err.Error() != "EOF" {
			//log.Log("No URL found for ", path, " got ", err.Error())
			log.Println("No URL found for ", path, " got ", err.Error())
			http.Redirect(response, request, configData.DefaultRedirect, 301)
		}
		redirectUrl = val

		log.Println("Redirected to ", redirectUrl)
		http.Redirect(response, request, redirectUrl, 301)
	})
	fmt.Println("Starting server on port " + strconv.Itoa(configData.Port))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(configData.Port), nil))
}
