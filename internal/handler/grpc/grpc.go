package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/captainhbb/movieexample-protoapis/gen"
	"movieexample-rating/internal/controller"
	"movieexample-rating/pkg/model"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	v, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}

	err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId),
		model.RecordType(req.RecordType),
		&model.Rating{UserID: model.UserID(req.UserId),
			Value: model.RatingValue(req.RatingValue)})
	if err != nil {
		return nil, err
	}

	return &gen.PutRatingResponse{}, nil
}
