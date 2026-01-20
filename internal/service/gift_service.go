package service

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
)

type IGiftService interface {
	CreateGift(ctx context.Context, gift domain.Gift) (*domain.Gift, error)
	ListGifts(ctx context.Context, pagination common.Pagination, sort common.Sort) ([]domain.Gift, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Gift, error)
}

type GiftService struct {
	giftRepo IGiftRepository
}

func NewGiftService(giftRepo IGiftRepository) *GiftService {
	return &GiftService{
		giftRepo: giftRepo,
	}
}

func (s *GiftService) CreateGift(ctx context.Context, gift domain.Gift) (*domain.Gift, error) {
	return s.giftRepo.Create(ctx, gift)
}

func (s *GiftService) ListGifts(ctx context.Context, pagination common.Pagination, sort common.Sort) ([]domain.Gift, int64, error) {
	gifts, err := s.giftRepo.FindAll(ctx, pagination, sort)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.giftRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return gifts, total, nil
}

func (s *GiftService) GetByID(ctx context.Context, id int64) (*domain.Gift, error) {
	return s.giftRepo.GetByID(ctx, id)
}
