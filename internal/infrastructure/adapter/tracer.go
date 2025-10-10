package adapter

import "github.com/google/uuid"

func GenerateTraceID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
