package oci

import (
	"context"
	"reflect"
	"testing"
)

func Test_bundleFileReader_FilterModules(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		dirs []string
		want []string
	}{
		{
			name: "Should return single modules",
			dirs: []string{"/abc/def"},
			want: []string{"/abc/def"},
		},
		{
			name: "Should return multiple modules",
			dirs: []string{"/abc/def", "/abcd/def"},
			want: []string{"/abc/def", "/abcd/def"},
		},
		{
			name: "Should filter dirs at lvl 1",
			dirs: []string{"/abc", "/abcd/def"},
			want: []string{"/abcd/def"},
		},
		{
			name: "Should filter dirs at lvl > 2",
			dirs: []string{"/abc/d/e", "/abc/d/e/f", "/abcd/def"},
			want: []string{"/abcd/def"},
		},
		{
			name: "Should filter root /",
			dirs: []string{"/", "/abcd/def"},
			want: []string{"/abcd/def"},
		},
		{
			name: "Should return empty when dirs are empty",
			dirs: []string{"/"},
			want: []string{},
		},
		{
			name: "Should skip hidden dirs at lvl 1",
			dirs: []string{"/_a/b", "/abcd/def"},
			want: []string{"/abcd/def"},
		},
		{
			name: "Should skip hidden dirs at lvl 2",
			dirs: []string{"/a/_b", "/abcd/def"},
			want: []string{"/abcd/def"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &bundleFileReader{}
			if got := r.FilterModules(ctx, tt.dirs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bundleFileReader.FilterModules() = %v, want %v", got, tt.want)
			}
		})
	}
}
