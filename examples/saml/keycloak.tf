provider "keycloak" {
  client_id     = "terraform"
  # secret is the same as what would come from
  # ${keycloak_client.terraform_client.client_secret}, but there is a chicken/egg
  # problem using that!
  client_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  api_base      = "https://keycloak.internal.example.com"
}
