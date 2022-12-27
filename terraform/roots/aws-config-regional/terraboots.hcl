root "aws-config-regional" {
  scopeTypes = ["org", "platform", "domain", "environment", "region"]

  dependency {
    root = "aws-config-global"
  }
}
