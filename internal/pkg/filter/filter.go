package filter

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"k8s.io/klog/v2"
)

// BuildGlobPattern will build a Glob object from a list of included/excluded Paths
func BuildGlobPattern(paths []string) glob.Glob {
	var globString string
	if len(paths) == 0 {
		globString = ""
	} else if len(paths) == 1 {
		globString = paths[0]
	} else {
		globString = fmt.Sprintf("{%s}", strings.Join(paths, ","))
	}

	return glob.MustCompile(globString)
}

// SecretMatchesGlob evaluates if the given path matches the inclusion/exclusion glob patterns
func SecretMatchesGlob(secretPath string, includeGlobPattern glob.Glob, excludeGlobPattern glob.Glob) bool {
	if !includeGlobPattern.Match(secretPath) {
		klog.V(5).Infof("Discarding secret %s as it doesn't match the inclusion paths %v", secretPath,
			includeGlobPattern)
		return false
	}
	if excludeGlobPattern.Match(secretPath) {
		klog.V(5).Infof("Discarding secret %s as it matches the exclusion paths %v", secretPath,
			excludeGlobPattern)
		return false
	}
	return true
}
