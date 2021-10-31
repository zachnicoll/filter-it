package util

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func FatalLog(msg string, err error) {
	log.Println(msg)
	log.Fatalln(err)
}

func safeFail(
	client *Clients,
	metaData *MetaData,
	msg *sqsTypes.Message,
) {
	documentID := msg.Body

	result, err := client.DynamoDb.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: metaData.ImageTable,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: *documentID},
		},
	})
	if err != nil {
		FatalLog("Failed to get sqs document for failsafe", err)
	}

	// Unmarshal image document
	var imageDocument ImageDocument
	err = attributevalue.UnmarshalMap(result.Item, &imageDocument)
	if err != nil {
		FatalLog("Failed to unmarshal sqs document for failsafe", err)
	}

	imageDocument.Progress = READY

	imageDocumentMap, err := attributevalue.MarshalMap(imageDocument)
	if err != nil {
		FatalLog("Failed to marshal sqs document for failsafe", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      imageDocumentMap,
		TableName: metaData.ImageTable,
	}
	_, err = client.DynamoDb.PutItem(context.TODO(), input)
	if err != nil {
		FatalLog("Failed to update sqs document for failsafe", err)
	}

	_, err = client.SQS.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: documentID,
		QueueUrl:    metaData.SQSUrl,
	})
	if err != nil {
		FatalLog("Failed to safely put SQS message back into queue for failsafe", err)
	}

	_, err = client.ASG.SetInstanceProtection(context.TODO(), &autoscaling.SetInstanceProtectionInput{
		InstanceIds:          []string{*metaData.InstanceID},
		AutoScalingGroupName: metaData.ASGName,
		ProtectedFromScaleIn: aws.Bool(true),
	})
	if err != nil {
		FatalLog("Failed to disable scale-in protection for failsafe", err)
	}
}

func SafeFailAndLog(clients *Clients,
	metaData *MetaData,
	sqsMsg *sqsTypes.Message,
	errMsg string, err error) {
	safeFail(clients, metaData, sqsMsg)
	FatalLog(errMsg, err)
}
