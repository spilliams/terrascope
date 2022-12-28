root "account-something" {
  scopeTypes = ["org", "platform", "domain", "environment"]
}

include {
  path = "../partials/backend.hcl"
}

generate "provider" {
  path = "provider.tf"
  contents = <<-EOF
  provider "aws" {
    profile = "${}"
    region  = "${}"
  }
  provider "aws" {
    alias = "domain_account"
    profile = "${}"
    region = "${}"
    assume_role {
      role+arn = "arn:aws:iam::${}:role/TerraformRole"
    }
  }
  EOF
}

inputs = {
  
}
