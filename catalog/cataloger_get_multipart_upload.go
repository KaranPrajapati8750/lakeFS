package catalog

import (
	"context"

	"github.com/treeverse/lakefs/db"
)

func (c *cataloger) GetMultipartUpload(ctx context.Context, repository string, uploadID string) (*MultipartUpload, error) {
	if err := Validate(ValidateFields{
		"repository": ValidateRepositoryName(repository),
		"uploadID":   ValidateUploadID(uploadID),
	}); err != nil {
		return nil, err
	}

	res, err := c.db.Transact(func(tx db.Tx) (interface{}, error) {
		repoID, err := getRepositoryID(tx, repository)
		if err != nil {
			return nil, err
		}
		var m MultipartUpload
		if err := tx.Get(&m, `
			SELECT r.name as repository, m.upload_id, m.path, m.creation_date, m.physical_address 
			FROM multipart_uploads m, repositories r
			WHERE r.id = m.repository_id AND m.repository_id = $1 AND m.upload_id = $2`,
			repoID, uploadID); err != nil {
			return nil, err
		}
		return &m, nil
	}, c.txOpts(ctx, db.ReadOnly())...)
	if err != nil {
		return nil, err
	}
	return res.(*MultipartUpload), nil
}