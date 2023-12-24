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
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var rootCmd = &cobra.Command{
	Use:   "ini-file",
	Short: "A tool to manage INI files",
	Long:  "A tool to manage INI files",
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func ini_load(filename string, createIfMissing bool) (*ini.File, error) {
	f, err := ini.LoadSources(ini_options(), filename)

	if err == nil {
		return f, nil
	}

	if createIfMissing && os.IsNotExist(err) {
		i := ini.Empty()
		i.DeleteSection(ini.DefaultSection)

		return i, nil
	}

	return nil, err
}

func ini_options() ini.LoadOptions {
	return ini.LoadOptions{
		AllowBooleanKeys:          true,
		AllowShadows:              true,
		Loose:                     false,
		IgnoreInlineComment:       true,
		UnescapeValueDoubleQuotes: true,
	}
}

func ini_save(filename string, iniFile *ini.File) error {
	f, err := os.OpenFile(filename, os.O_SYNC|os.O_RDWR|os.O_CREATE, os.FileMode(0644))

	if err != nil {
		return err
	}

	err = f.Truncate(0)

	if err != nil {
		return err
	}

	_, err = iniFile.WriteTo(f)

	if err != nil {
		return err
	}

	return f.Close()
}
