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
)

func UploadAndConvert(filename string, data []byte, token string) (urn string, err error) {

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

	//params := TranslationParams{
	//	Input: struct{ URN string }{URN: base64urn},
	//	Output: struct {
	//		Formats []struct {
	//			Type  string   `json:"type"`;
	//			Views []string `json:"views"`
	//		}
	//	}{Formats: }
	//
	//}


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

	//var response UploadFileToBucketResponse
	//result, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return
	//}
	//log.Printf("Received response: %s\n", string(result))
	//err = json.Unmarshal(result,response)


	if err != nil {
		return "", errors.New("Could not unmarshal upload response: " + err.Error())
	}

	objectId = response.ObjectID
	log.Printf("File '%s' was successfully uploaded into bucket '%s' and now has ID: %s\n",
		filename, bucketKey, objectId)
	return
}

func CreateTransientBucket(bucketName string, token string) (bucketKey string, err error) {
	params := BucketParamaters{
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

type BucketParamaters struct {
	BucketKey string `json:"bucketKey"`
	PolicyKey string `json:"policyKey"`
}

type TranslationParams struct {
	Input struct {
		URN string `json:"urn"`
	} `json:"input"`
	Output struct {
		Formats []struct {
			Type  string   `json:"type"`
			Views []string `json:"views"`
		} `json:"formats"`
	} `json:"output"`
}

type TranslationResponse struct {
	Result string `json:"result"`
	URN    string `json:"urn"`
}

type UploadFileToBucketResponse struct {
	BucketKey   string `json:"bucketKey"`
	ObjectID    string `json:"objectId"`
	ObjectKey   string `json:"objectKey"`
	SHA1        string `json:"sha1"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	Location    string `json:"location"`
}

/*
{
"bucketKey": "some_temp_bucket",
"objectId": "urn:adsk.objects:os.object:some_temp_bucket/giro_watch.f3d",
"objectKey": "giro_watch.f3d",
"sha1": "d0e31a2f8be8a47a5b1a514064a4a25ddf29f9a4",
"size": 1488828,
"contentType": "application/x-www-form-urlencoded",
"location": "https://developer.api.autodesk.com/oss/v2/buckets/some_temp_bucket/objects/giro_watch.f3d"
}

*/

type CreateBucketResponse struct {
	BucketKey   string `json:"bucketKey"`
	BucketOwner string `json:"bucketOwner"`
	CreatedDate int64  `json:"createdDate"`
	Permissions []struct {
		AuthID string `json:"authId"`
		Access string `json:"access"`
	} `json:"permissions"`
	PolicyKey string `json:"policyKey"`
}
