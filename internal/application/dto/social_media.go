// Package dto
package dto

type SocialMedia struct {
	IDSocialMedia uint   `json:"id_social_media"`
	IDEnterprise  uint   `json:"id_enterprise"`
	Name          string `json:"name"`
	URL           string `json:"url"`
}
