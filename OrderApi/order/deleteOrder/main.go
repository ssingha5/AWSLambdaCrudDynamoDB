package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//get the tablename from template.yaml
	OrdersTableName := os.Getenv("ORDERS_TABLE")
	//Get id parameter url : https://lrs93tlh1k.execute-api.us-east-1.amazonaws.com/Prod/orders?id=448
	orderId := request.QueryStringParameters["id"]

	//Start DynamoDB session
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	svc := dynamodb.New(sess)

	//Delete the order item from DynamoDB table.
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(OrdersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(orderId),
			},
		},
	}

	if _, err = svc.DeleteItem(input); err != nil {
		msg := "Could not find order with Id : " + orderId
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       string(msg),
		}, nil
	}

	//Return the response
	response, err := json.Marshal("Order with Id : " + orderId + " deleted successfully")
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
