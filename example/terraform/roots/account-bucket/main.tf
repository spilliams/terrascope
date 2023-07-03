resource "random_string" "slug" {
  length  = 6
  special = false
  upper   = false
  lower   = true
  numeric = true
}

variable "bucket_name" {
  description = "The name to give to the account bucket"
  type        = string
}

resource "aws_s3_bucket" "account" {
  bucket = "${var.bucket_name}-${random_string.slug.result}"
}
