/*
Copyright © 2024 Kovalev Pavel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Pavel7004/WebShop/pkg/infra/config"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "shop",
	Short: "Backend service for WebShop",
	Long: `Backend service for WebShop. 

Examples:
  $ shop serve`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Read)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}
