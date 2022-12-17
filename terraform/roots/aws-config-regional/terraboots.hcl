root "aws-config-regional" {
  scopes = ["org", "platform", "domain", "environment", "region"]
  dependency {
    root = "aws-config-global"
  }
}
