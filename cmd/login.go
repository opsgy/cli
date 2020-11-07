/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/opsgy/cli/opsgy"
	"github.com/spf13/cobra"
	"log"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with your Opsgy Account",
	Long:  `Login with your Opsgy Account`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := opsgy.Login(9932, []string{"all"})
		if err != nil {
			log.Fatal(err)
		}

		config := opsgy.LoadConfig()
		config.AccessToken = credentials.AccessToken
		config.RefreshToken = credentials.RefreshToken
		config.TokenExpiry = credentials.TokenExpiry
		config.TokenType = credentials.TokenType

		err = opsgy.SaveConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch default project
		client, err := opsgy.GetClient(config)
		if err != nil {
			log.Fatal(err)
		}
		projects, err := opsgy.GetProjects(client);
		if err != nil {
			log.Fatal(err)
		}
		if len(projects) > 0 {
			config.ProjectName = projects[0].Name
			err = opsgy.SaveConfig(config)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// loginCmd.Flags().Int("port", 9932, "List on this port for the OAuth callback request")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
