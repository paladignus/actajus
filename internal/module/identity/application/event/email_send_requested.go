// Package event
package event

import "time"

type EmailSendRequested struct {
	IDMessage string            `json:"id_message"`
	Template  string            `json:"template"`
	To        string            `json:"to"`
	Data      map[string]any    `json:"data"`
	Meta      map[string]string `json:"meta,omitzero"`
	SentAfter *time.Time        `json:"sent_after,omitzero"`
}
