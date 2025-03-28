package input_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/datastore"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/opa/input"
)

func Test_adapter_PREventToGHE(t *testing.T) {
	ctx := context.TODO()
	tests := []struct {
		name    string
		event   *github.PullRequestEvent
		want    input.GHE
		wantErr bool
	}{
		{
			name:    "Should convert PR event to GHE",
			event:   prEvent(ref("opened"), sampleID()),
			want:    sampleGHE("pull_request", "opened", sampleID()),
			wantErr: false,
		},
		{
			name:    "Should convert PR event to GHE with nil org",
			event:   prEmptyOrg(sampleID()),
			want:    gheEmptyOrg("pull_request", "opened", sampleID()),
			wantErr: false,
		},
		{
			name:    "Should convert PR event to GHE with nil repo",
			event:   prEmptyRepo(sampleID()),
			want:    gheEmptyRepo("pull_request", "opened", sampleID()),
			wantErr: false,
		},
		{
			name:    "Should convert PR event to GHE with nil PR",
			event:   prEmptyPR(sampleID()),
			want:    gheEmptyPR("pull_request", "opened", sampleID()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := gh.NewMockAPI(t)
			a := input.NewAdapter(mockDAO)
			got, err := a.PREventToGHE(ctx, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("adapter.PREventToGHE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adapter.PREventToGHE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adapter_MetadataToGHE(t *testing.T) {
	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	tests := []struct {
		name            string
		metadata        datastore.Metadata
		setExpectations func(*gh.MockAPI)
		want            input.GHE
		wantErr         bool
	}{
		{
			name:     "Should convert metadata to GHE",
			metadata: sampleMetadata(sampleID()),
			setExpectations: func(mockDAO *gh.MockAPI) {
				mockDAO.EXPECT().GetOrganization(ctx, sampleID()).Return(samplePREvent().Organization, nil)
				mockDAO.EXPECT().GetRepository(ctx, sampleID()).Return(samplePREvent().Repo, nil)
				mockDAO.EXPECT().GetPullRequest(ctx, sampleID()).Return(samplePREvent().PullRequest, nil)
			},
			want:    sampleGHE("kirkland", "ci-kirkland-jobs (us-west-2))", sampleID()),
			wantErr: false,
		},
		{
			name:     "Should throw error if GetOrganization fails",
			metadata: sampleMetadata(sampleID()),
			setExpectations: func(mockDAO *gh.MockAPI) {
				mockDAO.EXPECT().GetOrganization(ctx, sampleID()).Return(nil, randomErr)
			},
			want:    input.GHE{},
			wantErr: true,
		},
		{
			name:     "Should throw error if GetRepository fails",
			metadata: sampleMetadata(sampleID()),
			setExpectations: func(mockDAO *gh.MockAPI) {
				mockDAO.EXPECT().GetOrganization(ctx, sampleID()).Return(samplePREvent().Organization, nil)
				mockDAO.EXPECT().GetRepository(ctx, sampleID()).Return(nil, randomErr)
			},
			want:    input.GHE{},
			wantErr: true,
		},
		{
			name:     "Should throw error if GetPullRequest fails",
			metadata: sampleMetadata(sampleID()),
			setExpectations: func(mockDAO *gh.MockAPI) {
				mockDAO.EXPECT().GetOrganization(ctx, sampleID()).Return(samplePREvent().Organization, nil)
				mockDAO.EXPECT().GetRepository(ctx, sampleID()).Return(samplePREvent().Repo, nil)
				mockDAO.EXPECT().GetPullRequest(ctx, sampleID()).Return(nil, randomErr)
			},
			want:    input.GHE{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := gh.NewMockAPI(t)
			a := input.NewAdapter(mockDAO)
			tt.setExpectations(mockDAO)
			got, err := a.MetadataToGHE(ctx, &tt.metadata)
			if (err != nil) != tt.wantErr {
				t.Errorf("adapter.MetadataToGHE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adapter.MetadataToGHE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adapter_PRReviewEventToGHE(t *testing.T) {
	ctx := context.TODO()
	tests := []struct {
		name    string
		event   *github.PullRequestReviewEvent
		want    input.GHE
		wantErr bool
	}{
		{
			name:    "Should convert PR review event to GHE",
			event:   prrEvent(ref("submitted"), sampleID()),
			want:    sampleGHE("pull_request_review", "submitted", sampleID()),
			wantErr: false,
		},
		{
			name:    "Should convert PR review event to GHE with nil org",
			event:   prrEmptyOrg(sampleID()),
			want:    gheEmptyOrg("pull_request_review", "submitted", sampleID()),
			wantErr: false,
		},
		{
			name:    "Should convert PR review event to GHE with nil repo",
			event:   prrEmptyRepo(sampleID()),
			want:    gheEmptyRepo("pull_request_review", "submitted", sampleID()),
			wantErr: false,
		},
		{
			name:    "Should convert PR review event to GHE with nil PR",
			event:   prrEmptyPR(sampleID()),
			want:    gheEmptyPR("pull_request_review", "submitted", sampleID()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := gh.NewMockAPI(t)
			a := input.NewAdapter(mockDAO)
			got, err := a.PRReviewEventToGHE(ctx, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("adapter.PRReviewEventToGHE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adapter.PRReviewEventToGHE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sampleID() id.PR {

	return id.PR{
		Owner:        "owner1",
		Repo:         "repo1",
		Number:       1,
		NodeID:       "nodeid1",
		RepoFullName: "owner1/repo1",
		Author:       "user1",
		URL:          "owner1/repo1/1",
	}
}

func gheEmptyRepo(event string, action string, id id.PR) input.GHE {
	ghe := sampleGHE(event, action, id)
	ghe.Repository = nil
	return ghe
}

func gheEmptyOrg(event string, action string, id id.PR) input.GHE {
	ghe := sampleGHE(event, action, id)
	ghe.Organization = nil
	return ghe
}

func gheEmptyPR(event string, action string, id id.PR) input.GHE {
	ghe := sampleGHE(event, action, id)
	ghe.PullRequest = nil
	return ghe
}

func sampleGHE(event string, action string, id id.PR) input.GHE {
	prEvent := prEvent(&action, id)
	return input.GHE{
		Event:        event,
		Action:       action,
		PullRequest:  prEvent.PullRequest,
		Repository:   prEvent.Repo,
		Organization: prEvent.Organization,
	}
}

func prEmptyOrg(id id.PR) *github.PullRequestEvent {
	event := prEvent(ref("opened"), id)
	event.Organization = nil
	return event
}

func prEmptyRepo(id id.PR) *github.PullRequestEvent {
	event := prEvent(ref("opened"), id)
	event.Repo = nil
	return event
}

func prEmptyPR(id id.PR) *github.PullRequestEvent {
	event := prEvent(ref("opened"), id)
	event.PullRequest = nil
	return event
}

func samplePREvent() *github.PullRequestEvent {
	return prEvent(ref("opened"), sampleID())
}

func prEvent(action *string, id id.PR) *github.PullRequestEvent {
	url := fmt.Sprintf("%s/%s/%d", id.Owner, id.Repo, id.Number)
	return &github.PullRequestEvent{
		Action: action,
		PullRequest: &github.PullRequest{
			Number: &id.Number,
			NodeID: &id.NodeID,
			User: &github.User{
				Login: &id.Author,
			},
			HTMLURL: &url,
			Base: &github.PullRequestBranch{
				Repo: &github.Repository{
					AllowAutoMerge:   ref(false),
					AllowRebaseMerge: ref(false),
					AllowSquashMerge: ref(false),
					AllowMergeCommit: ref(false),
				},
			},
			ChangedFiles: github.Int(1),
		},
		Repo: &github.Repository{
			Owner: &github.User{
				Login: &id.Owner,
			},
			DefaultBranch: github.String("main"),
			Name:          &id.Repo,
			FullName:      &id.RepoFullName,
			Visibility:    github.String("public"),
		},
		Organization: &github.Organization{
			Login: &id.Owner,
		},
	}
}

func prrEmptyOrg(id id.PR) *github.PullRequestReviewEvent {
	event := prrEvent(ref("submitted"), id)
	event.Organization = nil
	return event
}

func prrEmptyRepo(id id.PR) *github.PullRequestReviewEvent {
	event := prrEvent(ref("submitted"), id)
	event.Repo = nil
	return event
}

func prrEmptyPR(id id.PR) *github.PullRequestReviewEvent {
	event := prrEvent(ref("submitted"), id)
	event.PullRequest = nil
	return event
}

func prrEvent(action *string, id id.PR) *github.PullRequestReviewEvent {
	prEvent := prEvent(action, id)
	prr := &github.PullRequestReviewEvent{
		Action:       action,
		PullRequest:  prEvent.PullRequest,
		Repo:         prEvent.Repo,
		Organization: prEvent.Organization,
	}
	return prr
}

func ref[T any](x T) *T {
	return &x
}

func sampleMetadata(pr id.PR) datastore.Metadata {
	return datastore.Metadata{
		PR:      pr,
		Service: "kirkland",
		Head:    "asdahasdsdjasdjacn1123h12",
		Base:    "1233511adsfaewdqffasdasda",
		Job:     "ci-kirkland-jobs (us-west-2))",
	}
}
