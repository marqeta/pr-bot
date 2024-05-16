package plugins

import (
	"context"
	"encoding/json"

	"github.com/go-chi/httplog"
	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
)

type FilesChanged struct {
	dao       gh.API
	sizeLimit int
}

// GetMessage implements input.Plugin.
func (fc *FilesChanged) GetInputMsg(ctx context.Context, ghe input.GHE) (json.RawMessage, error) {
	oplog := httplog.LogEntry(ctx)
	filesChanged, err := fc.dao.ListFilesChangedInPR(ctx, ghe.ToID())
	if err != nil {
		return json.RawMessage{}, err
	}
	smallFilesChanged := make([]*github.CommitFile, 0)
	size := 0
	isSkipped := false
	for _, file := range filesChanged {
		if (size + len(file.GetPatch())) > fc.sizeLimit {
			isSkipped = true
			continue
		}
		smallFilesChanged = append(smallFilesChanged, file)
		size += len(file.GetPatch())
	}
	if isSkipped {
		oplog.Info().Msg("Some files are skipped in input payload since it exceeds the size limit")
	}
	data, err := json.Marshal(smallFilesChanged)
	if err != nil {
		return json.RawMessage{}, err
	}
	return json.RawMessage(data), nil
}

// Name implements input.Plugin.
func (fc *FilesChanged) Name() string {
	return "files_changed"
}

func NewFilesChanged(dao gh.API, sizeLimit int) input.Plugin {
	return &FilesChanged{dao: dao, sizeLimit: sizeLimit}
}
