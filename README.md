# Scalable Image Filtering Platform (powered by AWS)

## Purpose
This application serves as a platform for uploading images and applying various filters to said images, before sharing the image to the rest of the users on the platform - akin to a mashup of Instagram and DeviantArt.

The project was created as part of a university assignment for a unit in Cloud Computing. It explores the implemenation of various AWS services to create a scalable and persistent application, where it dynamically responds to changes in load, and uses both long term (S3/DynamoDB) and short term storage (ElastiCache).

## Deployment
Terraform is used to deploy the various AWS services used in this application. To deploy the infrastructure:

```sh
cd deployments
terraform init
terraform apply
```

## AWS Services

### Architecture Diagram

![Architecture Diagram](./assets/architecture_diagram.png)

### Authors

- [Zachary Nicoll](https://github.com/zachnicoll)

- [Audi Bailey](https://github.com/audibailey)
