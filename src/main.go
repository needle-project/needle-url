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
	"bytes"
)

type ConfigJson struct {
	Port 			int 	`json:"port"`
	RedisHostname 	string 	`json:"redis_hostname"`
	RedisPort 		int 	`json:"redis_port"`
	RedisDb 		int 	`json:"redis_db"`
	DefaultRedirect	string 	`json:"default_redirect_path"`
	AdminFilePath	string 	`json:"admin_path"`
}

type UrlItem struct {
	FromUrl		string	`json:"from_url"`
	ToUrl		string	`json:"to_url"`
}

type ListResponse struct {
	Count 			int 	`json:"total"`
	List 			[]UrlItem
}

type SuccessResponse struct {
	Status 		string   `json:"status"`
	Created 	UrlItem  `json:"url"`
}

type ErrorResponse struct {
	Status 		string   `json:"status"`
	Response 	string   `json:"message"`
}

func createErrorResponse(message string) ErrorResponse {
	var responseMessage ErrorResponse
	responseMessage.Status = "error"
	responseMessage.Response = message
	return responseMessage
}

func createSuccessResponse(item UrlItem) SuccessResponse {
	var responseMessage SuccessResponse
	responseMessage.Status = "ok"
	responseMessage.Created = item
	return responseMessage
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
	var configData = getConfig()
	redisClient := createRedisClient(configData)

	// serve static - admin interface
	fs := http.FileServer(http.Dir(configData.AdminFilePath))
	http.Handle("/admin/", http.StripPrefix("/admin/", fs))

	// exclude favicon
	http.HandleFunc("/favicon.ico", func(response http.ResponseWriter, r *http.Request) {})
	// redirect Block
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		var path = strings.Trim(html.EscapeString(request.URL.Path), "/")
		var redirectUrl string

		val, err := redisClient.Get(path).Result()
		if err != nil && err.Error() != "EOF" {
			log.Println("[INFO] No URL found for <", path, "> got ", err.Error())
			http.Redirect(response, request, configData.DefaultRedirect, 307)
			return
		}
		redirectUrl = val

		log.Println("[INFO] Redirected <", path, "> to ", redirectUrl)
		http.Redirect(response, request, redirectUrl, 301)
		return
	})
	// redirect Block
	http.HandleFunc("/url", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch request.Method {
		case http.MethodPost:
			var urlItem UrlItem
			if request.Body == nil {
				response.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(response).Encode(createErrorResponse("Please provide a POST message with `from_url` and `to_url`"))
				return
			}

			err := json.NewDecoder(request.Body).Decode(&urlItem)
			if err != nil {
				buf := new(bytes.Buffer)
				buf.ReadFrom(request.Body)
				bodyString := buf.String()

				response.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(response).Encode(json.NewEncoder(response).Encode(createErrorResponse("Could not process request. Invalid json received, got <" + bodyString + ">")))
				return
			}

			duplicateValue, duplicateError := redisClient.Get(urlItem.FromUrl).Result()
			if duplicateError != redis.Nil && duplicateValue != "" {
				response.WriteHeader(http.StatusConflict)
				json.NewEncoder(response).Encode(createErrorResponse("A route for <" + urlItem.FromUrl + "> already exists!"))
				return
			}

			var listKeys []string
			availableKeys, err := redisClient.Get("list_keys").Result()
			if err != nil {
				availableKeys = ""
				log.Println("Could not retrieve keys, got ", err)
			}
			json.Unmarshal([]byte(availableKeys), &listKeys)
			listKeys = append(listKeys, urlItem.ToUrl)

			updateKeys, _encodeError := json.Marshal(listKeys)
			if _encodeError != nil {
				response.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(response).Encode(createErrorResponse("Unexpected error has occurred when saving list data!"))
				log.Println(_encodeError)
				return
			}
			setError := redisClient.Set("list_keys", updateKeys, 0).Err()
			if setError != nil {
				response.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(response).Encode(createErrorResponse("Unexpected error has occurred when saving list data!"))
				log.Println(setError)
				return
			}
			writeError := redisClient.Set(urlItem.FromUrl, urlItem.ToUrl, 0).Err()
			if writeError != nil {
				response.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(response).Encode(createErrorResponse("Unexpected error has occurred when saving list data!"))
				log.Println(writeError)
				return
			}

			response.WriteHeader(http.StatusCreated)
			json.NewEncoder(response).Encode(createSuccessResponse(urlItem))
			return
		case http.MethodPatch:
			var urlItem UrlItem
			if request.Body == nil {
				response.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(response).Encode(createErrorResponse("Please provide a PATCH message with `from_url` and `to_url`"))
				return
			}

			err := json.NewDecoder(request.Body).Decode(&urlItem)
			if err != nil {
				buf := new(bytes.Buffer)
				buf.ReadFrom(request.Body)
				bodyString := buf.String()

				response.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(response).Encode(json.NewEncoder(response).Encode(createErrorResponse("Could not process request. Invalid json received, got <" + bodyString + ">")))
				return
			}

			duplicateValue, duplicateError := redisClient.Get(urlItem.FromUrl).Result()
			if duplicateError == redis.Nil && duplicateValue == "" {
				response.WriteHeader(http.StatusConflict)
				json.NewEncoder(response).Encode(createErrorResponse("Could not find route for <" + urlItem.FromUrl + ">!"))
				return
			}

			writeError := redisClient.Set(urlItem.FromUrl, urlItem.ToUrl, 0).Err()
			if writeError != nil {
				response.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(response).Encode(createErrorResponse("Unexpected error has occurred when saving list data!"))
				log.Println(writeError)
				return
			}
			response.WriteHeader(http.StatusOK)
			json.NewEncoder(response).Encode(createSuccessResponse(urlItem))
			return
		case http.MethodDelete:
			var item = strings.Trim(html.EscapeString(request.URL.Path), "/url")
			fmt.Println(item)/*
			var urlItem UrlItem
			if request.Body == nil {
				response.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(response).Encode(createErrorResponse("Please provide a PATCH message with `from_url` and `to_url`"))
				return
			}
			err := json.NewDecoder(request.Body).Decode(&urlItem)
			if err != nil {
				buf := new(bytes.Buffer)
				buf.ReadFrom(request.Body)
				bodyString := buf.String()

				response.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(response).Encode(json.NewEncoder(response).Encode(createErrorResponse("Could not process request. Invalid json received, got <" + bodyString + ">")))
				return
			}

			duplicateValue, duplicateError := redisClient.Get(urlItem.FromUrl).Result()
			if duplicateError == redis.Nil && duplicateValue == "" {
				response.WriteHeader(http.StatusConflict)
				json.NewEncoder(response).Encode(createErrorResponse("Could not find route for <" + urlItem.FromUrl + ">!"))
				return
			}

			writeError := redisClient.Set(urlItem.FromUrl, urlItem.ToUrl, 0).Err()
			if writeError != nil {
				response.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(response).Encode(createErrorResponse("Unexpected error has occurred when saving list data!"))
				log.Println(writeError)
				return
			}
			response.WriteHeader(http.StatusOK)
			json.NewEncoder(response).Encode(createSuccessResponse(urlItem))*/
			return
		case http.MethodGet:
			// establish base items
			limit := request.URL.Query().Get("limit")
			if limit == "" {
				limit = "20"
			}
			offset := request.URL.Query().Get("offset")
			if offset == "" {
				offset = "0"
			}
			var fromLimit, offsetConversionError = strconv.Atoi(offset)
			if offsetConversionError != nil {
				log.Println("Could not convert offset request, got ", offsetConversionError)
				fromLimit = 0
			}
			var toLimit, limitConversionError = strconv.Atoi(limit)
			if limitConversionError != nil {
				log.Println("Could not convert limit request, got ", limitConversionError)
				toLimit = 20
			}
			toLimit = toLimit + fromLimit

			var listKeys []string
			availableKeys, err := redisClient.Get("list_keys").Result()
			if err != nil {
				availableKeys = ""
				log.Println("Could not retrieve keys, got ", err)
			}
			json.Unmarshal([]byte(availableKeys), &listKeys)

			var list []UrlItem
			for i := fromLimit; i < toLimit; i++ {
				var newUrlItem UrlItem
				newUrlItem.FromUrl = listKeys[i]

				// get to URL
				toUrlValue, _err := redisClient.Get(listKeys[i]).Result()
				if _err == nil {
					newUrlItem.ToUrl = toUrlValue
				}
				list = append(list, newUrlItem)
			}

			var listResponse ListResponse
			listResponse.Count = len(listKeys)
			listResponse.List = list

			response.WriteHeader(http.StatusOK)
			json.NewEncoder(response).Encode(listResponse)
			return
		}
	})

	fmt.Println("Starting server on port " + strconv.Itoa(configData.Port))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(configData.Port), nil))
}
