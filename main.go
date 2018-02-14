package main

import (
  "github.com/hashicorp/terraform/plugin"
  "github.com/lordbyron/terraform-provider-keycloak/provider"
)

func main() {
  plugin.Serve(provider.Provider())
}
