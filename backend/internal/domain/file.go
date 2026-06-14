package domain

import "time"

type UploadFile struct {
	ID               int
	Filename         string
	OriginalFilename string
	FilePath         string
	FileSize         int64
	MimeType         string
	FileType         string
	Extension        string
	UploaderID       int
	Purpose          *string
	Md5              *string
	ReferenceCount   int
	CreatedAt        time.Time
}

type UploadFileDetail struct {
	UploadFile
	UploaderName     string
	UploaderNickname *string
}
