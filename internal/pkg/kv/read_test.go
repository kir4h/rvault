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
)

func (tc *testClient) Read(searchPath string) (*api.Secret, error) {
	var relPath string
	if tc.readErr != nil {
		return nil, tc.readErr
	}

	if tc.kvVersion == "2" {
		relPath = strings.TrimPrefix(searchPath, path.Join(tc.engine+"/data"))
	} else {
		relPath = strings.TrimPrefix(searchPath, path.Clean(tc.engine))
	}

	if secret, ok := tc.secrets[relPath]; ok {
		if tc.kvVersion == "2" {
			return &api.Secret{
				Data: map[string]interface{}{
					"data": secret,
				},
			}, nil
		} else {
			return &api.Secret{
				Data: secret,
			}, nil
		}
	} else {
		return nil, fmt.Errorf("searchPath %s doesn't exist", searchPath)
	}
}

func TestRRead(t *testing.T) {
	type args struct {
		c            api2.VaultClient
		engine       string
		path         string
		concurrency  uint32
		includePaths []string
		excludePaths []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string]string
		wantErr bool
	}{
		{
			name: "Smoke Test V2",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "2",
					secretPaths: secretPaths,
					secrets:     secrets,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			want: map[string]map[string]string{
				"/spain/admin": {
					"admin.conf": "dsfdsflfrf43l4tlp",
				},
				"/spain/malaga/random": {
					"my.key": "d3ewf2323r21e2",
				},
				"/france/paris/key": {
					"id_rsa": "ewdfpelfr23pwlrp32l4[p23lp2k",
					"id_dsa": "fewfowefkfkwepfkewkfpweokfeowkfpk",
				},
				"/uk/london/mi5": {
					"mi5.conf": "salt, 324r23432, false",
				},
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
					secrets:     secrets,
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
				concurrency:  20,
			},
			want: map[string]map[string]string{
				"/spain/admin": {
					"admin.conf": "dsfdsflfrf43l4tlp",
				},
				"/spain/malaga/random": {
					"my.key": "d3ewf2323r21e2",
				},
				"/france/paris/key": {
					"id_rsa": "ewdfpelfr23pwlrp32l4[p23lp2k",
					"id_dsa": "fewfowefkfkwepfkewkfpweokfeowkfpk",
				},
				"/uk/london/mi5": {
					"mi5.conf": "salt, 324r23432, false",
				},
			},
			wantErr: false,
		},
		{
			name: "Read error",
			args: args{
				c: &testClient{
					engine:      engine,
					engineType:  "kv",
					kvVersion:   "2",
					secretPaths: secretPaths,
					secrets:     secrets,
					readErr:     fmt.Errorf("error reading"),
				},
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			wantErr: true,
			want:    map[string]map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kv.RRead(tt.args.c, tt.args.engine, tt.args.path, tt.args.includePaths, tt.args.excludePaths,
				tt.args.concurrency)
			if (err != nil) != tt.wantErr {
				t.Errorf("RRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RRead() got = %v, want %v", got, tt.want)
			}
		})
	}
}
