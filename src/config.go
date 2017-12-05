package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type ConfigJson struct {
	Port 			int 	`json:"port"`
	RedisHostname 	string 	`json:"redis_hostname"`
	RedisPassword	string	`json:"redis_password"`
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
