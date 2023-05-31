package entity

type FollowedChannel struct {
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterLogin string `json:"broadcaster_login"`
	BroadcasterName  string `json:"broadcaster_name"`
	FollowedAt       string `json:"followed_at"`
}

type GetFollowedChannelResponse struct {
	Data       []FollowedChannel      `json:"data"`
	Total      int                    `json:"total"`
	Pagination map[string]interface{} `json:"pagination"`
}

type MyChannel struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

type GetMyChannelResponse struct {
	Data []MyChannel `json:"data"`
}
