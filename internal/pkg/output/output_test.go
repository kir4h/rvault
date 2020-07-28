package output_test

import (
	"testing"

	"rvault/internal/pkg/output"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var populatedMemFs afero.Fs
var populatedMemFsFileName = "/france/paris/key/id_dsa"
var populatedMemFsFileContent = []byte("existingContent")
var readSecrets = map[string]map[string]string{
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

func init() {
	populatedMemFs = afero.NewMemMapFs()
	_ = afero.WriteFile(populatedMemFs, populatedMemFsFileName, populatedMemFsFileContent, 0700)
}
func TestDump(t *testing.T) {
	type args struct {
		secrets map[string]map[string]string
		fs      afero.Fs
		format  string
	}
	tests := []struct {
		name       string
		args       args
		overwrite  bool
		wantResult string
		wantFiles  map[string]string
		wantErr    bool
	}{
		{
			"Smoke Test File",
			args{
				secrets: readSecrets,
				fs:      afero.NewMemMapFs(),
				format:  "file",
			},
			false,
			"",
			map[string]string{
				"/france/paris/key/id_dsa":    "fewfowefkfkwepfkewkfpweokfeowkfpk",
				"/france/paris/key/id_rsa":    "ewdfpelfr23pwlrp32l4[p23lp2k",
				"/spain/admin/admin.conf":     "dsfdsflfrf43l4tlp",
				"/spain/malaga/random/my.key": "d3ewf2323r21e2",
				"/uk/london/mi5/mi5.conf":     "salt, 324r23432, false",
			},
			false,
		},
		{
			"Skip File Overwrite",
			args{
				secrets: readSecrets,
				fs:      populatedMemFs,
				format:  "file",
			},
			false,
			"",
			map[string]string{
				"/france/paris/key/id_dsa": "existingContent",
			},
			false,
		},
		{
			"File Overwrite",
			args{
				secrets: readSecrets,
				fs:      populatedMemFs,
				format:  "file",
			},
			true,
			"",
			map[string]string{
				"/france/paris/key/id_dsa": "fewfowefkfkwepfkewkfpweokfeowkfpk",
			},
			false,
		},

		{
			"Smoke Test JSON",
			args{
				secrets: readSecrets,
				fs:      afero.NewMemMapFs(),
				format:  "json",
			},
			false,
			`{"/france/paris/key":{"id_dsa":"fewfowefkfkwepfkewkfpweokfeowkfpk","id_rsa":"ewdfpelfr23pwlrp32l4[p23lp2k"},"/spain/admin":{"admin.conf":"dsfdsflfrf43l4tlp"},"/spain/malaga/random":{"my.key":"d3ewf2323r21e2"},"/uk/london/mi5":{"mi5.conf":"salt, 324r23432, false"}}`,
			nil,
			false,
		},
		{
			"Smoke Test YAML",
			args{
				secrets: readSecrets,
				fs:      afero.NewMemMapFs(),
				format:  "yaml",
			},
			false,
			`/france/paris/key:
  id_dsa: fewfowefkfkwepfkewkfpweokfeowkfpk
  id_rsa: ewdfpelfr23pwlrp32l4[p23lp2k
/spain/admin:
  admin.conf: dsfdsflfrf43l4tlp
/spain/malaga/random:
  my.key: d3ewf2323r21e2
/uk/london/mi5:
  mi5.conf: salt, 324r23432, false
`,
			nil,
			false,
		},
		{
			"Unsupported format",
			args{
				secrets: readSecrets,
				fs:      afero.NewMemMapFs(),
				format:  "fakeFormat",
			},
			false,
			"",
			nil,
			true,
		},
		{
			"Fail file exists with folder name",
			args{
				secrets: map[string]map[string]string{
					populatedMemFsFileName: {
						"admin.conf": "dsfdsflfrf43l4tlp",
					},
				},
				fs:     populatedMemFs,
				format: "file",
			},
			false,
			"",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("read.overwrite", tt.overwrite)
			gotResult, err := output.Dump(tt.args.secrets, tt.args.fs, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Dump() \ngotResult = \n%v, \nwant \n%v", gotResult, tt.wantResult)
			}

			if tt.wantFiles != nil {
				for wantFilePath, wantFileContent := range tt.wantFiles {
					exists, _ := afero.Exists(tt.args.fs, wantFilePath)
					if !exists {
						t.Errorf("File '%s' should exist and doesn't", wantFilePath)
					}
					gotFileContent, _ := afero.ReadFile(tt.args.fs, wantFilePath)
					if string(gotFileContent) != wantFileContent {
						t.Errorf("Wrong content for file '%s'.\n Wanted:\n%s \nGot:\n%v", wantFilePath,
							wantFileContent, string(gotFileContent))
					}

				}
			}
		})
	}
}
