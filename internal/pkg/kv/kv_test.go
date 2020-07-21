package kv_test

import (
	vapi "github.com/hashicorp/vault/api"
)

type testClient struct {
	engine      string
	kvVersion   string
	engineType  string
	secretPaths map[string][]interface{}
	secrets     map[string]map[string]interface{}
	listErr     error
	readErr     error
}

var secretPaths map[string][]interface{}
var secrets map[string]map[string]interface{}
var engine = "secret"

func init() {
	secretPaths = map[string][]interface{}{
		"/":               {"spain/", "france/", "uk/"},
		"/spain/":         {"malaga/", "admin"},
		"/spain/malaga/":  {"random"},
		"/france/":        {"paris/", "nantes/"},
		"/france/paris/":  {"key"},
		"/france/nantes/": {""},
		"/uk/":            {"london/"},
		"/uk/london/":     {"mi5"},
	}

	secrets = map[string]map[string]interface{}{
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
	}
}

func (tc *testClient) ListMounts() (map[string]*vapi.MountOutput, error) {
	return map[string]*vapi.MountOutput{
		tc.engine + "/": {
			Type: tc.engineType,
			Options: map[string]string{
				"version": tc.kvVersion,
			},
		},
	}, nil
}
