package filterit

import "fmt"

func applyFilters(filters []int) {
	// TODO: Use image magik to apply each of the filters to the image
}

func processMessage(id string) {
	// TODO: Fetch document form DynamoDb with respective id

	// TODO: Set document's progress attribute to PROCESSING

	// TODO: Fetch image from S3 based on the document's image attribute

	// TODO: Apply filters to image

	// TODO: Re-upload filtered image to S3

	// TODO: Write new filenname to document's image attribute

	// TODO: Set document progress attribute to DONE

	// TODO: Invalid cache with all keys containing filters (use KEYS command)
}

func WatchQueue() {
	//// Get instance ID
	//instanceID := util.FetchInstanceID()
	//
	//// Get AutoScaling Group
	//asg := os.Getenv("AS_GROUP")
	//if asg == "" {
	//	log.Fatalln("Unable to find autoscaling group")
	//}
	//
	//// Get S3 Bucket
	//s3Bucket := os.Getenv("S3_BUCKET")
	//if s3Bucket == "" {
	//	log.Fatalln("Unable to find s3 bucket")
	//}
	//
	//// Get DynamoDB Table
	//imageTable := os.Getenv("AWS_IMAGE_TABLE")
	//if imageTable == "" {
	//	log.Fatalln("Unable to find dynamodb table")
	//}
	//
	//// Get SQS Queue
	//sqsQueue := os.Getenv("AWS_SQS_QUEUE")
	//if sqsQueue == "" {
	//	log.Fatalln("Unable to find sqs queue name")
	//}
	//
	//// Get AWS Region
	//awsRegion := os.Getenv("AWS_REGION")
	//if awsRegion == "" {
	//	log.Fatalln("Unable to find aws region")
	//}
	//
	//// Initialise AWS Session Config
	//awsSessionConfig, err := session.NewSessionWithOptions(session.Options{
	//	Config: aws.Config{
	//		Region: aws.String(awsRegion),
	//	},
	//
	//	SharedConfigState: session.SharedConfigEnable,
	//})
	//if err != nil {
	//	util.FatalLog("Unable to initialise aws session config", err)
	//}
	//
	//// AWS DynamoDB session
	//awsDBSession := dynamodb.New(awsSessionConfig)
	//
	//// AWS SQS session
	//awsSQSSession := sqs.New(awsSessionConfig)
	//
	//// AWS S3 session
	//awsS3Session := s3.New(awsSessionConfig)
	//
	//// AWS ASG session
	//awsASGSession := autoscaling.New(awsSessionConfig)
	//
	//// SQS URL
	//urlResult, err := awsSQSSession.GetQueueUrl(&sqs.GetQueueUrlInput{
	//	QueueName: &sqsQueue,
	//})
	//if err != nil {
	//	util.FatalLog("Unable to fetch aws sqs url", err)
	//}
	//
	//for {
	//	// SQS Intake Parameters
	//	params := &sqs.ReceiveMessageInput{
	//		QueueUrl: urlResult.QueueUrl, // Required
	//		MaxNumberOfMessages: aws.Int64(1),
	//		MessageAttributeNames: []*string{
	//			aws.String("All"),
	//		},
	//		WaitTimeSeconds: aws.Int64(30),
	//	}
	//
	//	// Receive an SQS Message
	//	resp, err := awsSQSSession.ReceiveMessage(params)
	//	if err != nil {
	//		util.FatalLog("Unable to fetch aws sqs message", err)
	//	}
	//
	//	// Check a message was received
	//	if len(resp.Messages) == 1 {
	//		_, err = awsASGSession.SetInstanceProtection(&autoscaling.SetInstanceProtectionInput{
	//			InstanceIds: []*string{aws.String(instanceID)},
	//			AutoScalingGroupName: aws.String(asg),
	//			ProtectedFromScaleIn: aws.Bool(true),
	//		})
	//		if err != nil {
	//			util.FatalLog("Unable to enable scale-in protection", err)
	//		}
	//
	//		targetMessage := resp.Messages[0]
	//
	//		_, err =  awsSQSSession.DeleteMessage(&sqs.DeleteMessageInput{
	//			QueueUrl:      urlResult.QueueUrl, // Required
	//			ReceiptHandle: targetMessage.ReceiptHandle, // Required
	//		})
	//		if err != nil {
	//			util.FatalLog("Unable to delete aws sqs message", err)
	//		}
	//
	//		// Get DynamoDB image info from message body
	//		result, err := awsDBSession.GetItem(&dynamodb.GetItemInput{
	//			TableName: aws.String(imageTable),
	//			Key: map[string]*dynamodb.AttributeValue{
	//				"id": {
	//					N: targetMessage.Body,
	//				},
	//				"progress": {
	//					S: aws.String("READY"),
	//				},
	//			},
	//		})
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to check aws sqs dynamodb status", err)
	//		}
	//
	//		// Unmarshal image document
	//		var imageDocument util.ImageDocument
	//		err = dynamodbattribute.UnmarshalMap(result.Item, &imageDocument)
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to unmarshal aws sqs dynamodb status", err)
	//		}
	//
	//		imageDocument.Progress = util.PROCESSING
	//
	//		imageDocumentMap, err := dynamodbattribute.MarshalMap(imageDocument)
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to marshal aws sqs dynamodb (processing)", err)
	//		}
	//
	//		_, err = awsDBSession.PutItem(&dynamodb.PutItemInput{
	//			Item:      imageDocumentMap,
	//			TableName: aws.String(imageTable),
	//		})
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to update aws sqs dynamodb (processing)", err)
	//		}
	//
	//		// TODO: IMAGE
	//		_, err = awsS3Session.GetObject(&s3.GetObjectInput{
	//			Bucket: aws.String(s3Bucket),
	//			Key:    aws.String(imageDocument.Image),
	//		})
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to get S3 image", err)
	//		}
	//
	//		// PROCESS IMAGE
	//
	//		// Generate Image UUID
	//		imageID := uuid.New()
	//
	//		// String format S3 image name and generate a S3 put object
	//		imageName := fmt.Sprintf("uploads/%s.jpg", imageID.String())
	//		_, err = awsS3Session.PutObject(&s3.PutObjectInput{
	//			Bucket: aws.String(s3Bucket),
	//			Key:    aws.String(imageName),
	//			// TODO: IMAGE
	//			//Body: ,
	//		})
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to put new S3 image", err)
	//		}
	//
	//		imageDocument.Progress = util.DONE
	//		imageDocument.Image = imageName
	//
	//		imageDocumentMap, err = dynamodbattribute.MarshalMap(imageDocument)
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to marshal aws sqs dynamodb (processed)", err)
	//		}
	//
	//		_, err = awsDBSession.PutItem(&dynamodb.PutItemInput{
	//			Item:      imageDocumentMap,
	//			TableName: aws.String(imageTable),
	//		})
	//		if err != nil {
	//			util.SafeFail(awsSQSSession, awsDBSession, awsASGSession, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
	//			util.FatalLog("Unable to update aws sqs dynamodb (processed)", err)
	//		}
	//	}
	//}

	fmt.Println("Monka Moment")

	// TODO: In a loop, check if the SQS queue has a new message

	// TODO: If message, spin off a subroutine and process the message - processMessage(id)

	// TODO: Log custom metric to CloudWatch

}
