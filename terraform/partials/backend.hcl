terraform {
  backend "s3" {
    profile = ""
    region  = ""

    bucket         = ""
    dynamodb_table = ""
    key            = ""
    encrypt        = true
  }
}
