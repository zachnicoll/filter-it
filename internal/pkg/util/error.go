package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func FatalLog(msg string, err error) {
	log.Println(msg)
	log.Fatalln(err)
}

func safeFail(
	clients *Clients,
	metaData *MetaData,
	msg *QueueResponse,
) {

	result, err := clients.DynamoDb.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: metaData.ImageTable,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: msg.DocumentID,
			},
			"date_created": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", msg.DateCreated),
			},
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

	err = UpdateDocument(context.TODO(), clients, metaData.ImageTable, &imageDocument)
	if err != nil {
		FatalLog("Failed to update sqs document for failsafe", err)
	}

	queueMsgBytes, err := json.Marshal(msg)
	if err != nil {
		FatalLog("Failed to marshal sqs message for failsafe", err)
	}

	queueMsgStr := string(queueMsgBytes)

	_, err = clients.SQS.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(queueMsgStr),
		QueueUrl:    metaData.SQSUrl,
	})
	if err != nil {
		FatalLog("Failed to safely put SQS message back into queue for failsafe", err)
	}

	// _, err = client.ASG.SetInstanceProtection(context.TODO(), &autoscaling.SetInstanceProtectionInput{
	// 	InstanceIds:          []string{*metaData.InstanceID},
	// 	AutoScalingGroupName: metaData.ASGName,
	// 	ProtectedFromScaleIn: aws.Bool(true),
	// })
	// if err != nil {
	// 	FatalLog("Failed to disable scale-in protection for failsafe", err)
	// }
}

func SafeFailAndLog(clients *Clients,
	metaData *MetaData,
	sqsMsg *QueueResponse,
	errMsg string, err error) {
	// Log error here first in case safeFail fails and doesn't FatalLog original error
	log.Println(errMsg, err.Error())

	safeFail(clients, metaData, sqsMsg)
	FatalLog(errMsg, err)
}
