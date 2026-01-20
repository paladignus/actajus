// Package dto
package dto

type SocialMedia struct {
	IDSocialMedia uint   `json:"id_social_media"`
	IDCompany     uint   `json:"id_company"`
	Name          string `json:"name"`
	URL           string `json:"url"`
}
