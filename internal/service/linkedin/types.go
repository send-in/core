package service

import (
	model "core/internal/model"

	"github.com/google/uuid"
)

type ConnectionsResponse struct {
	Elements []ConnectionElement `json:"elements"`
	Paging   Paging              `json:"paging"`
}

type Paging struct {
	Count int `json:"count"`
	Start int `json:"start"`
}

type ConnectionElement struct {
	CreatedAt  int64       `json:"createdAt"`
	EntityURN  string      `json:"entityUrn"`
	MiniProfile MiniProfile `json:"miniProfile"`
}

type MiniProfile struct {
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Occupation       string  `json:"occupation"`
	PublicIdentifier string  `json:"publicIdentifier"`
	Picture          Picture `json:"picture"`
}

type Picture struct {
	VectorImage VectorImage `json:"com.linkedin.common.VectorImage"`
}

type VectorImage struct {
	RootURL  string     `json:"rootUrl"`
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	Width                         int    `json:"width"`
	Height                        int    `json:"height"`
	FileIdentifyingURLPathSegment string `json:"fileIdentifyingUrlPathSegment"`
}

func LinkedinConnection(
	accountID uuid.UUID,
	element ConnectionElement,
) model.Connection {
	profile := element.MiniProfile
	var picture string

	if len(profile.Picture.VectorImage.Artifacts) > 0 {
		artifact := profile.Picture.
			VectorImage.
			Artifacts[len(profile.Picture.VectorImage.Artifacts)-1]

		picture =
			profile.Picture.VectorImage.RootURL +
			artifact.FileIdentifyingURLPathSegment
	}

	return model.Connection{
		AccountID: &accountID,
		PublicID:  profile.PublicIdentifier,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.Occupation,
		Picture:   picture,
	}
}