package model

const DefaultPhotoID = "default-photo"

type Photo struct {
	ID            string `json:"id"`
	Base64Content string `json:"base64Content"`
}
