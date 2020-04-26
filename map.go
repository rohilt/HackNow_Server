package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func getCoordinates(address string) string {
	resp, err := http.Get("https://api.mapbox.com/geocoding/v5/mapbox.places/" + address + ".json?access_token=pk.eyJ1IjoiMTVkYW5pMSIsImEiOiJjazlmNWdvdG4wMGVvM2xubjdqcTducXM1In0.H7cu4oj3nkFtR23KeFEliQ")

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string][](map[string](map[string][]float64))
	err = json.Unmarshal(body, &result)
	// var result2 map[string]interface{}
	// err = json.Unmarshal(result, &result2)
	// log.Println(string(body))
	log.Println(result["features"][0]["geometry"]["coordinates"])
	lat := strconv.FormatFloat(result["features"][0]["geometry"]["coordinates"][0], 'f', -1, 64)
	//lat := (result["features"][0]["geometry"]["coordinates"][0]).string
	long := strconv.FormatFloat(result["features"][0]["geometry"]["coordinates"][1], 'f', -1, 64)
	//long := (result["features"][0]["geometry"]["coordinates"][1]).string
	return lat + "," + long
}

func main() {
	requestBody := "coordinates=" + getCoordinates("Land%20O%20Lakes") + ";" + getCoordinates("Orlando")

	//requestBody := "coordinates=-117.17282,32.71204;-117.17288,32.71225"
	// log.Println(requestBody)

	// var result2 map[string]interface{}
	// err = json.Unmarshal(result, &result2)
	// log.Println(string(body))
	// log.Println(result["features"][0]["geometry"]["coordinates"])

	// requestBody, err := json.Marshal(map[string]string{
	// 	"coordinates"+"="+getCoordinates("Land%20O%20Lakes") + ";" + getCoordinates("Orlando"),
	// // })
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	resp, err := http.Post("https://api.mapbox.com/directions/v5/mapbox/driving?access_token=pk.eyJ1IjoiMTVkYW5pMSIsImEiOiJjazlmNWdvdG4wMGVvM2xubjdqcTducXM1In0.H7cu4oj3nkFtR23KeFEliQ", "application/x-www-form-urlencoded", bytes.NewBufferString(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	// log.Println(resp)
}
