package kv

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"rvault/internal/pkg/api"

	vapi "github.com/hashicorp/vault/api"
	"k8s.io/klog/v2"
)

type readResult struct {
	path   string
	err    error
	secret *vapi.Secret
}

func read(c api.VaultClient, engine string, secretPath string, kvVersion string,
	wg *sync.WaitGroup, res chan<- *readResult,
	throttleC chan struct{}) {
	defer wg.Done()
	var err error
	var secret *vapi.Secret
	pathPrefix, err := api.GetReadBasePath(engine, kvVersion)

	if err == nil {
		// If channel is buffered
		if cap(throttleC) > 0 {
			// Block here if channel is full
			throttleC <- struct{}{}
		}

		secret, err = c.Read(path.Join(pathPrefix, secretPath))

		// If channel is buffered
		if cap(throttleC) > 0 {
			// Signal API call done
			<-throttleC
		}
	}

	res <- &readResult{
		path:   secretPath,
		err:    err,
		secret: secret,
	}
}

// RRead reads all secrets for a given path including every subpath. No more than 'concurrency' API queries to Vault
// will be done.
func RRead(c api.VaultClient, engine string, path string, includePaths []string, excludePaths []string,
	concurrency uint32) (map[string]map[string]string, error) {
	var dumpResults []readResult
	var errStrings []string
	secrets := make(map[string]map[string]string)
	kvVersion, err := getKVVersion(c, engine)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	resChan := make(chan *readResult)
	exitChan := make(chan struct{})
	throttleChan := make(chan struct{}, concurrency)
	secretPaths, err := RList(c, engine, path, includePaths, excludePaths, concurrency)
	klog.V(4).Infof("Listing returned %d secret paths", len(secretPaths))
	if err != nil {
		return nil, err
	}
	go func(dumpResults *[]readResult, resChan <-chan *readResult, exitC <-chan struct{}) {
		for {
			select {
			case res := <-resChan:
				*dumpResults = append(*dumpResults, *res)
			case <-exitChan:
				return
			}
		}
	}(&dumpResults, resChan, exitChan)

	wg.Add(len(secretPaths))
	for _, secretPath := range secretPaths {
		go read(c, engine, secretPath, kvVersion, &wg, resChan, throttleChan)
	}

	wg.Wait()
	// finish goroutine ensuring results are processed
	exitChan <- struct{}{}

	for _, dumpResult := range dumpResults {
		if dumpResult.err != nil {
			errStrings = append(errStrings, fmt.Sprintf("Error reading secret '%s': %s", dumpResult.path,
				dumpResult.err.Error()))
			continue
		}
		data, errString := parseSecretData(dumpResult, kvVersion)
		if errString != "" {
			errStrings = append(errStrings, errString)
			continue
		}

		if len(data) > 0 {
			secrets[dumpResult.path] = data
		}
	}
	err = nil
	if len(errStrings) > 0 {
		err = fmt.Errorf("errors found while reading secrets:\n%s", strings.Join(errStrings, "\n"))
	}
	return secrets, err

}
