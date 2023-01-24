terraboots "AcmeInc" {
  rootsDir  = "terraform/roots"
  scopeData = ["data.hcl"]

  scope {
    name        = "org"
    description = "The organization the resources belong to (e.g. Acme Inc)"
  }

  scope {
    name        = "platform"
    description = "Which platform the resources belong to"
  }

  scope {
    name        = "domain"
    description = "Which domain the resources belong to"
  }

  scope {
    name        = "environment"
    description = "Which domain environment the resources belong to"
    default     = "dev"
  }

  scope {
    name        = "region"
    description = "Which region the resources belong to"
    default     = "us-west-2"
  }
}
