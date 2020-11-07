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
	"fmt"
	"log"
	"os"

	"github.com/opsgy/cli/opsgy"

	"text/tabwriter"

	"github.com/spf13/cobra"
)

// projectsListCmd represents the projectsList command
var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects`,
	Run: func(cmd *cobra.Command, args []string) {
		config := opsgy.LoadConfig()
		client, err := opsgy.GetClient(config)
		if err != nil {
			log.Fatal(err)
		}

		projects, err := opsgy.GetProjects(client)
		if err != nil {
			log.Fatal(err)
		}

		// initialize tabwriter
		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		defer w.Flush()

		fmt.Fprintf(w, "%s\t%s\t%s\n", "NAME", "ID", "DESCRIPTION")

		for _, project := range projects {
			fmt.Fprintf(w, "%s\t%s\t%s\n", project.Name, project.ProjectID, project.Description)
		}

	},
}

func init() {
	projectsCmd.AddCommand(projectsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectsListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectsListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
