package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Item struct {
	Title string  `json:"title"`
	Body string   `json:"body"`
	Id int    `json:"id"`
}

func main() {
	startTime := time.Now()
	log.Println("Starting")

	strings := []string{"hello", "world","This, is, a, test, row",
		"This", "is", "another", "test", "row", "my", "dude",
		"100, 200, 300, 400, 500"}

	//setup a number of waitgroups for concurrency channels
	var wg sync.WaitGroup
	wg.Add(len(strings))
	// unbuffered channel
	results := make(chan string)

 	for i, s := range strings {
		go func(row int, stuff string) {
			log.Printf("Sending Data Row(%d) Value(%s)\n",row, stuff)
			log.Printf("WG state %v\n", wg)
			exData := Item{
				Title: stuff,
				Body: stuff,
				Id: row,
			}

			var jsonData []byte
			jsonData, err := json.Marshal(exData)
			if err != nil {
				log.Fatalf("Error marshalling the Json: %v",err)
			}
			log.Printf("Sending Json %s\n",string(jsonData))
			retString := postData("https://jsonplaceholder.typicode.com/posts",jsonData)
			results <- retString

			wg.Done()
		}(i,s)
	}

	// wait for all the channels to close out before finishing,  this is blocking
	go func() {
		wg.Wait()
		close(results)
	}()
	display(results)

	log.Printf("WG XXX state %v\n", wg)
	log.Printf("Took %v",time.Now().Sub(startTime).Round(time.Microsecond).String())
}

func display(results chan string) {
	for r :=  range results {
		log.Println("----------------")
		log.Printf("Returned: %v",r)
		log.Println("----------------")
	}
}


func postData(url string, payload []byte) string{
	//fmt.Printf("Posting to %s \n", url)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")

	//client := &http.Client{CheckRedirect: redirectPolicyFunc}
	client := &http.Client{}
	//request.Header.Add("Authorization", "Basic "+basicAuth(s.SnowAPIUser, s.SnowAPIPass))
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error: The HTTP POST request failed with error %s\n", err)
		return ""
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		return("OK:" + string(data))
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("kubernetes_user", "kubernetes_user"))
	return nil
}