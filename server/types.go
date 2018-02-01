package server

import (
	"github.com/apprentice3d/forge-api-go-client/oauth"
)

type ForgeServices struct {
	oauth oauth.AuthApi
}

type BucketParams struct {
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


type TranslationStatusResponse struct {
	Type string		`json:"type"`
	Status string	`json:"status"`
	Progress string	`json:"progress"`
	Region string	`json:"region"`
	URN string		`json:"urn"`
	Version string	`json:"version"`
}



/*
"type": "manifest",
    "hasThumbnail": "true",
    "status": "success",
    "progress": "complete",
    "region": "US",
    "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6c29tZV90ZW1wX2J1Y2tldC9naXJvX3dhdGNoLmYzZA",
    "version": "1.0",
 */