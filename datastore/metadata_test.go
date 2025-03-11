package datastore_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	chi "github.com/go-chi/chi/v5"

	store "github.com/marqeta/pr-bot/datastore"
	"github.com/marqeta/pr-bot/id"
)

func TestToMetadata(t *testing.T) {
	payload := []byte(`{"key": "value"}`)
	tests := []struct {
		name    string
		args    *store.Metadata
		payload any
		want    *store.Metadata
		wantErr bool
	}{
		{
			name:    "Should parse metadata from request",
			payload: payload,
			args:    randomMetadata(),
			want:    randomMetadata(),
			wantErr: false,
		},
		{
			name:    "Should Error when owner is missing",
			payload: payload,
			args:    emptyOwner(),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when repo is missing",
			payload: payload,
			args:    emptyRepo(),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when pr number is missing",
			payload: payload,
			args:    emptyPR(),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when service is missing",
			payload: payload,
			args:    emptyService(),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when number is negative",
			args:    negativeNumber(),
			payload: payload,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when head is missing",
			args:    emptyHead(),
			payload: payload,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when base is missing",
			args:    emptyBase(),
			payload: payload,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when job is missing",
			args:    emptyJob(),
			payload: payload,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Should Error when head and base are the same",
			args:    sameHeadAndBase(),
			payload: payload,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRequest(t, tt.args, tt.payload)
			got, err := store.ToMetadata(r)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sameHeadAndBase() *store.Metadata {
	m := randomMetadata()
	m.Head = m.Base
	return m
}

func emptyJob() *store.Metadata {
	m := randomMetadata()
	m.Job = ""
	return m
}

func emptyBase() *store.Metadata {
	m := randomMetadata()
	m.Base = ""
	return m
}

func emptyHead() *store.Metadata {
	m := randomMetadata()
	m.Head = ""
	return m
}

func negativeNumber() *store.Metadata {
	m := randomMetadata()
	m.Number = -1
	return m
}

func randomService() *store.Metadata {
	m := randomMetadata()
	m.Service = "random"
	return m
}

func emptyService() *store.Metadata {
	m := randomMetadata()
	m.Service = ""
	return m
}

func emptyPR() *store.Metadata {
	m := randomMetadata()
	m.Number = 0
	return m
}

func emptyRepo() *store.Metadata {
	m := randomMetadata()
	m.Repo = ""
	return m
}
func emptyOwner() *store.Metadata {
	m := randomMetadata()
	m.Owner = ""
	return m
}

func randomMetadata() *store.Metadata {
	return &store.Metadata{
		PR: id.PR{
			Owner:  "ci",
			Repo:   "ci-kirkland-jobs",
			Number: 1,
		},
		Service: "kirkland",
		Head:    "asdahasdsdjasdjacn1123h12",
		Base:    "1233511adsfaewdqffasdasda",
		Job:     "ci-kirkland-jobs (us-west-2))",
	}
}

func NewRequest(t *testing.T, m *store.Metadata, payload any) *http.Request {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	url := fmt.Sprintf("http://localhost/v1/data/%s/pr/%s/%s/%d?head=%s&base=%s&job=%s",
		m.Service, m.Owner, m.Repo, m.Number, m.Head, m.Base, m.Job)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("owner", m.Owner)
	chiCtx.URLParams.Add("repo", m.Repo)
	chiCtx.URLParams.Add("number", fmt.Sprintf("%d", m.Number))
	chiCtx.URLParams.Add("service", m.Service)
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	return req
}
