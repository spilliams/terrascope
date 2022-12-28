terraform {
  backend "s3" {
    profile = attributes.security_portal_profile
    region  = attributes.security_portal_region

    bucket         = "terraform-state"
    dynamodb_table = "terraform-state-lock"
    key            = "${join("/", scope.values)}/${root.id}/terraform.tfstate"
    encrypt        = true
  }
}
