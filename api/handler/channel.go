package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/issy20/go-oidc-client/usecase"
)

type ChannelHandler struct {
	cu usecase.IChannelUsecase
}

func NewChannelHandler(cu usecase.IChannelUsecase) *ChannelHandler {
	return &ChannelHandler{
		cu: cu,
	}
}

func (ch *ChannelHandler) GetFollowedChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")

	channels, err := ch.cu.GetFollowedChannel(ctx, userID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(channels)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
