module "db" {
  source  = "cloudposse/rds/aws"
  version = "1.1.2"

  allocated_storage           = 20
  storage_type                = "gp2"
  allow_major_version_upgrade = false
  security_group_ids = [
    aws_security_group.serve_vm.id
  ]
  database_name           = var.db_name
  database_user           = var.db_user
  database_password       = var.db_pass
  database_port           = var.db_port
  db_parameter_group      = "default.postgres15"
  db_subnet_group_name    = "serve-uw2"
  parameter_group_name    = "default.postgres15"
  option_group_name       = "default:postgres-15"
  engine                  = "postgres"
  engine_version          = "15.12"
  instance_class          = var.db_instance
  name                    = "serve-db-1"
  subnet_ids              = data.aws_db_subnet_group.this.subnet_ids
  vpc_id                  = data.aws_db_subnet_group.this.vpc_id
  skip_final_snapshot     = true
  backup_retention_period = 30
  backup_window           = "22:00-03:00"
}