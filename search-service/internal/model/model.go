package model

type Profile struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Description  string `json:"description"`
	CredentialID string `json:"credential_id"`
	AvatarURL    string `json:"avatar_url"`
}

type DeleteProfile struct {
	ID string `json:"id"`
}

type Playlist struct {
	ID string `json:"id"`
}

type DeletePlaylist struct {
	ID string `json:"id"`
}

type Track struct {
	ID string `json:"id"`
}

type DeleteTrack struct {
	ID string `json:"id"`
}
