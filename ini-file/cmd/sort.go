/*
Copyright © 2023 Julian Easterling <julian@julianscorner.com>

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
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var sortCmd = &cobra.Command{
	Use:   "sort [filename]",
	Short: "Sort sections and values in an INI file",
	Long:  "Sort sections and values in an INI file",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := ini_load(args[0], false)

		cobra.CheckErr(err)

		i, err := ini_load("", true)

		cobra.CheckErr(err)

		i.DeleteSection(ini.DefaultSection)

		sectionNames := f.SectionStrings()
		sort.Strings(sectionNames)

		for _, s := range sectionNames {
			section, err := f.GetSection(s)

			cobra.CheckErr(err)

			keyNames := section.KeyStrings()
			sort.Strings(keyNames)

			if s == ini.DefaultSection && len(keyNames) == 0 {
				continue
			}

			n, err := i.NewSection(s)

			cobra.CheckErr(err)

			for _, k := range keyNames {
				key, _ := section.GetKey(k)

				n.NewKey(key.Name(), key.Value())
			}
		}

		backup, _ := cmd.Flags().GetBool("backup")

		if backup {
			c := 0
			filename := fmt.Sprintf("%s.bak", args[0])

			for {
				if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
					input, err := os.ReadFile(args[0])

					cobra.CheckErr(err)

					err = os.WriteFile(filename, input, 0644)

					cobra.CheckErr(err)
					break
				} else {
					c++
					filename = fmt.Sprintf("%s.bak.%v", args[0], c)
				}
			}
		}

		cobra.CheckErr(i.SaveTo(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().Bool("backup", true, "Make a backup of the unsorted file")
}
