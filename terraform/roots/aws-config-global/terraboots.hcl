root "aws-config-global" {
  scopeTypes = ["org", "platform", "domain", "environment"]

  scopeMatch {
    scopeTypes = {
      org         = ".*"
      platform    = ".*"
      domain      = ".*"
      environment = ".*"
    }
  }
}
