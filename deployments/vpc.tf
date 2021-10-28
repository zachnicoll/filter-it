resource "aws_vpc" "filterit-vpc" {
  cidr_block = "10.0.0.0/16"
  enable_dns_support = "true"
  enable_dns_hostnames = "true"
  enable_classiclink = "false"
  instance_tenancy = "default"
}

resource "aws_subnet" "filterit-subnet-public-1" {
  vpc_id                  = aws_vpc.filterit-vpc.id
  cidr_block              = "10.0.0.0/21"
  map_public_ip_on_launch = true
  availability_zone = "ap-southeast-2a"
}

resource "aws_internet_gateway" "filterit-internet-gateway" {
  vpc_id = aws_vpc.filterit-vpc.id
}

resource "aws_route_table" "filterit-routetable-public-1" {
  vpc_id = aws_vpc.filterit-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.filterit-internet-gateway.id
  }
}

resource "aws_route_table_association" "filterit-routetableassoc-public-1" {
  subnet_id      = aws_subnet.filterit-subnet-public-1.id
  route_table_id = aws_route_table.filterit-routetable-public-1.id
}

resource "aws_eip" "filterit-eip" {
  vpc        = true
  depends_on = [aws_internet_gateway.filterit-internet-gateway]
}

resource "aws_nat_gateway" "filterit-nat-gateway" {
  allocation_id = aws_eip.filterit-eip.id
  subnet_id     = aws_subnet.filterit-subnet-public-1.id
}

// PRIVATE
resource "aws_subnet" "filterit-subnet-private-1" {
  vpc_id = aws_vpc.filterit-vpc.id
  cidr_block = "10.0.8.0/24"
  map_public_ip_on_launch = "false"
  availability_zone = "ap-southeast-2a"
}

resource "aws_route_table" "filterit-routetable-private-1" {
  vpc_id = aws_vpc.filterit-vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.filterit-nat-gateway.id
  }
}

resource "aws_route_table_association" "filterit-routetableassoc-private-1" {
  subnet_id = aws_subnet.filterit-subnet-private-1.id
  route_table_id = aws_route_table.filterit-routetable-private-1.id
}

// Security Group
resource "aws_security_group" "filterit-sg" {
  vpc_id = aws_vpc.filterit-vpc.id

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}