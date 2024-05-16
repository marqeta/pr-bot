package types_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/marqeta/pr-bot/opa/types"
)

func TestReviewType_Precedence(t *testing.T) {
	tests := []struct {
		name string
		rt   types.ReviewType
		want uint8
	}{
		{
			name: "SKIP",
			rt:   types.Skip,
			want: 0,
		},
		{
			name: "Approve",
			rt:   types.Approve,
			want: 1,
		},
		{
			name: "COMMENT",
			rt:   types.Comment,
			want: 2,
		},
		{
			name: "REQUEST_CHANGES",
			rt:   types.RequestChanges,
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uint8(tt.rt); got != tt.want {
				t.Errorf("ReviewType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestReviewType_String(t *testing.T) {
	tests := []struct {
		name string
		rt   types.ReviewType
		want string
	}{
		{
			name: "Approve",
			rt:   types.Approve,
			want: "APPROVE",
		},
		{
			name: "REQUEST_CHANGES",
			rt:   types.RequestChanges,
			want: "REQUEST_CHANGES",
		},
		{
			name: "COMMENT",
			rt:   types.Comment,
			want: "COMMENT",
		},
		{
			name: "SKIP",
			rt:   types.Skip,
			want: "SKIP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.String(); got != tt.want {
				t.Errorf("ReviewType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseReviewType(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    types.ReviewType
		wantErr bool
	}{
		{
			name:    "Should parse APPROVE",
			args:    "APPROVE",
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should parse REQUEST_CHANGES",
			args:    "REQUEST_CHANGES",
			want:    types.RequestChanges,
			wantErr: false,
		},
		{
			name:    "Should parse COMMENT",
			args:    "COMMENT",
			want:    types.Comment,
			wantErr: false,
		},
		{
			name:    "Should parse SKIP",
			args:    "SKIP",
			want:    types.Skip,
			wantErr: false,
		},
		{
			name:    "Should parse approve",
			args:    "approve",
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should parse request_changes",
			args:    "request_changes",
			want:    types.RequestChanges,
			wantErr: false,
		},
		{
			name:    "Should parse comment",
			args:    "comment",
			want:    types.Comment,
			wantErr: false,
		},
		{
			name:    "Should parse skip",
			args:    "skip",
			want:    types.Skip,
			wantErr: false,
		},
		{
			name:    "Should trim whitespace and parse",
			args:    "  APPROVE  ",
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should throw error for invalid review type",
			args:    "LGTM",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.ParseReviewType(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReviewType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseReviewType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseReviewState(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    types.ReviewType
		wantErr bool
	}{
		{
			name:    "Should parse APPROVED",
			args:    "APPROVED",
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should parse CHANGES_REQUESTED",
			args:    "CHANGES_REQUESTED",
			want:    types.RequestChanges,
			wantErr: false,
		},
		{
			name:    "Should parse COMMENTED",
			args:    "COMMENTED",
			want:    types.Comment,
			wantErr: false,
		},
		{
			name:    "Should parse SKIP",
			args:    "SKIP",
			want:    types.Skip,
			wantErr: false,
		},
		{
			name:    "Should parse approved",
			args:    "approved",
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should parse changes_requested",
			args:    "changes_requested",
			want:    types.RequestChanges,
			wantErr: false,
		},
		{
			name:    "Should parse commented",
			args:    "commented",
			want:    types.Comment,
			wantErr: false,
		},
		{
			name:    "Should parse skip",
			args:    "skip",
			want:    types.Skip,
			wantErr: false,
		},
		{
			name:    "Should trim whitespace and parse",
			args:    "  APPROVED  ",
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should throw error for invalid review state",
			args:    "LGTM",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.ParseReviewState(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReviewState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseReviewState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReviewType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		r       types.ReviewType
		want    []byte
		wantErr bool
	}{
		{
			name:    "Should marshal APPROVE",
			r:       types.Approve,
			want:    []byte(`"APPROVE"`),
			wantErr: false,
		},
		{
			name:    "Should marshal REQUEST_CHANGES",
			r:       types.RequestChanges,
			want:    []byte(`"REQUEST_CHANGES"`),
			wantErr: false,
		},
		{
			name:    "Should marshal COMMENT",
			r:       types.Comment,
			want:    []byte(`"COMMENT"`),
			wantErr: false,
		},
		{
			name:    "Should marshal SKIP",
			r:       types.Skip,
			want:    []byte(`"SKIP"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReviewType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReviewType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReviewType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    types.ReviewType
		wantErr bool
	}{
		{
			name:    "Should unmarshal SKIP",
			data:    []byte(`"SKIP"`),
			want:    types.Skip,
			wantErr: false,
		},
		{
			name:    "Should unmarshal COMMENT",
			data:    []byte(`"COMMENT"`),
			want:    types.Comment,
			wantErr: false,
		},
		{
			name:    "Should unmarshal REQUEST_CHANGES",
			data:    []byte(`"REQUEST_CHANGES"`),
			want:    types.RequestChanges,
			wantErr: false,
		},
		{
			name:    "Should unmarshal APPROVE",
			data:    []byte(`"APPROVE"`),
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Should unmarshal skip",
			data:    []byte(`"skip"`),
			want:    types.Skip,
			wantErr: false,
		},
		{
			name:    "Should unmarshal comment",
			data:    []byte(`"comment"`),
			want:    types.Comment,
			wantErr: false,
		},
		{
			name:    "Should unmarshal request_changes",
			data:    []byte(`"request_changes"`),
			want:    types.RequestChanges,
			wantErr: false,
		},
		{
			name:    "Should unmarshal approve",
			data:    []byte(`"approve"`),
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Unmarshal should trim spaces",
			data:    []byte(`"   approve   "`),
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Unmarshal should trim spaces",
			data:    []byte(`"   approve   "`),
			want:    types.Approve,
			wantErr: false,
		},
		{
			name:    "Unmarshal should throw error for invalid review type",
			data:    []byte(`"LGTM"`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r types.ReviewType
			err := r.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReviewType.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("ReviewType.UnmarshalJSON() = %v, want %v", r, tt.want)
			}
		})
	}
}

func TestReview_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    types.Review
		wantErr bool
	}{
		{
			name:    "Should unmarshal review",
			data:    []byte(`{"type":"APPROVE","body":"LGTM!! :100: :rocket: :tada:"}`),
			want:    types.Review{Type: types.Approve, Body: "LGTM!! :100: :rocket: :tada:"},
			wantErr: false,
		},
		{
			name:    "Should throw error for invalid review type",
			data:    []byte(`{"type":"APPROVED","body":"LGTM!! :100: :rocket: :tada:"}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r types.Review
			err := json.Unmarshal(tt.data, &r)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("json.Unmarshal() = %v, want %v", r, tt.want)
			}
		})
	}
}
