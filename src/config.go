package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

// Default ConfigObject
type ConfigJson struct {
	Port            int                     `json:"port"`
	RedisHostname   string                  `json:"redis_hostname"`
	RedisPassword   string                  `json:"redis_password"`
	RedisPort       int                     `json:"redis_port"`
	RedisDb         int                     `json:"redis_db"`
	DefaultRedirect string                  `json:"default_redirect_path"`
	AdminFilePath   string                  `json:"admin_path"`
	BasicAuth       BasicAuthCredentials    `json:"basic_auth"`
}

// Should map BasicAuth parameters
// It should be used as a part of ConfigJson
type BasicAuthCredentials struct {
	Username	string 	`json:"username"`
	Password	string 	`json:"password"`
}

// Fetch configuration data
func getConfig() ConfigJson {
	rawFile, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var configJson ConfigJson
	// fill with defaults that will be overridden by configJson
	configJson.Port = 80
	configJson.RedisHostname = "127.0.0.1"
	configJson.RedisPort = 6379
	configJson.RedisDb = 0

	configJson.AdminFilePath = ""

	json.Unmarshal(rawFile, &configJson)
	return configJson
}
