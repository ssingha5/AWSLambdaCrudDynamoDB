package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order/model"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//parse the request body
	order := model.Order{}
	if err := json.Unmarshal([]byte(request.Body), &order); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//gets the topic URL of an Amazon SNS
	topic := os.Getenv("ORDER_SNS_TOPIC")
	fmt.Println("topic ========== ", topic)
	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an SQS service client
	svc := sns.New(sess)

	//SendMsg sends a message to an Amazon SQS queue
	_, err := svc.Publish(&sns.PublishInput{
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Id": &sns.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(order.Id)),
			},
			"ItemName": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(order.ItemName),
			},
			"Quantity": &sns.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(order.Quantity)),
			},
		},
		Message:  aws.String("Publishing Book"),
		TopicArn: &topic,
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
