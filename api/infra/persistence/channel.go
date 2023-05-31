package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/issy20/go-oidc-client/domain/dto"
	"github.com/issy20/go-oidc-client/domain/entity"
	"github.com/issy20/go-oidc-client/domain/repository"
	"github.com/issy20/go-oidc-client/middleware"
	"github.com/issy20/go-oidc-client/util"
)

var _ repository.IChannelRepository = &ChannelRepository{}

type ChannelRepository struct{}

func NewChannelRepository() repository.IChannelRepository {
	return &ChannelRepository{}
}

func (cr *ChannelRepository) GetFollowedChannel(ctx context.Context, userID string) (*entity.GetFollowedChannelResponse, error) {
	url, err := util.CreateURL(dto.GetFollowedChannelEndpoint, "user_id", userID)
	if err != nil {
		return nil, fmt.Errorf("cannot create url: %w", err)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request err : %w", err)
	}
	authorizationHeader, ok := ctx.Value(middleware.AutorizationContextKey).(*middleware.AuthorizationHeader)
	if !ok {
		return nil, fmt.Errorf("cannot convert to *middleware.AuthorizationHeader")
	}

	authorization := util.AddBearer(authorizationHeader.AccessToken)

	req.Header.Set("Authorization", authorization)
	req.Header.Set("Client-Id", authorizationHeader.ClientID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request client.do err : %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var channels *entity.GetFollowedChannelResponse
	err = json.Unmarshal(body, &channels)
	if err != nil {
		return nil, fmt.Errorf("unmarshal err : %w", err)
	}

	log.Print(channels)

	return channels, nil
}
