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

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":iName": {
				S: aws.String(order.ItemName),
			},
			":qty": {
				N: aws.String(strconv.Itoa(order.Quantity)),
			},
		},
		TableName: aws.String(OrdersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.Itoa(order.Id)),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set itemName = :iName,quantity = :qty"),
	}

	//updating the order
	result, err := svc.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//map the updated result item from DB to Order object
	updatedOrder := model.Order{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &updatedOrder); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	updatedOrder.Id = order.Id

	//Return all the orders
	response, err := json.Marshal(updatedOrder)
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
