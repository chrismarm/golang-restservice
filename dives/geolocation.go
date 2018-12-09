package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Location struct {
	GlobalCode   string `json:"plus_code,omitempty"`
	LocationName string `json:"location,omitempty"`
	Locality     string `json:"locality,omitempty"`
	Country      string `json:"country,omitempty"`
}

func ReverseGeocode(latitude string, longitude string, key string) Location {
	resp, _ := getResponse(latitude, longitude, key)

	parseRes := gjson.Parse(resp)

	global_code := parseRes.Get("plus_code.global_code").String()
	location := Location{
		GlobalCode: global_code,
	}

	results := parseRes.Get("results").Array()
	if len(results) == 0 {
		compound_code := parseRes.Get("plus_code.compound_code").String()
		location.LocationName = compound_code
	} else if len(results) == 1 {
		result := results[0]
		formatted_address := result.Get("formatted_address").String()
		location.LocationName = formatted_address
		var locality, country string
		result.Get("address_components").ForEach(func(key, value gjson.Result) bool {
			typeR := value.Get("types.0").String()
			switch typeR {
			case "locality":
				locality = value.Get("long_name").String()
				fmt.Println(locality)
			case "country":
				country = value.Get("long_name").String() + " (" + value.Get("short_name").String() + ")"
				fmt.Println(country)
			}
			return true
		})
		location.Locality = locality
		location.Country = country
	}

	return location
}

func getResponse(latitude string, longitude string, apiKey string) (string, error) {
	// Google Maps Geocode API call
	u := url.URL{
		Host:   "maps.googleapis.com",
		Path:   "maps/api/geocode/json",
		Scheme: "https",
	}
	q := u.Query()
	q.Set("key", apiKey)
	q.Add("latlng", latitude+","+longitude)
	q.Add("result_type", "locality")
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())
	if err != nil {
		log.Fatal("The HTTP request failed with error %s\n", err)
		return "", err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		result := string(data)
		return result, nil
	}
}
