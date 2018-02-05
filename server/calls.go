package server

import (
	"strconv"
	"math/rand"
	"net/http"
	"encoding/json"
	"bytes"
	"errors"
	"io/ioutil"
	"encoding/base64"
	"log"
	"time"
)

// UploadAndConvert takes cares of file uploading and translation, given the file and access token
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

	//Checking the translation progress but not more than 360 times with interval of 10 sec => approx 1 hour
	counter := 359
	for {
		progress, err := CheckTranslationProgress(urn, token)
		if err != nil || progress == "complete" || counter < 0 {
			return urn, err
		}
		log.Printf("Translation for URN=%s not yet complete. [Will retry in 10 sec]", urn)
		time.Sleep(10 * time.Second)
		counter -= 1
	}

	return
}

// TranslateSourceToSVF takes care of base64 conversion of objectID and returns the URN
// for which translation was started
func TranslateSourceToSVF(objectId string, token string) (urn string, err error) {

	base64urn := base64.RawStdEncoding.EncodeToString([]byte(objectId))

	params := TranslationParams{}
	params.Input.URN = base64urn
	format := Format{
		Type:  "svf",
		Views: []string{"2d", "3d"},
	}
	params.Output.Formats = []Format{format}

	byteParams, err := json.Marshal(params)
	if err != nil {
		log.Println("Could not marshal the translation parameters")
		return
	}

	req, err := http.NewRequest("POST",
		"https://developer.api.autodesk.com/modelderivative/v2/designdata/job",
		bytes.NewBuffer(byteParams))

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

// UploadDataIntoBucket is responsible for uploading the received file into given bucket
func UploadDataIntoBucket(filename string, data []byte, bucketKey string, token string) (objectId string, err error) {

	url := "https://developer.api.autodesk.com/oss/v2/buckets/" + bucketKey + "/objects/" + filename

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))

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

		return "", errors.New("Could not upload file: " + strconv.Itoa(res.StatusCode) + " ==> " + string(data))
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

// CreateTransientBucket is responsible for creation of a transient bucket given the bucket name
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


// CheckTranslationProgress will check the status of the work and will return progress either "complete"
// or as percent value
func CheckTranslationProgress(urn string, token string) (progress string, err error) {

	url := "https://developer.api.autodesk.com/modelderivative/v2/designdata/" +
		urn + "/manifest"

	req, err := http.NewRequest("GET",
		url,
		nil)

	if err != nil {

		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)

	defer res.Body.Close()
	if err != nil {

		return
	}

	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)

		return "", errors.New(strconv.Itoa(res.StatusCode) + " ==> " + string(data))
	}

	decoder := json.NewDecoder(res.Body)
	var response TranslationStatusResponse
	err = decoder.Decode(&response)

	if err != nil {
		return "", errors.New("Could not unmarshal translation progress response: " + err.Error())
	}

	progress = response.Progress
	log.Printf("Checked translation status for URN=%s ==> %s\n",
		urn, progress)

	if response.Progress == "failed" {
		err = errors.New("Translation FAILED for URN=" + urn)
	}
	return

}
