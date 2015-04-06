# Go Location State
A Go Microservice to return a state for latitude and longitude.

#### Usage
Input must have 2 arguments:

1.  Latitude
2.  Longitude

#### Example Input - Output
-
Input:
```
go run main.go 28.1 81.6
```
Output (Success):
```
{"result":"Florida"}
```
-
Output (Failure):
```
{
  "error":"<error message>"
}
```

#### How to build docker image
Requirements:

1. Golang environment set up
2. Git
3. Boot2docker running

```
go get github.com/cloudspace/Go_Location_State
cd <Go_Location_State Directory>
docker run --rm -v $(pwd):/src centurylink/golang-builder
docker build -t <username>/go_location_state:0.1 ./

```

In order for the golang-builder to work, you need to have the github url on the top line of main.go. It should look like this:
```
package main // import "github.com/cloudspace/Go_Location_State"
```
You also must *push your code* to github *before building* the docker image.
