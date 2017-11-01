provider "manifold" {
  // uses the default and loads the API key from the ENV `MANIFOLD_API_TOKEN`.
}

resource "manifold_api_token" "manifold" {
  team        = "manifold"
  role        = "read"
  description = "New token"
}

output "token" {
  value = "${manifold_api_token.manifold.token}"
}
