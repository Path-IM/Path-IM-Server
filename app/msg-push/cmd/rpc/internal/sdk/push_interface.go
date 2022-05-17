package sdk

import "context"

type OfflinePusher interface {
	Push(ctx context.Context, userIDList []string, alert, detailContent string) (resp string, err error)
}
