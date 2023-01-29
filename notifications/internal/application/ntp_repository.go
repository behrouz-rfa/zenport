package application

import (
	"context"
)

type NtpRepository interface {
	Created(ctx context.Context) error
}
