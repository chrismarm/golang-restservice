package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type DiveLocation struct {
	Id    string `json:"_id,omitempty"`
	Name  string `json:"name"`
	Lat   string `json:"lat,omitempty"`
	Lon   string `json:"lon,omitempty"`
	Depth string `json:"depth,omitempty"`
}

// Simulates in-memory cache from database, that is a json file
var dives []DiveLocation

// Map for indexation by name and quick retrieval
var divesIndex map[string]DiveLocation

const filename = "dive_locations.json"

func main() {
	// Database (json file) initial retrieval
	loadDatabase()

	// API methods creation
	router := mux.NewRouter()
	router.HandleFunc("/dives", GetDives).Methods("GET")
	router.HandleFunc("/dives/{name}", GetDive).Methods("GET")
	router.HandleFunc("/dives/{name}", CreateOrUpdateDive).Methods("POST")
	router.HandleFunc("/dives/{name}", DeleteDive).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getData() []byte {
	var jsonBlob = []byte(`
        [{"_id":"1","name":"San Andres","lat":"36.0144638","lon":"-5.6090361","depth":"34"},{"_id":"32eaa763-3eeb-4616-bbd8-333c3756f7f5","name":"Calderas","lat":"222","lon":"333","depth":"20"}]
    `)
	return jsonBlob
}

func getDataFromFile() []byte {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonBlob, _ := ioutil.ReadAll(jsonFile)
	return jsonBlob
}

func loadDatabase() {
	jsonBlob := getDataFromFile()

	json.Unmarshal(jsonBlob, &dives)
	divesIndex = make(map[string]DiveLocation)

	// Print current DB content and populates the index map
	for i := 0; i < len(dives); i++ {
		currentDive := dives[i]
		diveName := currentDive.Name
		divesIndex[diveName] = currentDive
		fmt.Println("---Dive #" + strconv.Itoa(i+1))
		fmt.Println("Id: " + currentDive.Id)
		fmt.Println("Name: " + diveName)
		fmt.Println("Lat: " + currentDive.Lat)
		fmt.Println("Lon: " + currentDive.Lon)
		fmt.Println("Depth: " + currentDive.Depth)
	}
}

func updateDatabase() {
	// Writes to file
	divesJson, _ := json.Marshal(dives)
	ioutil.WriteFile(filename, divesJson, 0644)
}

func getDiveIndexInCache(name string) int {
	// Looks for a dive secuentially with the given name
	for i := 0; i < len(dives); i++ {
		currentDive := dives[i]
		if currentDive.Name == name {
			return i
		}
	}
	return -1
}

func generateId() string {
	var result string
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
		result = "defaultId"
	} else {
		// Remove last '\n' char
		result = string(uuid[:len(uuid)-1])
	}
	return result
}

// ---- API methods

func GetDives(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dives)
}

func GetDive(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	diveName := params["name"]
	// Looks for in the index map
	existingDive, present := divesIndex[diveName]
	if !present {
		// Empty response as there is no dive with that name
		json.NewEncoder(w).Encode(&DiveLocation{})
	} else {
		json.NewEncoder(w).Encode(existingDive)
	}
}

func CreateOrUpdateDive(w http.ResponseWriter, r *http.Request) {
	// Payload parsing in the form of a dive location
	var newDive DiveLocation
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&newDive); err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(&DiveLocation{})
		return
	}

	params := mux.Vars(r)
	diveName := params["name"]
	newDive.Name = diveName
	_, present := divesIndex[diveName]
	// Index map update
	divesIndex[diveName] = newDive
	if !present {
		// New dive with UUID generation
		newDive.Id = generateId()
		dives = append(dives, newDive)
		json.NewEncoder(w).Encode(dives)
	} else {
		// Update dive
		index := getDiveIndexInCache(diveName)
		dives[index] = newDive
		json.NewEncoder(w).Encode(dives)
	}
	updateDatabase()
}

func DeleteDive(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	diveName := params["name"]
	_, present := divesIndex[diveName]
	if !present {
		// Empty response as there is no dive with that name
		json.NewEncoder(w).Encode(&DiveLocation{})
	} else {
		// Update cache
		delete(divesIndex, diveName)
		index := getDiveIndexInCache(diveName)
		dives = append(dives[:index], dives[index+1:]...)
		json.NewEncoder(w).Encode(dives)
	}
	updateDatabase()
}
