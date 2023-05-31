package repository

import (
	"context"

	"github.com/issy20/go-oidc-client/domain/entity"
)

type IChannelRepository interface {
	GetFollowedChannel(ctx context.Context, userID string) (*entity.GetFollowedChannelResponse, error)
}
