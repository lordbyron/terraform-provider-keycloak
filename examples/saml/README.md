Keycloak as IdP with SAML to AWS
================================

This project creates a realm and client in Keycloak, and configures it as a
SAML provider. It also creates the provider and two roles in AWS, and links
them.

## Setup

Keycloak must already have a master realm, as well as a client which terraform
can use. The descriptions of those are in the `master.tf` configuration, but
that cannot be used to create them because that's how terraform will auth!

The file `tfstate.start` contains state for the master realm and `terraform`
client, as described above.

 A note on protocol mappers: there is an automatic pm that gets added to any
 new clients called "role list", which needs to be removed. Terraform doesn't
 make it easy to make sure something _doesn't_ exist, so delete it first.
