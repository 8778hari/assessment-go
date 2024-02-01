package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	request := make(chan string)
	response := make(chan string)

	go func() {
		// passing given request data to the worker via channel
		request <- string(reqBody)
	}()
	// run each process with goroutine
	go worker(request, response)

	fmt.Fprintf(w, string(<-response))
}

func worker(requests <-chan string, response chan<- string) {
	// converting given request data to desired response format
	givenReq := <-requests
	s := jsonToMap(givenReq)

	newAttr := make(map[string]interface{})
	newVal := make(map[string]interface{})
	result := make(map[string]interface{})

	for k := range s {
		if strings.Contains(k, "atr") {
			if strings.Contains(k, "uatr") {
				newVal[k] = s[k]
			} else {
				newAttr[k] = s[k]
			}
		} else {
			switch k {
			case "ev":
				result["event"] = s[k]
			case "et":
				result["event_type"] = s[k]
			case "id":
				result["app_id"] = s[k]
			case "uid":
				result["user_id"] = s[k]
			case "mid":
				result["message_id"] = s[k]
			case "p":
				result["page_url"] = s[k]
			case "t":
				result["page_title"] = s[k]
			case "l":
				result["browser_language"] = s[k]
			case "sc":
				result["screen_size"] = s[k]
			default:
				fmt.Println(k)
			}
		}
	}

	var x = map[string]map[string]string{}
	for a := range newAttr {
		if strings.Contains(a, "atrk") {
			str := strings.Replace(a, "atrk", "", 1)
			name := newAttr[a]
			value := newAttr["atrv"+str]
			typ := newAttr["atrt"+str]

			l1 := ""
			l2 := ""
			l3 := ""

			if str1, ok := name.(string); ok {
				l1 = str1
			}
			if str2, ok := value.(string); ok {
				l2 = str2
			}
			if str3, ok := typ.(string); ok {
				l3 = str3
			}

			x[l1] = map[string]string{}
			x[l1]["value"] = l2
			x[l1]["type"] = l3
		}
	}
	result["attributes"] = x

	var y = map[string]map[string]string{}
	for b := range newVal {
		if strings.Contains(b, "uatrk") {
			str := strings.Replace(b, "uatrk", "", 1)
			name := newVal[b]
			value := newVal["uatrv"+str]
			typ := newVal["uatrt"+str]

			l1 := ""
			l2 := ""
			l3 := ""

			if str1, ok := name.(string); ok {
				l1 = str1
			}
			if str2, ok := value.(string); ok {
				l2 = str2
			}
			if str3, ok := typ.(string); ok {
				l3 = str3
			}

			y[l1] = map[string]string{}
			y[l1]["value"] = l2
			y[l1]["type"] = l3
		}
	}
	result["traits"] = y

	j, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	response <- string(j)
}

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}
