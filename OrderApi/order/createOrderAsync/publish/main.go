package main

import (
	"encoding/json"
	"net/http"
	"order/model"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//parse the request body
	order := model.Order{}
	if err := json.Unmarshal([]byte(request.Body), &order); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//gets the URL of an Amazon SQS queue
	queue := os.Getenv("ORDER_QUEUE")

	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an SQS service client
	svc := sqs.New(sess)

	//SendMsg sends a message to an Amazon SQS queue
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Id": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(order.Id)),
			},
			"ItemName": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(order.ItemName),
			},
			"Quantity": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(order.Quantity)),
			},
		},
		MessageBody: aws.String("Publishing Order"),
		QueueUrl:    &queue,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response, err := json.Marshal("Order with Id : " + strconv.Itoa(order.Id) + " published successfully")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(handler)
}
