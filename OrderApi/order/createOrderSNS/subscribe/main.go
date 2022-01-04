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

func handler(snsEvent events.SNSEvent) error {
	//get the tablename from template.yaml
	OrdersTableName := os.Getenv("ORDERS_TABLE")
	//Start a DynamoDB session.
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	svc := dynamodb.New(sess)

	for _, message := range snsEvent.Records {
		snsRecord := message.SNS
		fmt.Println("message ======== ", message)
		fmt.Println("snsRecord ======== ", snsRecord)

		idMap := snsRecord.MessageAttributes["Id"].(map[string]interface{})
		idstr := fmt.Sprint(idMap["Value"])
		id, err := strconv.Atoi(idstr)
		if err != nil {
			fmt.Println("can't convert String to int")
		}
		itemNameMap := snsRecord.MessageAttributes["ItemName"].(map[string]interface{})
		itemName := fmt.Sprint(itemNameMap["Value"])

		quantityMap := snsRecord.MessageAttributes["Quantity"].(map[string]interface{})
		quantitystr := fmt.Sprint(quantityMap["Value"])
		quantity, err := strconv.Atoi(quantitystr)
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
