package pkg

import "github.com/google/uuid"

func GenerateUUIDV7() uuid.UUID {
	id, _ := uuid.NewV7()

	return id
}
