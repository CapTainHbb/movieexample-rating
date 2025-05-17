package rating

import (
	"context"
	"errors"
	"github.com/captainhbb/movieexample-rating/internal/controller"
	"github.com/captainhbb/movieexample-rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type Controller struct {
	repo ratingRepository
}

func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (r *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := r.repo.Get(ctx, recordID, recordType)
	if err != nil && !errors.Is(err, controller.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}

	return sum / float64(len(ratings)), nil
}

func (r *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return r.repo.Put(ctx, recordID, recordType, rating)
}
