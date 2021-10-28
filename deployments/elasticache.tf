variable "elasticache_name" {
  default = "filterit-cache"
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = var.elasticache_name
  engine               = "redis"
  node_type            = "cache.t2.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"
  engine_version       = "3.2.10"
  port                 = 6379
}
