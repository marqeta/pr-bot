package configstore

import (
	"context"
	"fmt"

	"github.com/marqeta/pr-bot/metrics"
)

type DynamicConfig interface {
	// a hook which would be called everytime the config is updated.
	Update() error
}

type Getter[T DynamicConfig] interface {
	Get() (T, error)
	Close()
}

func emitSuccess(m metrics.Emitter, name string) {
	m.EmitDist(context.Background(), "config.load.success", 1, []string{
		fmt.Sprintf("name:%s", name),
	})
}

func emitError(m metrics.Emitter, name string, errCode string) {
	m.EmitDist(context.Background(), "config.load.error", 1, []string{
		fmt.Sprintf("name:%s", name),
		fmt.Sprintf("code:%s", errCode),
	})
}
