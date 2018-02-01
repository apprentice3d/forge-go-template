package server

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

var currentURN string



func (service ForgeServices) getToken(writer http.ResponseWriter, request *http.Request) {

	bearer, err := service.oauth.Authenticate("viewables:read")
	if err != nil {
		writer.WriteHeader(http.StatusNotAcceptable)
		writer.Write([]byte(err.Error()))
		return
	}
	log.Printf("Received a token request: returning a token that will expire in %d\n", bearer.ExpiresIn)
	encoder := json.NewEncoder(writer)
	encoder.Encode(bearer)
}


func (service ForgeServices) getURN(writer http.ResponseWriter, request *http.Request) {

	//TODO: change this to dynamic URN
	//urn := "urn:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6c29tZV90ZW1wX2J1Y2tldC9naXJvX3dhdGNoLmYzZA"

	log.Printf("Received an URN request. Returning %s\n", currentURN)
	writer.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(writer)
	encoder.Encode(struct{
		URN string `json:"urn"`
	}{currentURN})
}

func (service *ForgeServices) uploadFiles(writer http.ResponseWriter, request *http.Request) {

	data, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError);
		writer.Write([]byte("Could not read the body"))
		return
	}

	headerData := request.Header["Filename"]
	if len(headerData) == 0 {
		writer.WriteHeader(http.StatusInternalServerError);
		writer.Write([]byte("Could not retrieve filename"))
		return
	}
	filename := headerData[0]

	log.Printf("Received request to translate file: %s\n", filename)
	bearer, err := service.oauth.Authenticate("data:read data:write bucket:create bucket:read viewables:read")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError);
		writer.Write([]byte("Could not acquire a token"))
		return
	}


	urn, err := UploadAndConvert(filename, data, bearer.AccessToken)

	log.Printf("Translation was successful. Got URN: %s\n", urn)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError);
		writer.Write([]byte(err.Error()))
		return
	}

	currentURN = urn
	log.Printf("Setting current URN to: %s\n", currentURN)

	writer.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(writer)
	encoder.Encode(struct{
		URN string `json:"urn"`
	}{urn})
}
