package asynqueue

import "context"

type HookFunc func(ctx context.Context, taskType string, payload any) error
