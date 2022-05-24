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
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wheat-os/slub/generate"
)

var (
	pPath string
	pMod  string
)

// startprojectCmd represents the startproject command
var startprojectCmd = &cobra.Command{
	Use:   "startproject",
	Short: "Create a slub project",
	Long: `This method is used to create a spider development
scaffold with slubby default configuration items.`,

	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		if ok, _ := regexp.MatchString(`^([a-z]|[A-Z])\w+$`, projectName); !ok {
			cobra.CheckErr("the project name must start with an English letter and be at least 2 characters long")
		}

		// 设置 mod
		if pMod == "" {
			pMod = strings.ToLower(projectName)
		}
		generate.SetProjectMod(pMod)

		proPath, err := os.Getwd()
		cobra.CheckErr(err)

		if pPath != "" {
			proPath = pPath
		}

		// 创建文件夹
		err = generate.MakeDirProjectDir(proPath, projectName)
		cobra.CheckErr(err)

		generate.SetProjectPath(path.Join(proPath, projectName))

		err = generate.MakeDefaultSettings()
		cobra.CheckErr(err)

	},

	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(startprojectCmd)

	startprojectCmd.Flags().StringVar(&pPath, "path", "", "--path=/home(specifies the build directory)")

	startprojectCmd.Flags().StringVarP(&pMod, "pkg-name", "m", "", "-m=firstproject(specifies the mod name)")
}
