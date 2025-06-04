variable "ami_id" {
  default = "ami-0ce0152a3f6225d58"
  type    = string
}

variable "instance_type" {
  default = "t3.medium"
  type    = string
}

variable "key_name" {
  type = string
}

variable "root_block_size" {
  default = 30
  type    = number
}

## App Config
variable "create_dns_record" {
  type    = bool
  default = false
}
variable "domain" {
  type = string
}
variable "dev_mode" {
  type    = bool
  default = true
}

variable "serve_day" {
  type    = string
  default = "07-12-25"
}

variable "api_port" {
  type    = number
  default = 8080
}

variable "db_port" {
  type    = number
  default = 5432
}

variable "db_user" {
  type    = string
  default = "postgres"
}

variable "db_pass" {
  type    = string
  default = "postgres"
}

variable "db_name" {
  type    = string
  default = "serve"
}

variable "db_instance" {
  type    = string
  default = "db.t4g.micro"
}

variable "auth0_domain" {
  type = string
}

variable "auth0_audience" {
  type = string
}

variable "auth0_client_id" {
  type = string
}

variable "auth0_client_secret" {
  type = string
}

variable "mailtrap_host" {
  type    = string
  default = "sandbox.smtp.mailtrap.io"
}

variable "mailtrap_key" {
  type = string
}

variable "mailtrap_from" {
  type = string
}

variable "mailtrap_port" {
  type = number
}

variable "mailtrap_user" {
  type = string
}

variable "mailtrap_pass" {
  type = string
}

variable "clearstream_api_key" {
  type = string
}

variable "clearstream_text_from" {
  type = string
}

variable "google_key" {
  type = string
}

variable "recaptcha_project" {
  type = string
}
variable "recaptcha_key" {
  type = string
}
variable "recaptcha_action" {
  type = string
}