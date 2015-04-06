package main // import "github.com/cloudspace/Go_Location_State"

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(errorStringAsJSON(fmt.Sprintf("Must have 3 argument: your are passing %v arguments", len(os.Args)-1)))
		return
	}

	strLat := os.Args[1]
	strLng := os.Args[2]

	lat, err := strconv.ParseFloat(strLat, 64)

	if err != nil {
		fmt.Println(getJSONError(err))
	}

	lng, err := strconv.ParseFloat(strLng, 64)

	if err != nil {
		fmt.Println(getJSONError(err))
	}

	address, err := reverseGeocodeState(lat, lng)

	if err != nil {
		fmt.Println(getJSONError(err))
	}

	result := make(map[string]interface{}, 0)
	result["result"] = address
	result["error"] = ""

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

	return "{\"result\": \"\"\n\"error\": \"" + errorString + "\"}"
}

func reverseGeocodeState(lat float64, lng float64) (*string, error) {
	geoPoint := &point{
		Lat: lat,
		Lng: lng}

	result, err := reverseGeocode(*geoPoint)
	if err != nil {
		return nil, err
	}

	if result.StatusCode != 200 {
		return nil, fmt.Errorf("Got non 200 status code %d for request %v", result.StatusCode, result.QueryString)
	}

	if result.Count == 0 {
		return nil, fmt.Errorf("Got 0 results for request %v", result.QueryString)
	}

	state, err := getAddressNameForType(result.googleResponse.Results[0], cAddressPartAdministrativeAreaLevel1)

	if err != nil {
		return nil, err
	}

	return state, nil
}

const (
	cGoogle                              = "http://maps.googleapis.com/maps/api/geocode/json"
	cAddressPartStreetNumber             = "street_number"
	cAddressPartRoute                    = "route"
	cAddressPartLocality                 = "locality"
	cAddressPartPolitical                = "political"
	cAddressPartAdministrativeAreaLevel2 = "administrative_area_level_2"
	cAddressPartAdministrativeAreaLevel1 = "administrative_area_level_1"
	cAddressPartCountry                  = "country"
	cAddressPartPostalCode               = "postal_code"
	cAddressPartPostalCodeSuffix         = "postal_code_suffix"
)

//http://maps.googleapis.com/maps/api/geocode/json?sensor=false&latlng=28.1,-81.6
type point struct {
	Lat, Lng float64
}

type bounds struct {
	NorthEast, SouthWest point
}

type googleResponse struct {
	Results []*googleResult
}

type googleResult struct {
	Address      string               `json:"formatted_address"`
	AddressParts []*googleAddressPart `json:"address_components"`
	Geometry     *geometry
	Types        []string
}

type googleAddressPart struct {
	Name      string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types     []string
}

type geometry struct {
	Bounds   bounds
	Location point
	Type     string
	Viewport bounds
}

type response struct {
	StatusCode  int
	QueryString string
	Found       string
	Count       int
	*googleResponse
}

func getAddressNameForType(gresult *googleResult, addressType string) (*string, error) {
	addressPart, err := getAddressPartForType(gresult, addressType)

	if err != nil {
		return nil, err
	}

	return &addressPart.Name, nil
}

func getAddressPartForType(gresult *googleResult, addressType string) (*googleAddressPart, error) {
	for _, eachAddressPart := range gresult.AddressParts {
		for _, eachAddressType := range eachAddressPart.Types {
			if eachAddressType == addressType {
				return eachAddressPart, nil
			}
		}
	}

	return nil, errors.New("Could not find state")
}

func reverseGeocode(latlng point) (*response, error) {

	client := &http.Client{}

	url := fmt.Sprintf("%s?sensor=false&latlng=%f,%f", cGoogle, latlng.Lat, latlng.Lng)

	req, err := http.NewRequest("GET", url, nil)

	clientResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer clientResp.Body.Close()

	resp := new(response)
	resp.QueryString = url

	if clientResp.StatusCode == 200 { // OK
		err = json.NewDecoder(clientResp.Body).Decode(resp)
		// reverse geocoding
		resp.Count = len(resp.googleResponse.Results)
		if resp.Count >= 1 {
			resp.Found = resp.googleResponse.Results[0].Address
		}
		if err != nil {
			return nil, err
		}
	}

	resp.StatusCode = clientResp.StatusCode
	return resp, nil
}
