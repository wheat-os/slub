/*
Copyright © 2022 NAME HERE <bandl 1658002533@qq.com>

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
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wheat-os/slub/generate"
)

var (
	uid string
)

func IsProject(pwd string) bool {
	proSet := path.Join(pwd, strings.ToLower(path.Base(pwd)))
	proSpider := path.Join(pwd, generate.SpiderDir)
	proConf := path.Join(pwd, "conf.toml")

	checkFile := []string{proSet, proSpider, proConf}

	for _, file := range checkFile {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// genspiderCmd represents the genspider command
var genspiderCmd = &cobra.Command{
	Use:   "genspider",
	Short: "Create a spider instance",
	Long: `To create a crawler instance of slub with this command, 
it is important to note that each spider should have a separate uid 
and need to be given fqdn.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		if !IsProject(pwd) {
			cobra.CheckErr(
				fmt.Errorf("this does not depend on a valid slub project, please recheck, path:%s", pwd),
			)
		}

		// check args
		if ok, _ := regexp.MatchString(`([a-z]|[A-Z])\w+`, args[0]); !ok {
			cobra.CheckErr(`the spider name must be multiple characters 
			(including letters and numbers or underline) that begin with a letter`)
		}

		reg, _ := regexp.Compile("[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+.?")
		if ok := reg.MatchString(args[1]); !ok {
			cobra.CheckErr("please check that the domain name is correct")
		}

		if uid == "" {
			uid = strings.ToLower(args[0])
		}

		cobra.CheckErr(generate.MakeSpider(args[0], args[1], uid))
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(genspiderCmd)

	genspiderCmd.Flags().StringVarP(&uid, "uid", "u", "", `-u=spider_id(Specify an id for the spider, 
noting that it is unique within a project)`)
}
