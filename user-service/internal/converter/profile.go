package converter

import (
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/d1rtyloudx/spotiby/user-service/internal/dto"
)

func ProfileToProfileDTO(profile model.Profile) dto.Profile {
	return dto.Profile{
		ID:           profile.ID,
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		Description:  profile.Description,
		CredentialID: profile.CredentialID,
		AvatarURL:    profile.AvatarURL,
	}
}

func ProfileDTOToProfile(profile dto.Profile) model.Profile {
	return model.Profile{
		ID:           profile.ID,
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		Description:  profile.Description,
		CredentialID: profile.CredentialID,
		AvatarURL:    profile.AvatarURL,
	}
}

func ProfileListToProfileDTO(profiles []model.Profile) []dto.Profile {
	res := make([]dto.Profile, 0, len(profiles))

	for _, profile := range profiles {
		res = append(res, ProfileToProfileDTO(profile))
	}

	return res
}
