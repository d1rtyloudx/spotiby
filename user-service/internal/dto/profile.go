package dto

import "github.com/d1rtyloudx/spotiby-pkg/lib"

type Profile struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Description  string `json:"description"`
	CredentialID string `json:"credential_id"`
	AvatarURL    string `json:"avatar_url"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

type UpdateAvatarProfileRequest struct {
	ID        string `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

type FollowingRequest struct {
	TargetID string `json:"target_id"`
}

type PagedProfileResponse struct {
	Profiles   []Profile              `json:"profiles"`
	Pagination lib.PaginationResponse `json:"pagination"`
}
