package service

import (
	"context"

	"oa-nsdiy/backend/internal/domain"
)

// FileRepository defines the interface for file data access required by FileService.
type FileRepository interface {
	List(ctx context.Context, keyword, fileType string, uploaderID *int, page, pageSize int) ([]*domain.UploadFileDetail, int64, error)
	GetDetailByID(ctx context.Context, id int) (*domain.UploadFileDetail, error)
	GetByID(ctx context.Context, id int) (*domain.UploadFile, error)
	Create(ctx context.Context, file *domain.UploadFile) error
	Update(ctx context.Context, file *domain.UploadFile) error
	Delete(ctx context.Context, id int) error
}

type FileService struct {
	repo FileRepository
}

func NewFileService(repo FileRepository) *FileService {
	return &FileService{repo: repo}
}

type FileCreateInput struct {
	Filename         string
	OriginalFilename string
	FilePath         string
	FileSize         int64
	MimeType         string
	FileType         string
	Extension        string
	Purpose          string
	Md5              string
}

type FileUpdateInput struct {
	Filename         string
	OriginalFilename string
	FilePath         string
	FileSize         int64
	MimeType         string
	FileType         string
	Extension        string
	Purpose          string
	Md5              string
}

type FileListResult struct {
	Items []*domain.UploadFileDetail
	Total int64
}

func (s *FileService) ListFiles(ctx context.Context, keyword, fileType string, uploaderID *int, page, pageSize int) (*FileListResult, error) {
	items, total, err := s.repo.List(ctx, keyword, fileType, uploaderID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &FileListResult{Items: items, Total: total}, nil
}

func (s *FileService) GetFile(ctx context.Context, id int) (*domain.UploadFileDetail, error) {
	return s.repo.GetDetailByID(ctx, id)
}

func (s *FileService) CreateFile(ctx context.Context, input FileCreateInput, uploaderID int) (*domain.UploadFile, error) {
	file := &domain.UploadFile{
		Filename:         input.Filename,
		OriginalFilename: input.OriginalFilename,
		FilePath:         input.FilePath,
		FileSize:         input.FileSize,
		MimeType:         input.MimeType,
		FileType:         input.FileType,
		Extension:        input.Extension,
		UploaderID:       uploaderID,
		Purpose:          &input.Purpose,
		Md5:              &input.Md5,
		ReferenceCount:   1,
	}

	if err := s.repo.Create(ctx, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) UpdateFile(ctx context.Context, id int, input FileUpdateInput) (*domain.UploadFile, error) {
	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "file_not_found", "File not found")
	}

	if input.Filename != "" {
		file.Filename = input.Filename
	}
	if input.OriginalFilename != "" {
		file.OriginalFilename = input.OriginalFilename
	}
	if input.FilePath != "" {
		file.FilePath = input.FilePath
	}
	if input.FileSize > 0 {
		file.FileSize = input.FileSize
	}
	if input.MimeType != "" {
		file.MimeType = input.MimeType
	}
	if input.FileType != "" {
		file.FileType = input.FileType
	}
	if input.Extension != "" {
		file.Extension = input.Extension
	}
	if input.Purpose != "" {
		file.Purpose = &input.Purpose
	}
	if input.Md5 != "" {
		file.Md5 = &input.Md5
	}

	if err := s.repo.Update(ctx, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) DeleteFile(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "file_not_found", "File not found")
	}

	return s.repo.Delete(ctx, id)
}
