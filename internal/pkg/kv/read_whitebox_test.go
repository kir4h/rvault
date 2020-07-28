package kv

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
)

var wantSmokeTest map[string]map[string]string

func init() {
	wantSmokeTest = map[string]map[string]string{
		"/my/path": {
			"foo": "bar",
		},
	}
}
func Test_parseReadResults(t *testing.T) {
	type args struct {
		dumpResults []readResult
		kvVersion   string
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
				[]readResult{
					{
						path: "/my/path",
						err:  nil,
						secret: &api.Secret{
							Data: map[string]interface{}{
								"data": map[string]interface{}{
									"foo": "bar",
								},
							},
						},
					},
				},
				"2",
			},
			want:    wantSmokeTest,
			wantErr: false,
		},
		{
			name: "Smoke Test V1",
			args: args{
				[]readResult{
					{
						path: "/my/path",
						err:  nil,
						secret: &api.Secret{
							Data: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				"1",
			},
			want:    wantSmokeTest,
			wantErr: false,
		},
		{
			name: "Nil secret should be ignored",
			args: args{
				[]readResult{
					{"/my/emptySecret",
						nil,
						nil,
					},
				},
				"2",
			},
			want:    map[string]map[string]string{},
			wantErr: false,
		},
		{
			name: "Error reading from Vault",
			args: args{
				[]readResult{
					{"/my/fakePath",
						fmt.Errorf("error reading from Vault"),
						nil,
					},
				},
				"2",
			},
			want:    map[string]map[string]string{},
			wantErr: true,
		},
		{
			name: "Unexpected Secret Content",
			args: args{
				[]readResult{
					{"/my/fakePath",
						nil,
						&api.Secret{
							Data: map[string]interface{}{
								"NotData": map[string]interface{}{
									"foo": "bar",
								},
							},
						},
					},
				},
				"2",
			},
			want:    map[string]map[string]string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReadResults(tt.args.dumpResults, tt.args.kvVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReadResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseReadResults() got = %v, want %v", got, tt.want)
			}
		})
	}
}
