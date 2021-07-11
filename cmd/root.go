/*
Copyright © 2021 Tanmay Chaudhry <tanmay.chaudhry@gmail.com>

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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tchaudhry91/tls-check/tlsverify"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tls-check [url]",
	Short: "Check TLS Certificate Information for remote servers",
	Long: `This is a simple command line utility to check the TLS information for a particular remote server.
You can check individual sites or batch operations by supplying a CSV`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please supply a URL to test")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			cert := tlsverify.GetTLSCertDetails(arg)
			if cert.Error != nil {
				fmt.Printf("Error while checking url: %s, Error: %s", arg, cert.Error.Error())
			}
			tlsverify.JSONPrint(cert)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("output", "o", "json", "Output Format")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
