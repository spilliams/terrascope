terraboots "AcmeInc" {
  modulesDir = "terraform/modules"
  rootsDir   = "terraform/roots"

  scope {
    name        = "org"
    description = "The organization the resources belong to (e.g. Acme Inc)"
    default     = "Acme Inc"
  }

  scope {
    name        = "platform"
    description = "Which platform the resources belong to"
    validation {
      condition     = contains(["gold", "silver", "bronze"], scope.value)
      error_message = "The platform must be one of gold, silver or bronze."
    }
  }

  scope {
    name        = "domain"
    description = "Which domain the resources belong to"
    validation {
      condition     = length(scope.value) < 15
      error_message = "The domain must be less than 15 characters long"
    }
  }

  scope {
    name        = "environment"
    description = "Which domain environment the resources belong to"
    default     = "dev"
    validation {
      condition     = contains(["dev", "stage", "prod"], scope.value)
      error_message = "The environment must be one of dev, stage or prod."
    }
  }

  scope {
    name        = "region"
    description = "Which region the resources belong to"
    default     = "us-west-2"
    validation {
      condition     = contains(["us-west-1", "us-west-2", "us-east-1", "us-east-2"], scope.value)
      error_message = "The region must be one of us-west-1, us-west-2, us-east-1 or us-east-2."
    }
  }
}
