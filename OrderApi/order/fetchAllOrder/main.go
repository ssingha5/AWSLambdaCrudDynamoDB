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

	//Start DynamoDB session
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	svc := dynamodb.New(sess)

	//get Table details
	input := &dynamodb.ScanInput{
		TableName: aws.String(OrdersTableName),
	}

	//get all the records from the table
	result, err := svc.Scan(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//mab it to array of Order
	obj := []model.Order{}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//Return all the orders
	response, err := json.Marshal(obj)
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
