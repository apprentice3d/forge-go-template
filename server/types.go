package server

import (
	"github.com/apprentice3d/forge-api-go-client/oauth"
)

// ForgeServices holds the necessary references to Forge services
type ForgeServices struct {
	oauth oauth.TwoLeggedAuth
}

// BucketParams struct reflects the Bucket Creation parameters
type BucketParams struct {
	BucketKey string `json:"bucketKey"`
	PolicyKey string `json:"policyKey"`
}

// TranslationResponse struct reflects the necessary data returned upon starting file translation
type TranslationResponse struct {
	Result string `json:"result"`
	URN    string `json:"urn"`
}

// UploadFileToBucketResponse struct reflects the data returned upon file uploading into a bucket
type UploadFileToBucketResponse struct {
	BucketKey   string `json:"bucketKey"`
	ObjectID    string `json:"objectId"`
	ObjectKey   string `json:"objectKey"`
	SHA1        string `json:"sha1"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	Location    string `json:"location"`
}

// CreateBucketResponse reflects the the data returned upon creating a bucket
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

// TranslationStatusResponse reflect the data returned when checking for translation status
type TranslationStatusResponse struct {
	Type     string `json:"type"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
	Region   string `json:"region"`
	URN      string `json:"urn"`
	Version  string `json:"version"`
}

// TranslationParams reflects the data and the structure necessary to start a translation job
type TranslationParams struct {
	Input struct {
		URN string `json:"urn"`
	} `json:"input"`
	Output struct {
		Formats []Format `json:"formats"`
	} `json:"output"`
}

// Format struct is part of TranslationParams struct and reflects the type to be translated to and the views
type Format struct {
	Type  string   `json:"type"`
	Views []string `json:"views"`
}

/*

`{
  "input": {
    "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6c29tZV90ZW1wX2J1Y2tldC9jZWFzaWsuZjNk"
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
}`
*/

/*
"type": "manifest",
    "hasThumbnail": "true",
    "status": "success",
    "progress": "complete",
    "region": "US",
    "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6c29tZV90ZW1wX2J1Y2tldC9naXJvX3dhdGNoLmYzZA",
    "version": "1.0",
*/
