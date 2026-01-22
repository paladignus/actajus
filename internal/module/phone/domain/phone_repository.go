// Package domain
package domain

import (
	"context"
)

type PhoneRepository interface {
	Create(ctx context.Context, phone *Phone) error
}
