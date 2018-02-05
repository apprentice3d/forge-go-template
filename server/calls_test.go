package server

import (
	oauth2 "github.com/apprentice3d/forge-api-go-client/oauth"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateTransientBucket(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	oauth := oauth2.NewTwoLeggedClient(clientID, clientSecret)

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

	oauth := oauth2.NewTwoLeggedClient(clientID, clientSecret)

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

	objectID, err := UploadDataIntoBucket("tester.txt", data, "bucket5577006791947779410", bearer.AccessToken)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(objectID) == 0 {
		t.Fatal("ObjectId is empty")
	}

}

func TestTranslateSourceToSVF(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	oauth := oauth2.NewTwoLeggedClient(clientID, clientSecret)

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

