/*
Copyright © 2022  <>

Licensed under the HLT License, Version 0.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/shammishailaj/taskfactory/pkg/schemas"
	"github.com/spf13/cobra"
)

// cleanCmd represents the cleanCmd command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version information",
	Long:  `Displays version information for the App`,
	Run: func(cmd *cobra.Command, args []string) {

		semVer := &schemas.SemanticVersion{
			GitBranch:  GitBranch,
			GitState:   GitState,
			GitSummary: GitSummary,
			BuildDate:  BuildDate,
			Version:    Version,
			GitCommit:  GitCommit,
		}
		fmt.Println(semVer.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
