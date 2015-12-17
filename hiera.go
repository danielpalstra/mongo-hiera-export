package main

import (
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	"log"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type HieraEntry struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Source string        `bson:"source,omitempty"`
	Key    string        `bson:"key,omitempty"`
	Value  interface{}   `bson:"value"`
}

var output = "hieradata"
var uri = os.Getenv("MONGOHQ_URL")
var db = "puppet"
var collection = "hiera"

func main() {

	// Connect to Mongo
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	// Collection hiera
	c := session.DB(db).C(collection)

	// Cound # of docs
	count, _ := c.Count()
	log.Printf("Found %s Hiera documents in MongoDB: %s\n", strconv.Itoa(count), uri)

	// Retrieve all available sources from mongo
	var sources []string
	c.Find(nil).Distinct("source", &sources)
	log.Printf("Found %s sources in MongoDB", strconv.Itoa(len(sources)))

	// Keep track of the mongo entries handled
	var entries int
	var duplicates int

	// Iterate over all the distinct sources
	for _, source := range sources {

		// Retrieve the k,v pairs per source
		var results []HieraEntry
		err = c.Find(bson.M{"source": source}).Sort("source").All(&results)
		if err != nil {
			panic(err)
		}

		// Create a new map containing the hiera k,v entries and fill it with the
		// mongo result set
		h := make(map[string]interface{})
		for _, result := range results {
			if val, ok := h[result.Key]; ok {
				log.Printf("Warn: Found already existing key in source: %s\nkey: %s\nfirst value : %s\nsecond value: %s \n", source, result.Key, val, result.Value)
				duplicates += 1
			}
			h[result.Key] = result.Value
		}

		// Create a new file containing the hieradata
		hieraData, err := yaml.Marshal(h)
		if err != nil {
			panic(err)
		}

		// Build yaml header and footer
		header := []byte("---\n")
		footer := []byte("\n\n")
		hieraFile := append(append(header, hieraData...), footer...)

		// Build yaml output
		err = ioutil.WriteFile(output+"/"+source+".yaml", hieraFile, 0644)
		if err != nil {
			log.Fatalf("Error wring %s to file: %s\n", source, err)
		}

		// Keep track of handled entries
		entries += len(results)

	}
	log.Printf("Found %s duplicates", duplicates)
	log.Printf("Finished converting %s mongo entries to yaml folder %s", strconv.Itoa(entries), output)
}
