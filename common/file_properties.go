package common

type UploadFileProperties struct {
	File            ReaderAtSeekerCloser
	Size            int64
	Category        string
	SubCategory     string
	Filename        string
	MimeType        string
	OtherProperties *UnknownMap
}
