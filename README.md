# Scalable Image Filtering Platform (powered by AWS)

## Purpose
This application serves as a platform for uploading images and applying various filters to said images, before sharing the image to the rest of the users on the platform - akin to a mashup of Instagram and DeviantArt.

The project was created as part of a university assignment for a unit in Cloud Computing. It explores the implemenation of various AWS services to create a scalable and persistent application, where it can dynamically respond to changes in load, and uses both long term (S3/DynamoDB) and short term storage (ElastiCache).

## AWS Services

### Architecture Diagram

![Architecture Diagram](./content/architecture_diagram.png)

#### S3

#### DynamoDB

#### ElastiCache

#### Lambda & API Gateway

### SQS

#### CloudWatch

#### AutoScaler

## Monorepo Structure

```
_
|- content // Local files accessed by this README
|- frontend // NextJS clientside application files
|- lambda // Functions deployed to AWS Lambda
    |- feed // Function for retrieving images based on tag
    |- filter_progress // Function for retrieving current progress of image processing
    |- upload // Function for uploading an image and queueing it for processing
|- compute // Services deployed to AWS EC2 instances
    |- image_processor // Service for consuming SQS messages and applying filters to images
```

### Authors
[Zachary Nicoll](https://github.com/zachnicoll)
[Audi Bailey](https://github.com/audibailey)
