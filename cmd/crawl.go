/*
Copyright © 2022 bandl

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
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/wheat-os/slub/generate"
)

var (
	running []string
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Run a spider instance.",
	Long: `Using this command to run an instance of a crawler, 
all spiders are started by default, and if you need to specify running crawlers, 
you should add --run=<spider name>, <spider name2>.`,

	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		if !IsProject(pwd) {
			cobra.CheckErr(
				fmt.Errorf("this does not depend on a valid slub project, please recheck, path:%s", pwd),
			)
		}

		// 创建 bin data
		generate.SetProjectPath(pwd)
		err = generate.SetProjectModByFile(path.Join(pwd, "go.mod"))
		cobra.CheckErr(err)

		for i, v := range running {
			running[i] = fmt.Sprintf("%s.go", v)
		}

		err = generate.MakeRegisterSpider(running)
		cobra.CheckErr(err)

		// run project
		runs := exec.Command("go", "run", ".")
		runs.Dir = pwd
		runs.Stdout = os.Stdout

		cobra.CheckErr(runs.Run())
	},

	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	crawlCmd.Flags().StringSliceVar(&running, "run", nil, "--run=[first, second]")
}
