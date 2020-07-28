package kv

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"

	"rvault/internal/pkg/api"
	"rvault/internal/pkg/filter"

	"github.com/gobwas/glob"
	vapi "github.com/hashicorp/vault/api"
	"k8s.io/klog/v2"
)

func lookup(c *vapi.Client, engine string, searchPath string, includeGlobPattern glob.Glob,
	excludeGlobPattern glob.Glob, kvVersion string,
	wg *sync.WaitGroup, secC chan<- string,
	errC chan<- error, throttleC chan struct{}) {
	defer wg.Done()
	pathPrefix, err := api.GetListBasePath(engine, kvVersion)
	if err != nil {
		errC <- err
		return
	}

	// If channel is buffered
	if cap(throttleC) > 0 {
		// Block here if channel is full
		throttleC <- struct{}{}
	}

	klog.V(5).Infof("Listing for %s", path.Join(pathPrefix, searchPath)+"/")
	secret, err := c.Logical().List(path.Join(pathPrefix, searchPath) + "/")

	// If channel is buffered
	if cap(throttleC) > 0 {
		// Signal API call done
		<-throttleC
	}

	if err != nil {
		errC <- err
		return
	} else if secret == nil {
		klog.Infof("No secrets found for path %s", searchPath)
		errC <- nil
		return
	}
	items := secret.Data["keys"].([]interface{})
	for _, item := range items {
		itemPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(searchPath, "/"), item)
		if strings.HasSuffix(item.(string), "/") {
			wg.Add(1)
			go lookup(c, engine, itemPath, includeGlobPattern, excludeGlobPattern, kvVersion, wg, secC, errC,
				throttleC)
		} else if item != "" && filter.SecretMatchesGlob(itemPath, includeGlobPattern, excludeGlobPattern) {
			secC <- itemPath
		}
	}

	errC <- nil
}

// RList lists all secrets for a given 'path' including every subpath as long as they match one of the 'includePaths'.
// No more than 'concurrency' API queries to Vault will be done.
func RList(c *vapi.Client, engine string, path string, includePaths []string, excludePaths []string,
	concurrency uint32) ([]string, error) {
	var secretPaths []string
	var errors []error
	kvVersion, err := getKVVersion(c, engine)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error)
	secretChan := make(chan string)
	exitChan := make(chan struct{})
	throttleChan := make(chan struct{}, concurrency)
	go func(secretPaths *[]string, errorList *[]error, secC <-chan string, errC <-chan error, exitC <-chan struct{}) {
		for {
			select {
			case err := <-errC:
				if err != nil {
					*errorList = append(*errorList, err)
				}
			case sec := <-secretChan:
				*secretPaths = append(*secretPaths, sec)
			case <-exitC:
				return
			}
		}
	}(&secretPaths, &errors, secretChan, errChan, exitChan)

	wg.Add(1)
	includeGlobPattern := filter.BuildGlobPattern(includePaths)
	excludeGlobPattern := filter.BuildGlobPattern(excludePaths)
	go lookup(c, engine, path, includeGlobPattern, excludeGlobPattern, kvVersion, &wg, secretChan, errChan,
		throttleChan)

	wg.Wait()

	// finish goroutine ensuring all results are processed
	exitChan <- struct{}{}

	if len(errors) > 0 {
		var errStrings []string
		for _, errItem := range errors {
			errStrings = append(errStrings, errItem.Error())
		}
		return secretPaths, fmt.Errorf("there were %d error(s) while listing secrets: %s", len(errors),
			strings.Join(errStrings, "\n"))
	}
	sort.Strings(secretPaths)
	return secretPaths, nil

}
