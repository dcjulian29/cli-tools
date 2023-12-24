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
	"fmt"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump [filename]",
	Short: "Dump sections and values in an INI file",
	Long:  "Dump sections and values in an INI file",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := ini_load(args[0], false)

		cobra.CheckErr(err)

		for _, section := range f.Sections() {
			cmd.Print(fmt.Sprintf("\n[\033[1;31m%s\033[0m]\n", section.Name()))
			for _, k := range section.Keys() {
				cmd.Print(fmt.Sprintf("  \033[1;33m%s\033[0m = \033[1;32m%s\033[0m\n", k.Name(), k.String()))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
