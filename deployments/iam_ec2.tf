resource "aws_iam_instance_profile" "ec2_iam_profile" {
  name = "ec2_iam_profile"
  role = aws_iam_role.ec2_exec_role.name
}


resource "aws_iam_role" "ec2_exec_role" {
  name = "ec2_exec_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ec2_exec_sqs" {
  name        = "ec2_exec_sqs"
  path        = "/"
  description = "IAM policy for sqs from ec2"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "sqs:RecieveMessage",
        "sqs:SendMessage",
        "sqs:GetQueueUrl"
      ],
      "Resource": "${aws_sqs_queue.filterit-queue.arn}",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ec2_exec_s3" {
  name        = "ec2_exec_s3"
  path        = "/"
  description = "IAM policy for s3 from ec2"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "${aws_s3_bucket.image_bucket.arn}/*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ec2_exec_dynamodb" {
  name        = "ec2_exec_dynamodb"
  path        = "/"
  description = "IAM policy for dynamodb from ec2"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "dynamodb:BatchGetItem",
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:BatchWriteItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem"
      ],
      "Resource": [
        "${aws_dynamodb_table.ddbtable.arn}",
        "${aws_dynamodb_table.ddbtable.arn}/index/*"
      ],
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ec2_policy" {
  role       = aws_iam_role.ec2_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

resource "aws_iam_role_policy_attachment" "ec2_sqs" {
  role       = aws_iam_role.ec2_exec_role.name
  policy_arn = aws_iam_policy.ec2_exec_sqs.arn
}

resource "aws_iam_role_policy_attachment" "ec2_s3" {
  role       = aws_iam_role.ec2_exec_role.name
  policy_arn = aws_iam_policy.ec2_exec_s3.arn
}

resource "aws_iam_role_policy_attachment" "ec2_dynamodb" {
  role       = aws_iam_role.ec2_exec_role.name
  policy_arn = aws_iam_policy.ec2_exec_dynamodb.arn
}
