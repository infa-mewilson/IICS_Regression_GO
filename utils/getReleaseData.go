package utils

import (
	"Golangcode/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io/ioutil"
	"log"
)

// GetReleaseData ...
func GetReleaseData(release_Iteration string, build string, usermetric string, userIndex string, release_Number string) map[string]int {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://asviicsperf03:9200",
		},
		// ...
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		panic(err)
	}
	//log.Println(es)

	//ping to check connectivity

	info, err := es.Ping()
	if err != nil {
		// Handle error
		log.Fatal("Unable to connect to elastic search")
		panic(err)

	}
	log.Println(info)
	//building query of this type
	//GET jmeter-aggregate-jtl/_search
	//{
	//	"query": {
	//	"bool": {
	//		"must": [
	//{"match": {"ReleaseNumber":202301
	//}},
	//{"match": {"BuildNumber":2 }}
	//]
	//}   // match_all
	//},
	//"size": 200
	//}

	var buffer bytes.Buffer
	query := map[string]interface{}{
		"size": 2000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"Release_Iteration": release_Iteration},
					}, {
						"match": map[string]interface{}{
							"BuildNumber": build},
					}, {
						"match": map[string]interface{}{
							"ReleaseNumber": release_Number,
						},
					},
				},
			},
		},
	}
	json.NewEncoder(&buffer).Encode(query)
	log.Println(&buffer)
	ctx := context.Background()
	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(userIndex),
		es.Search.WithBody(&buffer),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	//defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	var result config.Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON", err)
	}
	if result.Hits.Total.Value == 0 {
		log.Println("No Data for the selected release")
	}
	if result.Hits.Total.Value > 0 {
		log.Printf("Found a total of %d hits\n", result.Hits.Total.Value)
	}
	responsedata := make(map[string]int)
	var key2 int
	for _, hit := range result.Hits.Hits {

		if usermetric == "95th" {
			key2 = hit.Source.Nine5Th
		}
		if usermetric == "90th" {
			key2 = hit.Source.Nine0Th
		}
		if usermetric == "99th" {
			key2 = hit.Source.Nine9Th
		}
		if usermetric == "Average" {
			key2 = hit.Source.Average
		}
		responsedata[hit.Source.Label] = key2
	}
	return responsedata
}
