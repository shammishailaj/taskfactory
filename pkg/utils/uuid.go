package utils

import (
	"github.com/google/uuid"
)

func (u *Utils) NewUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}
