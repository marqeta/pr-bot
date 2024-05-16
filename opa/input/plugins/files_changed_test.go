package plugins_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/input/plugins"
)

func TestFilesChanged_GetMessage(t *testing.T) {

	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	type args struct {
		ghe            input.GHE
		setExpectaions func(d *gh.MockAPI)
		sizeLimit      int
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Should return single file changed",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().ListFilesChangedInPR(ctx, randomGHE().ToID()).
						Return(randomFilesChanged(1), nil)
				},
				sizeLimit: 100,
			},
			want:    toJSON(t, randomFilesChanged(1)),
			wantErr: false,
		},
		{
			name: "Should return multiple files changed",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().ListFilesChangedInPR(ctx, randomGHE().ToID()).
						Return(randomFilesChanged(5), nil)
				},
				sizeLimit: 100,
			},
			want:    toJSON(t, randomFilesChanged(5)),
			wantErr: false,
		},
		{
			name: "Should skip single file with large patch",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().ListFilesChangedInPR(ctx, randomGHE().ToID()).
						Return(randomFilesChanged(1), nil)
				},
				sizeLimit: 0,
			},
			want:    toJSON(t, []*github.CommitFile{}),
			wantErr: false,
		},
		{
			name: "Should skip files with large patch",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().ListFilesChangedInPR(ctx, randomGHE().ToID()).
						Return(slices.Concat(large(22), randomFilesChanged(3), large(1)), nil)
				},
				sizeLimit: 7 * 3, // 7 is the size of each file in randomFilesChanged(3)
			},
			want:    toJSON(t, randomFilesChanged(3)),
			wantErr: false,
		},
		{
			name: "Should skip files after total size is > limit",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().ListFilesChangedInPR(ctx, randomGHE().ToID()).
						Return(slices.Concat(large(20), large(10), randomFilesChanged(5), large(2)), nil)
				},
				sizeLimit: 20 + 7, // 20 is the size of large(1) and 7 is the size of each file in randomFilesChanged(5)
			},
			want:    toJSON(t, slices.Concat(large(20), randomFilesChanged(1))),
			wantErr: false,
		},
		{
			name: "Should return error when dao returns error",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().ListFilesChangedInPR(ctx, randomGHE().ToID()).
						Return(nil, randomErr)
				},
				sizeLimit: 100,
			},
			want:    json.RawMessage([]byte{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := gh.NewMockAPI(t)
			tt.args.setExpectaions(d)
			fc := plugins.NewFilesChanged(d, tt.args.sizeLimit)
			got, err := fc.GetInputMsg(ctx, tt.args.ghe)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesChanged.GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilesChanged.GetMessage() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func toJSON(t *testing.T, f interface{}) json.RawMessage {
	bytes, err := json.Marshal(f)
	if err != nil {
		t.Errorf("error marshalling CommitFiles into json %v", err)
		return nil
	}
	return json.RawMessage(bytes)
}

func randomGHE() input.GHE {
	return input.GHE{
		Event:  "random Event",
		Action: "random Action",
		PullRequest: &github.PullRequest{
			Number: aws.Int(1),
			NodeID: aws.String("random NodeID"),
			User: &github.User{
				Login: aws.String("random Login"),
			},
			Base: &github.PullRequestBranch{
				Ref: aws.String("random Ref"),
			},
		},
		Repository: &github.Repository{
			Name:     aws.String("random Name"),
			Owner:    &github.User{Login: aws.String("random Login")},
			FullName: aws.String("random FullName"),
		},
	}
}

func randomFilesChanged(n int) []*github.CommitFile {

	files := make([]*github.CommitFile, 0, n)
	for i := 0; i < n; i++ {
		files = append(files, &github.CommitFile{
			Filename: aws.String(fmt.Sprintf("file-%d", i)),
			Patch:    aws.String(fmt.Sprintf("patch-%d", i)),
		})
	}
	return files
}

func large(n int) []*github.CommitFile {
	s := ""
	for i := 0; i < n; i++ {
		s += "a"
	}
	return []*github.CommitFile{
		{
			Filename: aws.String("large-file"),
			Patch:    &s,
		},
	}
}
