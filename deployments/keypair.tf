variable "key_name" {
  default = "filterit-terraform"
}

resource "aws_key_pair" "image_processor_key" {
  key_name   = var.key_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCSB/ZHcGkC5bZNb2EcNqyDOk7HnzbGSWhBcP/MuQuqwr+Dt+rUy9cPIooFAUhI2ixFDiEFayt1Ias3N5BEqpMgYdk3e2mpelmr9W06ojpcNdiw5Oxtf9zUs2xV9ieXJfRGS+u/ILsAX+DyW5inxNCk3u2fg87PBoYfpRJ1cyxv7Uyi2/S2v4Qq9u19YFU75NunRovpfheYSOSVh3it3u3Vc+QHRnl6ziRH3cn0zrPZ8PLt2wDLsH4DS+jLSczKAA+ZhaQq4JIbNFQfixB3sWoRJ8Z4mjciFVIO9EX/IIF84WImNhSN0QfaxwNnmsIR7da7zCrFQIhu9JLP8WFez7zn"
}
