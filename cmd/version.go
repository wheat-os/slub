/*
Copyright © 2022 bandl HERE <1658002533@qq.com>

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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wheat-os/slub/generate"
)

var (
	checkSlubby bool
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Gets the version number of the framework",
	Long: `Use this command to query slub, as well as the 
latest or historical version number of slubby.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("slub: %s\n", generate.Version)
		if !checkSlubby {
			return
		}

		tags, err := generate.GetSlubbyVersion()
		if err != nil {
			fmt.Println(err)
		}

		if len(tags) == 0 {
			cobra.CheckErr("didn't get a proper slubby version")
			return
		}
		fmt.Printf("slubby latest version: %s\n", tags[0].Name)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&checkSlubby, "slubby", "s", false, "-s=true (print slubby latest version)")
}
