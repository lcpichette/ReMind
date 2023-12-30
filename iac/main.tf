# CORE

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.31"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "us-west-2"
  profile = "default"
}

# DEV S3

resource "aws_s3_bucket" "example" {
  bucket = "dev.re.mind"
}

resource "aws_s3_bucket_ownership_controls" "example" {
  bucket = aws_s3_bucket.example.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "example" {
  bucket = aws_s3_bucket.example.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "example" {
  depends_on = [
    aws_s3_bucket_ownership_controls.example,
    aws_s3_bucket_public_access_block.example,
  ]

  bucket = aws_s3_bucket.example.id
  acl    = "public-read"
}

# DEV EC2

data "aws_ami" "amzn2" {
  owners      = ["amazon"]
  most_recent = true
  filter {
    name = "name"
    values = ["amzn2-ami-hvm-2.*-x86_64-gp2"]
  }
}

resource "aws_key_pair" "mac2" {
  key_name   = "mac2"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCzK2QPQC1+pEXiqINbWrcJ0Y5dtTI6qP4KdA4mZBjBGaomiTsbf1ZQgoCHyS5zhIi8eSHzwXSYNMpLb9XJ+WXK+GH3wfR06KSqBcdpn23HUHEVzxZMfdulUzs1dufH1TBXww7vXjeVLdcAOCXE7JLSstj3zB7I6uTimX5u4iiWid6UKNJDr9t8dbHvhTglXCs+WiX0wFPMxC/oEOaiwvdCIftffpuqobEDce0DEecPQihXc9EE0mxc0GjRYOlzqRNaELhwF+xb3mCgguAUiTfzvardCIz/QpoAtXfaet4NEybajHn+8NTYEICfQF81OdoLt9YMtJjI4TIpzGFTZ5ERTlyjTpMuGQ2zUb5JZ8BY8pUHeRY0gWQGWU7lVgXnJHqY4v1lc7sAlEB2Xr5KzGqYE68ezdq3a44+wVx+htLa/k/xn7k5Jb0qyFLCB73eqQ1L+kld7BO6KgXwG4btkYL52Owq4QsBgNBTKp6NHlUJzTvgi2beaQn8L4TbJN6O43M= lucaspichette@Lucass-MacBook-Pro.local"
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.amzn2.id
  instance_type = "t2.micro"
  key_name = aws_key_pair.mac2.key_name

  tags = {
    Name = "ReMind"
  }
}
