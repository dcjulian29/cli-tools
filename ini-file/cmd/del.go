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
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:     "delete [filename]",
	Aliases: []string{"del"},
	Short:   "Delete values in an INI file",
	Long:    "Delete values in an INI file",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := ini_load(args[0], false)

		cobra.CheckErr(err)

		sectionName, _ := cmd.Flags().GetString("section")
		section := f.Section(sectionName)
		key, _ := cmd.Flags().GetString("key")

		if section.HasKey(key) {
			section.DeleteKey(key)
			cobra.CheckErr(ini_save(args[0], f))
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	delCmd.Flags().StringP("section", "s", "", "section name")
	delCmd.Flags().StringP("key", "k", "", "key name")

	delCmd.MarkFlagRequired("section")
	delCmd.MarkFlagRequired("key")
}
