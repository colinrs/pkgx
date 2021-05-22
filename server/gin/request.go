package gin

import "context"

type Request interface {
	Validator(ctx context.Context) error
}