// Package readmodel
package readmodel

import "time"

type SocialMediaReadModel struct {
	ID        int64
	IDCompany int64
	Platform  string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
