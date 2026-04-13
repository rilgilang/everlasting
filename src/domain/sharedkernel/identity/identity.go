package identity

import "github.com/google/uuid"

type ID struct {
	uuid.UUID
}

func (id ID) IsNil() bool {
	return id.UUID == uuid.Nil
}

// NewID represent constructor
func NewID() ID {
	id, _ := uuid.NewRandom()
	return ID{UUID: id}
}

// NewZeroID represent constructor
func NewZeroID() ID {
	return ID{UUID: uuid.Nil}
}

// Generate new id from string
func FromStringOrNil(val string) (id ID) {
	u, err := uuid.Parse(val)
	if err != nil {
		return NewZeroID()
	}

	return ID{
		UUID: u,
	}
}
