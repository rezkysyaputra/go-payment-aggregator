package pkg

import "github.com/google/uuid"

func GenerateUUIDV7() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
