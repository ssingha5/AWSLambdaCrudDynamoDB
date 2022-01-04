package main

import (
	"fmt"
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

func handler(sqsEvent events.SQSEvent) error {
	//get the tablename from template.yaml
	OrdersTableName := os.Getenv("ORDERS_TABLE")
	//Start a DynamoDB session.
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	svc := dynamodb.New(sess)

	for _, message := range sqsEvent.Records {
		id, err := strconv.Atoi(*message.MessageAttributes["Id"].StringValue)
		if err != nil {
			fmt.Println("can't convert String to int")
		}
		itemName := *message.MessageAttributes["ItemName"].StringValue
		quantity, err := strconv.Atoi(*message.MessageAttributes["Quantity"].StringValue)
		if err != nil {
			fmt.Println("can't convert String to int")
		}
		order := model.Order{id, itemName, quantity}
		// Marshal order to dynamodb attribute value map
		atrrval, err := dynamodbattribute.MarshalMap(order)
		if err != nil {
			fmt.Println(err)
		}

		// Insert order into DynamoDB table.
		input := &dynamodb.PutItemInput{
			Item:      atrrval,
			TableName: aws.String(OrdersTableName),
		}

		if _, err = svc.PutItem(input); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Order with id " + strconv.Itoa(order.Id) + " inserted successfully")

	}

	return nil
}

func main() {
	lambda.Start(handler)
}
