# Opsgy CLI
`opsgy` CLI for interacting with the Opsgy DevOps Suite.
For more info, visit: https://www.opsgy.com.

## Get Started
1. Download the Opsgy binary from the [Release page](https://github.com/opsgy/cli/releases) that corresponds with your system.
2. Place the Opsgy binary somewhere on your `PATH` (eg. `/usr/local/bin/opsgy`)
3. Make the binary executable: `chmod +x /usr/local/bin/opsgy`
4. Login with your Opsgy account: `opsgy login`

Now you're ready to use the `opsgy` CLI. 

## Configure Kubectl to access your Opsgy Kubernetes cluster
To configure `kubectl`, run: `opsgy clusters login <cluster_name> --project=<project_name>`
