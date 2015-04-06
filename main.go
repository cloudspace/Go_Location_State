package main // import "github.com/cloudspace/Go_Location_State"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(errorStringAsJSON(fmt.Sprintf("Must have 1 argument: your are passing %v arguments", len(os.Args)-1)))
		return
	}

	urlToExtend := os.Args[1]

	resp, err := http.Get(urlToExtend)
	if err != nil {
		fmt.Println(getJSONError(err))
	}

	finalURL := resp.Request.URL.String()

	result := make(map[string]interface{}, 0)
	result["result"] = finalURL

	fmt.Println(asJSON(result))
}

func asJSON(anything interface{}) string {

	jsonData, err := json.Marshal(anything)
	if err != nil {
		return getJSONError(err)
	}
	return string(jsonData)
}

func getJSONError(myError error) string {

	errorJSON := make(map[string]interface{})
	errorJSON["error"] = myError.Error()
	jsonData, err := json.Marshal(errorJSON)
	if err != nil {
		return errorStringAsJSON("There was an error generatoring the error.. goodluck")
	}
	return string(jsonData)
}

func errorStringAsJSON(errorString string) string {

	return "{\"error\": \"" + errorString + "\"}"
}
