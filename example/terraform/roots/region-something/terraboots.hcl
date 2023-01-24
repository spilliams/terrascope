root "region-something" {
  scopeTypes = ["org", "platform", "domain", "environment", "region"]

  dependency {
    root = "account-something"
  }
}
