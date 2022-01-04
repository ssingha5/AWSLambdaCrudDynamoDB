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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//get the tablename from template.yaml
	OrdersTableName := os.Getenv("ORDERS_TABLE")
	//Get id parameter
	orderId := request.PathParameters["id"]
	//Start DynamoDB session
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	svc := dynamodb.New(sess)

	//Read order item from DynamoDB table.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(OrdersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(orderId),
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if result.Item == nil {
		msg := "Could not find order with Id : " + orderId
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       string(msg),
		}, nil
	}
	//map the result item from DB to Order object
	order := model.Order{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &order); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//Return all the orders
	response, err := json.Marshal(order)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil

}

func main() {
	lambda.Start(Handler)
}
