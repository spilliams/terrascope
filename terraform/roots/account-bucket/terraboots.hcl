root "account-bucket" {
  scopeTypes = ["org", "platform", "domain", "environment"]
}

include {
  path = "../partials/backend.hcl"
}

generate "provider" {
  path     = "provider.tf"
  contents = <<-EOF
  provider "aws" {
    profile = "${attributes.security_portal_profile}"
    region  = "${attributes.security_portal_region}"
  }
  provider "aws" {
    alias = "domain_account"
    profile = "${attributes.security_portal_profile}"
    region = "${attributes.security_portal_region}"
    assume_role {
      role_arn = "arn:aws:iam::${attributes.aws_account_id}:role/TerraformRole"
    }
  }
  EOF
}

inputs = {
  bucket_name = "${lower(scope.domain)}-${lower(scope.environment)}"
}
