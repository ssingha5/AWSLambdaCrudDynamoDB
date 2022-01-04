package main

import (
	"encoding/json"
	"net/http"
	"order/model"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//get the tablename from template.yaml
	OrdersTableName := os.Getenv("ORDERS_TABLE")
	//parse the request body
	order := model.Order{}
	if err := json.Unmarshal([]byte(request.Body), &order); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//Start a DynamoDB session.
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	svc := dynamodb.New(sess)

	// Marshal order to dynamodb attribute value map
	atrrval, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Insert order into DynamoDB table.
	input := &dynamodb.PutItemInput{
		Item:      atrrval,
		TableName: aws.String(OrdersTableName),
	}

	if _, err = svc.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Return the order id
	response, err := json.Marshal(&model.Response{Id: order.Id})
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
