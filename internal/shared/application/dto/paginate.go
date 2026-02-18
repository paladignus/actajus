// Package dto
package dto

type PageInfo struct {
	HasNextPage     bool    `json:"has_next_page"`
	HasPreviousPage bool    `json:"has_previous_page"`
	NextURL         *string `json:"next_url,omitzero"`
	PreviousURL     *string `json:"previous_url,omitzero"`
}
