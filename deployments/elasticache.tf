variable "elasticache_name" {
  default = "filterit-cache"
}

resource "aws_elasticache_subnet_group" "redis_subnet" {
  name       = "redis-subnet"
  subnet_ids = [aws_subnet.filterit-subnet-public-1.id, aws_subnet.filterit-subnet-private-1.id]
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = var.elasticache_name
  engine               = "redis"
  node_type            = "cache.t2.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"
  engine_version       = "3.2.10"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.redis_subnet.name
  security_group_ids   = [aws_security_group.filterit-sg.id]
}
