package kv_test

import (
	"fmt"
	"path"
	"reflect"
	"strings"
	"testing"

	api2 "rvault/internal/pkg/api"
	"rvault/internal/pkg/kv"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func (tc *testClient) List(searchPath string) (*api.Secret, error) {
	var relPath string
	if tc.listErr != nil {
		return nil, tc.listErr
	}
	if tc.kvVersion == "2" {
		relPath = strings.TrimPrefix(searchPath, path.Join(tc.engine+"/metadata"))
	} else {
		relPath = strings.TrimPrefix(searchPath, path.Clean(tc.engine))
	}

	if paths, ok := tc.secretPaths[relPath]; ok {
		if len(paths) == 1 && paths[0] == "" {
			return nil, nil
		} else {
			return &api.Secret{
				Data: map[string]interface{}{
					"keys": paths,
				},
			}, nil
		}
	} else {
		return nil, fmt.Errorf("searchPath %s doesn't exist", searchPath)
	}

}

func TestRList(t *testing.T) {
	type args struct {
		c            api2.VaultClient
		engine       string
		path         string
		concurrency  uint32
		includePaths []string
		excludePaths []string
	}
	tests := []struct {
		name       string
		args       args
		viperFlags map[string]interface{}
		want       []string
		wantErr    bool
	}{
		{
			name: "Smoke Test V2",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "2",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			viperFlags: map[string]interface{}{
				"global.kv_version": "2",
			},
			want: []string{
				"/france/paris/key",
				"/spain/admin",
				"/spain/malaga/random",
				"/uk/london/mi5",
			},
			wantErr: false,
		},
		{
			name: "Smoke Test V1 Buffered",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "1",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
				concurrency:  10,
			},
			want: []string{
				"/france/paris/key",
				"/spain/admin",
				"/spain/malaga/random",
				"/uk/london/mi5",
			},
			wantErr: false,
		},
		{
			name: "Filter by Path",
			args: args{
				c: &testClient{
					engine:      engine,
					kvVersion:   "2",
					engineType:  "kv",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"/spain/malaga/*", "/france/*"},
			},
			want: []string{
				"/france/paris/key",
				"/spain/malaga/random",
			},
			wantErr: false,
		},
		{
			name: "Filter by Path with exclusions",
			args: args{
				c: &testClient{
					engine:      engine,
					kvVersion:   "2",
					engineType:  "kv",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"/spain/malaga/*", "/france/*"},
				excludePaths: []string{"*/random"},
			},
			want: []string{
				"/france/paris/key",
			},
			wantErr: false,
		},
		{
			name: "Path with no secrets",
			args: args{
				c: &testClient{
					engine:      engine,
					kvVersion:   "2",
					engineType:  "kv",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/france/nantes",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Non Existing Path Should Fail",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "2",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/nonexisting",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Non Existing Engine Should Fail",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "2",
					secretPaths: secretPaths,
				},
				engine:       "nonexisting",
				path:         "/",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unsupported Version Should Fail",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "99",
					secretPaths: secretPaths,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "List Error",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "2",
					secretPaths: secretPaths,
					listErr:     fmt.Errorf("error while listing"),
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.viperFlags {
				viper.Set(k, v)
			}
			got, err := kv.RList(tt.args.c, tt.args.engine, tt.args.path, tt.args.includePaths, tt.args.excludePaths,
				tt.args.concurrency)
			viper.Reset()
			if (err != nil) != tt.wantErr {
				t.Errorf("RList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
