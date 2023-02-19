variable "keys" {
  type    = list(string)
  default = []
}

resource "random_string" "this" {
  for_each = toset(var.keys)

  length = 16

  upper            = false
  lower            = false
  numeric          = true
  special          = true
  override_special = "abcdef"
}
