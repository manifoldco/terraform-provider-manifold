data "manifold_credential" "my-credential" {
  project  = "terraform"
  resource = "custom-resource1-1"
  key      = "MY_CREDENTIAL_KEY"
  default  = "my-value"
}
