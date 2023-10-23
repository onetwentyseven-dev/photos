package mysql

import (
	"context"
	"photos"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
)

var imageTableColumns = []string{
	"id", "user_id", "name", "description", "status", "processing_errors", "image_exif_data", "ts_created", "ts_updated",
}

type ImageRepository struct {
	db *sqlx.DB
}

func NewImageRepository(db *sqlx.DB) *ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

// ImagesByUserID
func (r *ImageRepository) ImagesByUserID(ctx context.Context, userID string) ([]*photos.Image, error) {

	var images []*photos.Image
	query, args, err := sq.Select(imageTableColumns...).From("images").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.SelectContext(ctx, &images, query, args...)
	if err != nil {
		return nil, err
	}

	return images, err

}

func (r *ImageRepository) Image(ctx context.Context, id string) (*photos.Image, error) {

	var image photos.Image
	query, args, err := sq.Select(imageTableColumns...).From("images").Where(sq.Eq{"id": id}).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &image, query, args...)
	if err != nil {
		return nil, err
	}

	return &image, err
}

func (r *ImageRepository) CreateImage(ctx context.Context, image *photos.Image) error {

	now := time.Now()
	image.TSCreated = now
	image.TSUpdated = now

	_, err := sq.Insert("images").SetMap(structs.Map(image)).RunWith(r.db).ExecContext(ctx)
	return err
}

// UpdateImageByImageID
func (r *ImageRepository) UpdateImageByImageID(ctx context.Context, image *photos.Image) error {

	now := time.Now()
	image.TSUpdated = now

	query, args, err := sq.Update("images").SetMap(structs.Map(image)).Where(sq.Eq{"id": image.ID}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteImageByImageID
func (r *ImageRepository) DeleteImageByImageID(ctx context.Context, imageID string) error {

	query, args, err := sq.Delete("images").Where(sq.Eq{"id": imageID}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
