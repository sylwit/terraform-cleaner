variable "name" {
  type = string
  default = "terraform-cleaner"
}

variable "region" {
  type = string
  default = "ca-central-1"
}

variable "instance_ids" {
  type = list(string)
  default = ["i-123", "i-456"]
}

variable "legacy" {}