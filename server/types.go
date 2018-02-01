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
