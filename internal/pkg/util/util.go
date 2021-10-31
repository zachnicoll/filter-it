package util

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func InternalServerError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET, PUT, PATCH, POST, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Authorization, Content-Type",
		},
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}
}

func JSONStringResponse(body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET, PUT, PATCH, POST, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Authorization, Content-Type",
		},
		Body: body,
	}
}

func SortFilters(filters []int) {
	sort.Slice(
		filters,
		func(i, j int) bool {
			return i > j
		},
	)
}

/*
Sort ImageDocument slice by DateCreated, with latest first.
*/
func SortDocuments(documents []ImageDocument) {
	sort.Slice(
		documents,
		func(i, j int) bool {
			return documents[i].DateCreated < documents[j].DateCreated
		},
	)
}

/*
Create a DynamoDB Expression containing conditions where each document retrieved
must include each of the supplied filters in the "filters" attribute.
*/
func BuildFilterConditions(filter string) (expression.Expression, error) {
	builder := expression.NewBuilder()

	filterInt, err := strconv.Atoi(filter)

	if err == nil {
		filterCondition := expression.Name("filter").Equal(expression.Value(filterInt))
		builder = builder.WithCondition(filterCondition)
	}

	// Make sure that only DONE documents are selected from DynamoDB
	progressCondition := expression.Name("progress").Equal(expression.Value(DONE))

	builder = builder.WithCondition(progressCondition)

	return builder.Build()
}

func FetchInstanceID() string {
	response, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.Fatalln("Unable to find instance ID")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("Unable to close instance ID response")
		}
	}(response.Body)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("Unable to read instance ID response")
	}

	instanceID := string(responseData)
	fmt.Println(instanceID)

	return instanceID
}

func FatalLog(msg string, err error) {
	log.Println(msg)
	log.Fatalln(err)
}

func SafeFail(
	sqsService *sqs.Client,
	dbService *dynamodb.Client,
	asgService *autoscaling.Client,
	asg string,
	instanceID string,
	imageTable string,
	queueURL string,
	documentID string,
) {
	result, err := dbService.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(imageTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: documentID},
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
		TableName: aws.String(imageTable),
	}
	_, err = dbService.PutItem(context.TODO(), input)
	if err != nil {
		FatalLog("Failed to update sqs document for failsafe", err)
	}

	_, err = sqsService.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(documentID),
		QueueUrl:    aws.String(queueURL),
	})
	if err != nil {
		FatalLog("Failed to safely put SQS message back into queue for failsafe", err)
	}

	_, err = asgService.SetInstanceProtection(context.TODO(), &autoscaling.SetInstanceProtectionInput{
		InstanceIds:          []string{instanceID},
		AutoScalingGroupName: aws.String(asg),
		ProtectedFromScaleIn: aws.Bool(true),
	})
	if err != nil {
		FatalLog("Failed to disable scale-in protection for failsafe", err)
	}
}
