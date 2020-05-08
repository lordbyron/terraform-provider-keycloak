Terraform Keycloak Provider
===========================

This project implements a [Terraform provider][] for declaratively configuring
API resources in [Keycloak][].

## Status

This provider can currently manage Keycloak `client` resources, roles, and
protocol mappings. This is enough to set up Keycloak as the SAML provider for
aws (see examples).

Not all fields of those resources are supported at the moment.

There is a single Data Source, `keycloak_client`, which can simply be used to change between client_id (which is name) and guid (which is id).

## Installation

Grab a binary release for your operating system from the [releases][] page and drop it into
`~/.terraform.d/plugins`.

Run `terraform init` to initialise the new provider in the folder containing your configuration
files and `terraform providers` to check that it has been loaded correctly.

**Note**: The targeted version of Terraform is currently **v0.12.0**.

## Building from source

For "vanilla"-builds just do this:

1. Install and configure Go
2. `go get github.com/lordbyron/terraform-provider-keycloak`

## Setup instructions

The Keycloak instance to manage needs to be configured with a client that has
permission to change the resources in Keycloak.

Create a Client in Keycloak with appropriate settings:
```
Client Protocol: openid-connect
Access Type: confidential
Service Accounts Enabled: On
(others disabled)
Scope/ Full Scope Allowed: Off
Scope/ Assigned Roles: admin
Service Account Roles/ Assigned Roles: admin, offline_access, uma_authorization
```

The provider needs to be configured with credentials to access the API (see Credentials tab on the Client):

```
provider "keycloak" {
  # These parameters are required:
  client_id     = "dingus"
  client_secret = "Oox7luexoofeuquaosh5ti3aequie7sh"
  api_base      = "https://keycloak.my-company.acme"
  
  # These parameters are optional:
  realm = "my-company"  # defaults to 'master'
}
```
[Terraform provider]: https://www.terraform.io/docs/plugins/provider.html
[Keycloak]: http://www.keycloak.org/
[configure]: https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin
