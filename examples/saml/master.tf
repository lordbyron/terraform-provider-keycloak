#### MASTER REALM
# This realm already exists, and must be configured appropriately in order for
# terraform to do anything. Let this section be a reference for how to
# configure the master realm and clients.
# IMPORTANT: Again, don't use this to create the realm. Manually create the
# realm and terraform client (see below) and make sure that `terraform plan`
# shows no changes for them.
####
resource "keycloak_realm" "master_realm" {
  realm = "master"
  display_name = "Keycloak"
  default_roles = ["offline_access", "uma_authorization"]
  remember_me = false
  verify_email = false
  reset_password_allowed = false
  edit_username_allowed = false
  brute_force_protected = false

}

resource "keycloak_client" "terraform_client" {
  realm = "master"
  client_id = "terraform"
  redirect_uris = []
  protocol = "openid-connect"
  public_client = false
  bearer_only = false
  service_accounts_enabled = true
  full_scope_allowed = false

  attributes = {
    saml.assertion.signature = "false"
    saml.authnstatement = "false"
    saml.client.signature = "false"
    saml.encrypt = "false"
    saml.force.post.binding = "false"
    saml.multivalued.roles = "false"
    saml.onetimeuse.condition = "false"
    saml.server.signature = "false"
    saml.server.signature.keyinfo.ext = "false"
    saml_force_name_id_format = "false"
  }
}
