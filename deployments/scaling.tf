variable "asg_name" {
  default = "filterit-asg"
}

variable "asp_name_down" {
  default = "filterit-asp-down"
}

variable "asp_name_up" {
  default = "filterit-asp-up"
}

variable "asp_alarm_down" {
  default = "filterit-alarm-down"
}

variable "asp_alarm_up" {
  default = "filterit-alarm-up"
}

resource "aws_autoscaling_group" "filterit-asg" {
  name = var.asg_name

  min_size             = 1
  desired_capacity     = 1
  max_size             = 4

  vpc_zone_identifier = [aws_subnet.filterit-subnet-public-1.id, aws_subnet.filterit-subnet-private-1.id]

  launch_configuration = aws_launch_configuration.filterit-lc.name

  enabled_metrics = [
    "GroupMinSize",
    "GroupMaxSize",
    "GroupDesiredCapacity",
    "GroupInServiceInstances",
    "GroupTotalInstances"
  ]

  metrics_granularity = "1Minute"

  # Required to redeploy without an outage.
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_launch_configuration" "filterit-lc" {
  name_prefix = "filterit-"

  image_id = aws_ami_from_instance.image_processor_ami.id
  instance_type = "t2.micro"
  key_name = aws_key_pair.image_processor_key.key_name
  iam_instance_profile = aws_iam_instance_profile.ec2_iam_profile.name

  security_groups = [ aws_security_group.filterit-sg.id ]
  associate_public_ip_address = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_policy" "filterit-servers-down" {
  name = var.asp_name_down
  autoscaling_group_name = aws_autoscaling_group.filterit-asg.name
  scaling_adjustment = -1
  cooldown = 300
  adjustment_type = "ChangeInCapacity"
}

resource "aws_autoscaling_policy" "filterit-servers-up" {
  name = var.asp_name_up
  autoscaling_group_name = aws_autoscaling_group.filterit-asg.name
  scaling_adjustment = 1
  cooldown = 300
  adjustment_type = "ChangeInCapacity"
}

resource "aws_cloudwatch_metric_alarm" "filter-it-add-workers" {
  alarm_name = var.asp_alarm_up
  metric_name = "ApproximateNumberOfMessagesVisible"
  namespace = "AWS/SQS"
  statistic = "Sum"
  period = "60"
  threshold = "6"
  comparison_operator = "GreaterThanOrEqualToThreshold"

  evaluation_periods  = "1"

  dimensions = {
    QueueName = aws_sqs_queue.filterit-queue.name
  }

  alarm_actions = [aws_autoscaling_policy.filterit-servers-up.arn]
}

resource "aws_cloudwatch_metric_alarm" "filter-it-remove-workers" {
  alarm_name = var.asp_alarm_down
  metric_name = "ApproximateNumberOfMessagesVisible"
  namespace = "AWS/SQS"
  statistic = "Sum"
  period = "60"
  threshold = "2"
  comparison_operator = "LessThanOrEqualToThreshold"

  evaluation_periods  = "1"

  dimensions = {
    QueueName = aws_sqs_queue.filterit-queue.name
  }

  alarm_actions = [aws_autoscaling_policy.filterit-servers-down.arn]
}