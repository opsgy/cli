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
	"os/exec"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/opsgy/cli/opsgy"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var kubeFile string

// clustersLoginCmd represents the clustersLogin command
var clustersLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Configure credentials to use kubectl",
	Long:  `Configure credentials to use kubectl`,
	Run: func(cmd *cobra.Command, args []string) {
		if kubeFile == "" {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			kubeFile = home + "/.kube/config"
		}

		// Load opsgy config
		config := opsgy.LoadConfig()
		if config.AccessToken == "" {
			// Not authenticated
			fmt.Println("You're not autenticated, please run first: opsgy login")
			os.Exit(1)
		}
		token, err := jwt.Parse(config.AccessToken, nil)
		if token == nil {
			fmt.Println(err)
			os.Exit(1)
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		// Find cluster
		client, err := opsgy.GetClient(config)
		if err != nil {
			log.Fatal(err)
		}

		project, err := opsgy.GetProjectByName(client, ProjectName)
		if err != nil {
			log.Fatal(err)
		}

		cluster, err := opsgy.GetClusterByName(client, project.ProjectID, args[0])
		if err != nil {
			log.Fatal(err)
		}
		if cluster.Status == nil || cluster.Status.KubeApiURL == "" {
			fmt.Println("Cluster is not ready yet")
			os.Exit(1)
		}

		// Create token with only clusters scope
		credentials, err := opsgy.Login(9932, []string{"kubernetes", "email", "profile", "openid"})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Configuring kubectl...")

		// Set user
		username := "opsgy_" + strings.ReplaceAll(strings.ReplaceAll(claims["email"].(string), ".", "_"), "@", "_")
		_, err = exec.Command("kubectl", "config",
			"--kubeconfig", kubeFile,
			"set-credentials", username,
			"--auth-provider", "oidc",
			"--auth-provider-arg", "idp-issuer-url="+opsgy.OpsgyAccountUrl,
			"--auth-provider-arg", "client-id="+opsgy.ClientID,
			"--auth-provider-arg", "client-secret="+opsgy.ClientSecret,
			"--auth-provider-arg", "id-token="+credentials.AccessToken,
			"--auth-provider-arg", "refresh-token="+credentials.RefreshToken,
		).Output()

		// Set cluster
		clusterName := "opsgy_" + cluster.Name
		_, err = exec.Command("kubectl", "config",
			"--kubeconfig", kubeFile,
			"set-cluster", clusterName,
			"--server", cluster.Status.KubeApiURL,
		).Output()

		// Set context
		contextName := "opsgy_" + cluster.Name
		_, err = exec.Command("kubectl", "config",
			"--kubeconfig", kubeFile,
			"set-context", contextName,
			"--cluster", clusterName,
			"--user", username,
		).Output()

		// use context
		_, err = exec.Command("kubectl", "config",
			"--kubeconfig", kubeFile,
			"use-context", contextName,
		).Output()

		fmt.Println("All set! You can now use 'kubectl' to access your cluster")
	},
}

func init() {
	clustersCmd.AddCommand(clustersLoginCmd)

	clustersLoginCmd.Flags().StringVar(&kubeFile, "kube-config", "", "Kubernetes config file (default is $HOME/.kube/config)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustersLoginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustersLoginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
