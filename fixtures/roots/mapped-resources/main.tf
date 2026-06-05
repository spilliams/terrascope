variable "keys" {
  type = list(string)
}
resource "random_string" "this" {
  for_each = toset(var.keys)
  length   = 3
}
