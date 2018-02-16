#### Employee Realm
# The realm that employees use.
####

resource "keycloak_realm" "employee_realm" {
  realm = "EmployeeRealm"
  display_name = "TheRealm"
  default_roles = ["offline_access", "uma_authorization"]
  remember_me = false
  verify_email = false
  reset_password_allowed = false
  edit_username_allowed = false
}

resource "keycloak_client" "aws_saml" {
  realm = "${keycloak_realm.employee_realm.realm}"
  client_id = "urn:amazon:webservices"
  protocol = "saml"
  public_client = false
  redirect_uris = ["https://signin.aws.amazon.com/saml"]
  web_origins = ["https://signin.aws.amazon.com"]
  base_url = "/auth/realms/${keycloak_realm.employee_realm.realm}/protocol/saml/clients/amazon-aws"
  full_scope_allowed = false

  # This is the same as installing the saml metadata xml from
  # https://signin.aws.amazon.com/static/saml-metadata.xml
  attributes = {
    saml.assertion.signature = "true"
    saml.force.post.binding = "true"
    saml.multivalued.roles = "false"
    saml.encrypt = "false"
    saml_assertion_consumer_url_post = "https://signin.aws.amazon.com/saml"
    saml.server.signature = "true"
    saml_idp_initiated_sso_url_name = "amazon-aws"
    saml.server.signature.keyinfo.ext = "false"
    saml.signing.certificate = "MIIDbTCCAlWgAwIBAgIEdcdzXTANBgkqhkiG9w0BAQsFADBnMR8wHQYDVQQDExZ1cm46YW1hem9uOndlYnNlcnZpY2VzMSIwIAYDVQQKExlBbWF6b24gV2ViIFNlcnZpY2VzLCBJbmMuMRMwEQYDVQQIEwpXYXNoaW5ndG9uMQswCQYDVQQGEwJVUzAeFw0xODAxMTkwMDAwMDBaFw0xOTAxMTkwMDAwMDBaMGcxHzAdBgNVBAMTFnVybjphbWF6b246d2Vic2VydmljZXMxIjAgBgNVBAoTGUFtYXpvbiBXZWIgU2VydmljZXMsIEluYy4xEzARBgNVBAgTCldhc2hpbmd0b24xCzAJBgNVBAYTAlVTMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvFEhKRreEuXHqUEu2eFSDpZEnEGTW8eXfLMWEuSMh4s5b/VJ6tIXN8W/gdVPOzi4trNRqZ/3gqCQhWR0AAA+QjlHb/PdMt9hXzgCkm2MFq4Zsx0w1csudKBMQUA6kK1sNFXSvo86CDlGFEJYpM6NmHwd699lBdYSuTm9J8R8qSjFJe5d8gU71qTUB2g1GVjEcEZRboSF9BZdrV7wm+ytw4NtDxRO/hFKIeAYy8BuI5JdO65NZ8cFLL8i4tEh1tFd561NMhb0S8BrRRncw7XoQL9N0ug2j417Jzkg9i8dbHMU7FcAgfScTcm+HvbLswTi2Ml9xkVsoHbS9KPqjD0ZEQIDAQABoyEwHzAdBgNVHQ4EFgQUpa9rLa3W+cM+74SJ1JbSlfGSbbIwDQYJKoZIhvcNAQELBQADggEBAHdqSjkGlxKfB7+Sp/VPhVFE2X5RNHt7LFxrpAhSJdCbPUDlGvNGrKQWi2da+lM63+fRRwO8m3/AJA8KAXwORddnGQZi9YVtL2roDV499yVP6Nfctxo8Hu5BmOLlU7575CAP09iApMJiN3UmzXkEdixqTJYQUvMWsiO5ObxISTamrC+Pey014L5gdYbEcFIjZC0oxZEYuB4bcIZ9DSYdDEYN+bVqwQIWcbxYsUpayEXxMbE42J5FOxjWp+2jE+19czwuUapHstkqo1TZSd4iQluKKknPCo7P34MQkIPcIa3Q/AEibRN7OnS2RH0ZRulaAwyhJpIvIuccRRCiz0uRY6s="
    saml.signature.algorithm = "RSA_SHA256"
    saml_force_name_id_format = "false"
    saml.client.signature = "true"
    saml.authnstatement = "true"
    saml_name_id_format = "transient"
    saml_signature_canonicalization_method = "http://www.w3.org/2001/10/xml-exc-c14n#"
    saml.onetimeuse.condition = "false"
  }
}

resource "keycloak_protocol_mapper" "employee_pm_session_name" {
  realm = "${keycloak_realm.employee_realm.realm}"
  client_id = "${keycloak_client.aws_saml.id}"

  name = "Session Name"
  protocol = "saml"
  protocol_mapper = "saml-user-property-mapper"
  consent_required = "false"
  config = {
    single = "",
    attribute.nameformat = "Basic",
    user.attribute = "username",
    friendly.name = "Session Name",
    attribute.name = "https://aws.amazon.com/SAML/Attributes/RoleSessionName"
  }
}

resource "keycloak_protocol_mapper" "employee_pm_session_role" {
  realm = "${keycloak_realm.employee_realm.realm}"
  client_id = "${keycloak_client.aws_saml.id}"

  name = "Session Role"
  protocol = "saml"
  protocol_mapper = "saml-role-list-mapper"
  consent_required = "false"
  config = {
      single = "",
      attribute.nameformat = "Basic",
      role = "",
      new.role.name = "",
      friendly.name = "Session Role",
      attribute.name = "https://aws.amazon.com/SAML/Attributes/Role"
  }
}

resource "keycloak_protocol_mapper" "employee_pm_session_duration" {
  realm = "${keycloak_realm.employee_realm.realm}"
  client_id = "${keycloak_client.aws_saml.id}"

  name = "Session Duration"
  protocol = "saml"
  protocol_mapper = "saml-hardcode-attribute-mapper"
  consent_required = "false"
  config = {
      attribute.value = "28800",
      attribute.nameformat = "Basic",
      friendly.name = "Session Duration",
      attribute.name = "https://aws.amazon.com/SAML/Attributes/SessionDuration"
  }
}
