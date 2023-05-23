package common

type APIModeType string

const (
	APIMODE_DEBUG APIModeType = "DEBUG"
	APIMODE_QA    APIModeType = "QA"
	APIMODE_PROD  APIModeType = "PROD"
	APIMODE_TEST  APIModeType = "TEST"
)
