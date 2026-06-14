package repository

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/uploadfile"
	"oa-nsdiy/backend/internal/domain"

	"entgo.io/ent/dialect/sql"
)

type FileRepository struct {
	client *ent.Client
}

func NewFileRepository(client *ent.Client) *FileRepository {
	return &FileRepository{client: client}
}

// Type aliases for backward compatibility
type UploadFile = domain.UploadFile
type UploadFileDetail = domain.UploadFileDetail

func toUploadFile(e *ent.UploadFile) *UploadFile {
	return &UploadFile{
		ID:               e.ID,
		Filename:         e.Filename,
		OriginalFilename: e.OriginalFilename,
		FilePath:         e.FilePath,
		FileSize:         e.FileSize,
		MimeType:         e.MimeType,
		FileType:         e.FileType,
		Extension:        e.Extension,
		UploaderID:       e.UploaderID,
		Purpose:          stringPtr(e.Purpose),
		Md5:              stringPtr(e.Md5),
		ReferenceCount:   e.ReferenceCount,
		CreatedAt:        e.CreatedAt,
	}
}

func toUploadFileDetail(e *ent.UploadFile) *UploadFileDetail {
	d := &UploadFileDetail{
		UploadFile: *toUploadFile(e),
	}
	if e.Edges.Uploader != nil {
		d.UploaderName = e.Edges.Uploader.Username
		d.UploaderNickname = stringPtr(e.Edges.Uploader.Nickname)
	}
	return d
}

func toUploadFileDetails(es []*ent.UploadFile) []*UploadFileDetail {
	result := make([]*UploadFileDetail, len(es))
	for i, e := range es {
		result[i] = toUploadFileDetail(e)
	}
	return result
}

func (r *FileRepository) GetByID(ctx context.Context, id int) (*UploadFile, error) {
	e, err := r.client.UploadFile.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUploadFile(e), nil
}

func (r *FileRepository) GetDetailByID(ctx context.Context, id int) (*UploadFileDetail, error) {
	e, err := r.client.UploadFile.Query().
		Where(uploadfile.IDEQ(id)).
		WithUploader().
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toUploadFileDetail(e), nil
}

func (r *FileRepository) Create(ctx context.Context, f *UploadFile) error {
	e, err := r.client.UploadFile.Create().
		SetFilename(f.Filename).
		SetOriginalFilename(f.OriginalFilename).
		SetFilePath(f.FilePath).
		SetFileSize(f.FileSize).
		SetMimeType(f.MimeType).
		SetFileType(f.FileType).
		SetExtension(f.Extension).
		SetUploaderID(f.UploaderID).
		SetNillablePurpose(f.Purpose).
		SetNillableMd5(f.Md5).
		SetReferenceCount(f.ReferenceCount).
		Save(ctx)
	if err != nil {
		return err
	}
	f.ID = e.ID
	return nil
}

func (r *FileRepository) Update(ctx context.Context, f *UploadFile) error {
	_, err := r.client.UploadFile.UpdateOneID(f.ID).
		SetFilename(f.Filename).
		SetOriginalFilename(f.OriginalFilename).
		SetFilePath(f.FilePath).
		SetFileSize(f.FileSize).
		SetMimeType(f.MimeType).
		SetFileType(f.FileType).
		SetExtension(f.Extension).
		SetNillablePurpose(f.Purpose).
		SetNillableMd5(f.Md5).
		Save(ctx)
	return err
}

func (r *FileRepository) Delete(ctx context.Context, id int) error {
	return r.client.UploadFile.DeleteOneID(id).Exec(ctx)
}

func (r *FileRepository) List(ctx context.Context, keyword, fileType string, uploaderID *int, page, pageSize int) ([]*UploadFileDetail, int64, error) {
	q := r.client.UploadFile.Query()

	if keyword != "" {
		q.Where(uploadfile.Or(
			uploadfile.FilenameContains(keyword),
			uploadfile.OriginalFilenameContains(keyword),
		))
	}
	if fileType != "" {
		q.Where(uploadfile.FileTypeEQ(fileType))
	}
	if uploaderID != nil {
		q.Where(uploadfile.UploaderIDEQ(*uploaderID))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	es, err := q.
		Order(uploadfile.ByID(sql.OrderDesc())).
		Limit(pageSize).
		Offset(offset).
		WithUploader().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toUploadFileDetails(es), int64(total), nil
}
