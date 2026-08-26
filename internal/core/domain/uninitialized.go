package domain

import (
	"time"

	"github.com/google/uuid"
)

var (
	UninitializedUUID = uuid.Nil
	UninitializedRole = "user"
	UninitializedTime = time.Time{}
)
