package kv

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
)

func Test_parseSecretData(t *testing.T) {
	type args struct {
		result    readResult
		kvVersion string
	}
	tests := []struct {
		name           string
		args           args
		wantSecretData map[string]string
		wantErr        bool
	}{
		{
			name: "Smoke Test",
			args: args{
				result: readResult{
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
				kvVersion: "2",
			},
			wantSecretData: map[string]string{"foo": "bar"},
			wantErr:        false,
		},
		{
			name: "Secret with nil data is ignored",
			args: args{
				result: readResult{
					path: "/my/path",
					err:  nil,
					secret: &api.Secret{
						Data: map[string]interface{}{
							"data": nil,
						},
					},
				},
				kvVersion: "2",
			},
			wantSecretData: map[string]string{},
			wantErr:        false,
		},
		{
			name: "Secret with unexpected data type should fail",
			args: args{
				result: readResult{
					path: "/my/path",
					err:  nil,
					secret: &api.Secret{
						Data: map[string]interface{}{
							"data": map[string]string{
								"foo": "bar",
							},
						},
					},
				},
				kvVersion: "2",
			},
			wantSecretData: map[string]string{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSecretData, gotErrString := parseSecretData(tt.args.result, tt.args.kvVersion)
			if !reflect.DeepEqual(gotSecretData, tt.wantSecretData) {
				t.Errorf("parseSecretData() gotSecretData = %v, want %v", gotSecretData, tt.wantSecretData)
			}
			gotErr := len(gotErrString) > 0
			if gotErr != tt.wantErr {
				t.Errorf("parseSecretData() gotErrString = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
