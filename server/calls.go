package server

import (
	"strconv"
	"math/rand"
	"net/http"
	"encoding/json"
	"bytes"
	"errors"
	"mime/multipart"
	"io/ioutil"
	"encoding/base64"
	"log"
	"time"
)

func UploadAndConvert(filename string, data []byte, token string) (urn string, err error) {

	rand.Seed(time.Now().UTC().UnixNano())
	randomBucketName := "bucket" + strconv.Itoa(rand.Int())
	log.Printf("Creating a transient bucket with name: %s\n", randomBucketName)
	bucketKey, err := CreateTransientBucket(randomBucketName, token)
	if err != nil {
		log.Printf("FAIL: could not create bucket: %s\n", err.Error())
		return
	}

	log.Printf("Uploading file '%s' into bucket '%s'\n", filename, bucketKey)
	objectId, err := UploadDataIntoBucket(filename, data, bucketKey, token)
	if err != nil {
		return
	}
	log.Printf("Request object '%s' to be translated into SVF\n", objectId)
	urn, err = TranslateSourceToSVF(objectId, token)

	return
}


func TranslateSourceToSVF(objectId string, token string) (urn string, err error) {

	base64urn := base64.RawStdEncoding.EncodeToString([]byte(objectId))

	//TODO: replace hardcodings with TranslationParams struct
	var params = []byte(`{
  "input": {
    "urn": "`)
	params = append(params, []byte(base64urn)...)

	params = append(params, []byte(`"
  },
  "output": {
    "formats": [{
      "type": "svf",
      "views": [
        "2d",
        "3d"
        ]
    }]
  }
}`)...)

	req, err := http.NewRequest("POST",
		"https://developer.api.autodesk.com/modelderivative/v2/designdata/job",
		bytes.NewBuffer(params))

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)

		return "", errors.New("Fail, received status code " + strconv.Itoa(res.StatusCode) + " ==> " + string(data))
	}

	decoder := json.NewDecoder(res.Body)
	var response TranslationResponse

	err = decoder.Decode(&response)

	if err != nil {
		return
	}

	urn = response.URN
	log.Printf("Object '%s' was successfully sent for translation\n",
		objectId)
	return

}

func UploadDataIntoBucket(filename string, data []byte, bucketKey string, token string) (objectId string, err error) {

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		return
	}
	fileWriter.Write(data)

	url := "https://developer.api.autodesk.com/oss/v2/buckets/"+bucketKey+"/objects/" + filename

	req, err := http.NewRequest("PUT", url, bodyBuffer)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)

		return "", errors.New("Fail, received status code " + strconv.Itoa(res.StatusCode) + " ==> " + string(data))
	}

	decoder := json.NewDecoder(res.Body)
	var response UploadFileToBucketResponse
	err = decoder.Decode(&response)

	if err != nil {
		return "", errors.New("Could not unmarshal upload response: " + err.Error())
	}

	objectId = response.ObjectID
	log.Printf("File '%s' was successfully uploaded into bucket '%s' and now has ID: %s\n",
		filename, bucketKey, objectId)
	return
}

func CreateTransientBucket(bucketName string, token string) (bucketKey string, err error) {
	params := BucketParams{
		bucketName,
		"transient",
	}
	payload, err := json.Marshal(params)

	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		"https://developer.api.autodesk.com/oss/v2/buckets",
		bytes.NewBuffer(payload))

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return
	}

	if res.StatusCode == http.StatusConflict {
		// Bucket already exists
		log.Printf("Bucket '%s' already exist, writting into it\n", bucketName)
		return bucketName, nil
	}

	decoder := json.NewDecoder(res.Body)
	var response CreateBucketResponse

	err = decoder.Decode(&response)

	if err != nil {
		return
	}

	bucketKey = response.BucketKey
	log.Printf("Bucket '%s' successfully created\n", bucketKey)
	return
}

