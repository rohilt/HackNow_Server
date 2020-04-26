package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// requestBody, err := json.Marshal(map[string]string{
	// 	"coordinates": "-117.17282,32.71204;-117.17288,32.71225",
	// })

	requestBody := "coordinates=-117.17282,32.71204;-117.17288,32.71225"

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

}
