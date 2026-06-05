variable "qty" {
  type = number
}
resource "random_string" "this" {
  count  = var.qty
  length = 3
}
