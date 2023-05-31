package usecase

import (
	"context"

	"github.com/issy20/go-oidc-client/domain/entity"
	"github.com/issy20/go-oidc-client/domain/repository"
	"golang.org/x/xerrors"
)

var _ IChannelUsecase = &ChannelUsecase{}

type ChannelUsecase struct {
	cr repository.IChannelRepository
}

type IChannelUsecase interface {
	GetFollowedChannel(ctx context.Context, userID string) (*entity.GetFollowedChannelResponse, error)
}

func NewChannelUsecase(cr repository.IChannelRepository) IChannelUsecase {
	return &ChannelUsecase{
		cr: cr,
	}
}

func (cu *ChannelUsecase) GetFollowedChannel(ctx context.Context, userID string) (*entity.GetFollowedChannelResponse, error) {
	channles, err := cu.cr.GetFollowedChannel(ctx, userID)
	if err != nil {
		return nil, xerrors.Errorf("ChannelUsecase.GetFollowedChannel Error : %w", err)
	}
	return channles, nil
}
