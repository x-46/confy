package fileprocessing

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
)

func worker(id int, callback func(path string) error, jobs <-chan string, results chan<- error) {
	for s := range jobs {
		results <- callback(s)
	}
}

func RecursiveFileProcessing(rootDir string, fileExtensions []string, callback func(path string) error) error {
	jobs := make(chan string, 1)
	results := make(chan error, 1)
	defer close(results)

	for i := range 5 { // TODO configurable worker count
		go worker(i, callback, jobs, results)
	}

	jbs := []string{}
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if slices.Contains(fileExtensions, filepath.Ext(path)) {
				jbs = append(jbs, path)
				return nil
			}
		}
		return nil
	})

	numJobs := len(jbs)
	go func() {
		for _, j := range jbs {
			jobs <- j
		}
		close(jobs)
	}()

	errs := []error{err}
	for range numJobs {
		errs = append(errs, <-results)
	}
	<-jobs
	return errors.Join(errs...)
}
