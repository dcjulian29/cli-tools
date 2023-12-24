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

var getCmd = &cobra.Command{
	Use:   "get [filename]",
	Short: "Get values in an INI file",
	Long:  "Get values in an INI file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := ini_load(args[0], false)

		cobra.CheckErr(err)

		sectionName, _ := cmd.Flags().GetString("section")
		section := f.Section(sectionName)
		key, _ := cmd.Flags().GetString("key")

		if len(key) == 0 {
			for _, k := range section.Keys() {
				cmd.Print(fmt.Sprintf("%s=%s\n", k.Name(), k.String()))
			}
		} else {
			k, err := section.GetKey(key)

			cobra.CheckErr(err)

			cmd.Print(k.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("section", "s", "", "section name")
	getCmd.Flags().StringP("key", "k", "", "key name")

	getCmd.MarkFlagRequired("section")
}
