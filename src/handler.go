package main

import (
	"log"
	"net/http"
	"github.com/go-redis/redis"
	"encoding/json"
	"bytes"
	"strconv"
	"hash/adler32"
	"time"
)

// Main form of a ShortUrl item
type UrlItem struct {
	FromUrl string  `json:"from_url"`
	ToUrl   string  `json:"to_url"`
}

// handle redirect
func RedirectHandler(response http.ResponseWriter, request *http.Request, path string, redisClient *redis.Client, configData ConfigJson) {
	var redirectUrl string

	val, err := redisClient.Get(path).Result()
	if err != nil && err.Error() != "EOF" {
		log.Println("[INFO] No URL found for ", path, " got ", err.Error())
		http.Redirect(response, request, configData.DefaultRedirect, http.StatusTemporaryRedirect)
		return
	}
	redirectUrl = val

	log.Println("[INFO] Redirected ", path, " to ", redirectUrl)
	http.Redirect(response, request, redirectUrl, http.StatusMovedPermanently)
	return
}

// Handle deletion of items
func DeleteHandler(response http.ResponseWriter, request *http.Request, path string, redisClient *redis.Client) {
	_, err := redisClient.Get(path).Result()
	if err != nil && err.Error() != "EOF" {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(buildErrorResponse("Could not find route for " + path + "!"))
		return
	}

	_, delError := redisClient.Del(path).Result()
	if delError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Could not delete current item!"))
		log.Println(delError)
		return
	}

	// remove from list_keys
	var listKeys []string
	availableKeys, err := redisClient.Get("list_keys").Result()
	if err != nil {
		availableKeys = ""
		log.Println("Could not retrieve keys, got ", err)
	}
	json.Unmarshal([]byte(availableKeys), &listKeys)

	// remove from general list
	for index, key := range listKeys {
		if key == path {
			listKeys = append(listKeys[:index], listKeys[index+1:]...)
			break
		}
	}
	updateKeys, _encodeError := json.Marshal(listKeys)
	if _encodeError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Unexpected error has occurred when saving list data!"))
		log.Println(_encodeError)
		return
	}
	setError := redisClient.Set("list_keys", updateKeys, 0).Err()
	if setError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Unexpected error has occurred when saving list data!"))
		log.Println(setError)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(buildSuccessResponse(path + " deleted with success!"))
	return
}

// Handle create resource
func createHandler(response http.ResponseWriter, request *http.Request, redisClient *redis.Client) {
	var urlItem UrlItem
	if request.Body == nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(response).Encode(buildErrorResponse("Please provide a POST message with `from_url` and `to_url`"))
		return
	}

	err := json.NewDecoder(request.Body).Decode(&urlItem)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(request.Body)
		bodyString := buf.String()

		response.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(response).Encode(json.NewEncoder(response).Encode(buildErrorResponse("Could not process request. Invalid json received, got " + bodyString)))
		return
	}

	if urlItem.ToUrl == "" {
		response.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(response).Encode(buildErrorResponse("No destination URL provided!"))
		return
	}

	if urlItem.FromUrl == "" {
		for {
			var pendingGeneratedToUrl = generateHash()
			_, duplicateError := redisClient.Get(pendingGeneratedToUrl).Result()
			if duplicateError == redis.Nil || duplicateError == nil {
				urlItem.FromUrl = pendingGeneratedToUrl
				break
			}
		}
	}

	duplicateValue, duplicateError := redisClient.Get(urlItem.FromUrl).Result()
	if duplicateError != redis.Nil && duplicateValue != "" {
		response.WriteHeader(http.StatusConflict)
		json.NewEncoder(response).Encode(buildErrorResponse("A route for " + urlItem.FromUrl + " already exists!"))
		return
	}

	var listKeys []string
	availableKeys, err := redisClient.Get("list_keys").Result()
	if err != nil {
		availableKeys = ""
		log.Println("Could not retrieve keys, got", err)
	}
	json.Unmarshal([]byte(availableKeys), &listKeys)
	listKeys = append(listKeys, urlItem.FromUrl)

	updateKeys, _encodeError := json.Marshal(listKeys)
	if _encodeError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Unexpected error has occurred when saving list data!"))
		log.Println(_encodeError)
		return
	}
	setError := redisClient.Set("list_keys", updateKeys, 0).Err()
	if setError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Unexpected error has occurred when saving list data!"))
		log.Println(setError)
		return
	}
	writeError := redisClient.Set(urlItem.FromUrl, urlItem.ToUrl, 0).Err()
	if writeError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Unexpected error has occurred when saving list data!"))
		log.Println(writeError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(buildCreateResponse(urlItem))
	return
}

// Update Resource
func updateHandler(response http.ResponseWriter, request *http.Request, redisClient *redis.Client) {
	var urlItem UrlItem
	if request.Body == nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(response).Encode(buildErrorResponse("Please provide a PATCH message with `from_url` and `to_url`"))
		return
	}

	err := json.NewDecoder(request.Body).Decode(&urlItem)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(request.Body)
		bodyString := buf.String()

		response.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(response).Encode(json.NewEncoder(response).Encode(buildErrorResponse("Could not process request. Invalid json received, got " + bodyString)))
		return
	}

	duplicateValue, duplicateError := redisClient.Get(urlItem.FromUrl).Result()
	if duplicateError == redis.Nil && duplicateValue == "" {
		response.WriteHeader(http.StatusConflict)
		json.NewEncoder(response).Encode(buildErrorResponse("Could not find route for " + urlItem.FromUrl + "!"))
		return
	}

	writeError := redisClient.Set(urlItem.FromUrl, urlItem.ToUrl, 0).Err()
	if writeError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(buildErrorResponse("Unexpected error has occurred when saving list data!"))
		log.Println(writeError)
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(buildCreateResponse(urlItem))
	return
}

// Fetch list
func fetchHandler(response http.ResponseWriter, request *http.Request, redisClient *redis.Client) {
	// establish base items
	limit := request.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}
	offset := request.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	// if we ask for a specific item
	item := request.URL.Query().Get("item")

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

	// reverse the list - FIFO
	listKeys = reverseList(listKeys)

	// if we ask only for one item
	if item != "" {
		listKeys = nil
		listKeys = append(listKeys, item)
	}

	if (len(listKeys) - 1) < toLimit {
		toLimit = len(listKeys)
	}
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

// Reverse a list showing the last item as first
func reverseList(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseList(input[1:]), input[0])
}

// Generate a random HASH
func generateHash() string {
	adler32Int := adler32.Checksum([]byte(time.Now().String()))
	return strconv.FormatUint(uint64(adler32Int), 16)
}