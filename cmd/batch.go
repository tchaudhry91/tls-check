/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/tchaudhry91/tls-check/tlsverify"
)

var outputFile string
var outputFormat string
var concurrency int

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch [input_file]",
	Short: "Run the check across a list of urls",
	Long:  `Supply an input file that contains newline separate urls that will all be batched. This produces a CSV/JSON with all the details.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please supply an input file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Read File
		urls, err := readURLsFromFile(args[0])
		if err != nil {
			fmt.Printf("Unable to read URLs from %s, Error: %s", args[0], err.Error())
			return
		}
		// Fetch the urls in batch
		certs := concurrentVerify(concurrency, urls)

		// Write out the certs as a JSON file
		ofile, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Unable to create output file: %s", err.Error())
		}
		defer ofile.Close()

		err = json.NewEncoder(ofile).Encode(certs)
		if err != nil {
			fmt.Printf("Unable to write to output file: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&outputFile, "output-file", "f", "tls-check.json", "Select a different file to output to.")
	batchCmd.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output Format: json")
	batchCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 100, "Max concurrent connections")
}

func readURLsFromFile(fname string) ([]string, error) {
	urls := []string{}
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return urls, err
	}
	urlstr := string(data)
	urls = strings.Split(urlstr, "\n")
	return urls, nil
}

func concurrentVerify(goroutines int, urls []string) (certs []tlsverify.CertDetails) {
	certs = []tlsverify.CertDetails{}
	chunkStart := 0
	chunkEnd := chunkStart + goroutines
	if chunkEnd >= len(urls) {
		chunkEnd = len(urls)
	}
	for chunkStart != chunkEnd {
		var wg sync.WaitGroup
		resultsChan := make(chan tlsverify.CertDetails, goroutines)
		for i := chunkStart; i < chunkEnd; i++ {
			wg.Add(1)
			go func(rawURL string, resultsChan chan<- tlsverify.CertDetails) {
				defer wg.Done()
				cert := tlsverify.GetTLSCertDetails(rawURL)
				resultsChan <- cert
			}(urls[i], resultsChan)
		}
		wg.Wait()
		close(resultsChan)
		chunkStart = chunkEnd
		chunkEnd = chunkStart + goroutines
		if chunkEnd >= len(urls) {
			chunkEnd = len(urls)
		}
		for c := range resultsChan {
			certs = append(certs, c)
		}
	}
	return certs
}
