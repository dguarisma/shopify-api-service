package server

import "github.com/aws/aws-lambda-go/events"

func newApiResponse() (ApiResponse events.APIGatewayProxyResponse) {
	ApiResponse = events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE",
			//			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Headers": "*",
			"AllowCredentials":             "true",
			"Content-Type":                 "application/json",
		},
		Body: "null",
	}
	return ApiResponse
}
