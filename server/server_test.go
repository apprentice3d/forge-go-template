package server

import (
	"testing"
	"os"
	oauth2 "github.com/apprentice3d/forge-api-go-client/oauth"
	"io/ioutil"
)

func TestCreateTransientBucket(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	oauth := oauth2.NewTwoLeggedClient(clientID,clientSecret)

	bearer, err := oauth.Authenticate("data:read data:write bucket:create bucket:read viewables:read")

	if err != nil {
		t.Fatal(err.Error())
	}

	bucketKey, err := CreateTransientBucket("something2", bearer.AccessToken)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(bucketKey) == 0 {
		t.Fatal("BucketKey is empty")
	}
}


func TestUploadDataIntoBucket(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	oauth := oauth2.NewTwoLeggedClient(clientID,clientSecret)

	bearer, err := oauth.Authenticate("data:write")

	if err != nil {
		t.Fatal(err.Error())
	}

	file, err := os.Open("../tmpfile.txt")
	if err != nil {
		t.Fatal(err.Error())
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		t.Fatal(err.Error())
	}

	//bearer.AccessToken = "eyJhbGciOiJIUzI1NiIsImtpZCI6Imp3dF9zeW1tZXRyaWNfa2V5In0.eyJjbGllbnRfaWQiOiJlMTMwcXRBSnNKZDMwQ1kzV3dBY0IzVFBXT05JUDFZQSIsImV4cCI6MTUxNzQ2MDQ1MSwic2NvcGUiOlsiZGF0YTpyZWFkIiwiZGF0YTp3cml0ZSIsImJ1Y2tldDpjcmVhdGUiLCJidWNrZXQ6cmVhZCIsInZpZXdhYmxlczpyZWFkIl0sImF1ZCI6Imh0dHBzOi8vYXV0b2Rlc2suY29tL2F1ZC9qd3RleHA2MCIsImp0aSI6IjZrOEVxaWhma0FHMFExMlRGYXlETHdTdkRiZ3VZa201N01QYkJZZ1A3eDFza3duNlJGRmhSZXFzaHpENFVxU0MifQ._CRyc4SdM4JPv-br9vVGkIRSXlJ1JqwJU3YAeM9gvZk"


	objectId, err := UploadDataIntoBucket("tester.txt", data, "bucket5577006791947779410", bearer.AccessToken)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(objectId) == 0 {
		t.Fatal("ObjectId is empty")
	}


}


func TestTranslateSourceToSVF(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	oauth := oauth2.NewTwoLeggedClient(clientID,clientSecret)

	bearer, err := oauth.Authenticate("data:read data:write bucket:create bucket:read viewables:read")

	if err != nil {
		t.Fatal(err.Error())
	}

	objectID := "urn:adsk.objects:os.object:some_temp_bucket/giro_watch.f3d"
	urn, err := TranslateSourceToSVF(objectID, bearer.AccessToken)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(urn) == 0 {
		t.Fatal("URN is empty")
	}

}

