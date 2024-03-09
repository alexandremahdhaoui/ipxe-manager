package types

import (
	"github.com/google/uuid"
)

// ---------------------------------------------------- Assignment ----------------------------------------------------- //

type Assignment struct {
	ProfileName      string
	SubjectSelectors SubjectSelectors
}

type SubjectSelectors struct {
	UUIDs     []uuid.UUID
	Buildarch string
}
