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

var setCmd = &cobra.Command{
	Use:   "set [filename]",
	Short: "Set values in an INI file",
	Long:  "Set values in an INI file",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := ini_load(args[0], true)

		cobra.CheckErr(err)

		sectionName, _ := cmd.Flags().GetString("section")
		section := f.Section(sectionName)
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")

		if section.HasKey(key) && len(section.Key(key).ValueWithShadows()) == 1 {
			section.Key(key).SetValue(value)
		} else {
			section.NewKey(key, value)
		}

		cobra.CheckErr(ini_save(args[0], f))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("section", "s", "", "section name")
	setCmd.Flags().StringP("key", "k", "", "key name")
	setCmd.Flags().StringP("value", "v", "", "key value")

	setCmd.MarkFlagRequired("section")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("value")
}
