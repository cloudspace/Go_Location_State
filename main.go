package main // import "github.com/cloudspace/Go_Location_State"

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(errorStringAsJSON(fmt.Sprintf("Must have 2 argument: your are passing %v arguments", len(os.Args)-1)))
		return
	}

	strLat := os.Args[1]
	strLng := os.Args[2]

	lat, err := strconv.ParseFloat(strLat, 64)

	if err != nil {
		fmt.Println(getJSONError(err))
		return
	}

	lng, err := strconv.ParseFloat(strLng, 64)

	if err != nil {
		fmt.Println(getJSONError(err))
		return
	}

	cmd := exec.Command("sh", "-c", "service postgresql start")
	err = cmd.Run()

	if err != nil {
		fmt.Println(getJSONError(err))
		return
	}

	connectionURI := "host=127.0.0.1 port=5432 user=docker password=docker dbname=geolocation"
	query := fmt.Sprintf("SELECT name FROM ne_110m_admin_1_states_provinces WHERE ST_Contains(geom, ST_GeometryFromText('POINT(%f %f)', 4326))", lng, lat)

	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		fmt.Println(getJSONError(err))
		return
	}
	defer db.Close()

	result, err := getJSONResultOfQuery(query, db)
	if err != nil {
		fmt.Println(getJSONError(err))
		return
	}
	fmt.Println(result)
}

func getJSONResultOfQuery(sqlString string, db *sql.DB) (string, error) {
	rows, err := db.Query(sqlString)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	var tableData []map[string]interface{}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	result := make(map[string]interface{}, 0)
	if len(tableData) == 0 {
		result["state"] = ""
		result["error"] = ""
	} else {
		result["state"] = tableData[0]["name"]
		result["error"] = ""
	}
	return asJSON(result), nil
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

	return "{\"state\": \"\"\n\"error\": \"" + errorString + "\"}"
}
