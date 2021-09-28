package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type resp struct {
	Body 	body 		`json:"body,omitempty"`
}

type body struct {
	Stores	[]store
}

type store struct {
	StoreName         string            `json:"storeName,omitempty"`
	Storelistnumber   int               `json:"storelistnumber,omitempty"`
	PartsAvailability map[string]partsAvailability `json:"partsAvailability,omitempty"`
}

type partsAvailability struct {
	PickupSearchQuote     string `json:"pickupSearchQuote,omitempty"`
	PickupQuote           string `json:"pickupQuote,omitempty"`
	StoreSelectionEnabled bool `json:"storeSelectionEnabled,omitempty"`
	PickupDisplay         string `json:"pickupDisplay,omitempty"`
}

func main() {
	//args := os.Args[1:]

	partsModel := "MLKV3LL/A"
	location := "95051"
	url := fmt.Sprintf("https://www.apple.com/shop/retail/availabilitySearch?parts.0=%s&location=%s", partsModel, location)
	fmt.Println(url)
	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := resp{}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return 
	}
	stores := response.Body.Stores

	var availableStores []string
	for _, store := range stores {
		availability := store.PartsAvailability[partsModel]
		if availability.StoreSelectionEnabled && availability.PickupDisplay != "ineligible" {
			_ = append(availableStores, store.StoreName)
		}
	}
	fmt.Println(availableStores)
}
