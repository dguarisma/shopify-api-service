package server

import (
	"desarrollosmoyan/lambda/src/database"
	"desarrollosmoyan/lambda/src/response"
	"desarrollosmoyan/lambda/src/service"
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"
)

func Lambda() {
	db, errDb := database.GetDb()
	lambda.Start(HandleRequest2(
		db,
		//Product,
		// Reception,
		Inventory,
		// CategoryOne,
		// CategoryTwo,
		// CategoryThree,
		// Reception,
		// Substance,
		// Supplier,
		// Typesproduct,
		// Warehouse,
		// Pack,
		//Trademark,
		// Maker,
		//Purchase,
		errDb,
	))
}

var token = os.Getenv("token")

type handleRequest func(request events.APIGatewayProxyRequest) (ApiResponse events.APIGatewayProxyResponse, err error)

func HandleRequest2(db *gorm.DB, method uint8, errDb error) handleRequest {
	return func(request events.APIGatewayProxyRequest) (ApiResponse events.APIGatewayProxyResponse, err error) {

		ApiResponse = newApiResponse()

		tokenHead := request.MultiValueHeaders["x-amz-security-token"]
		if len(tokenHead) == 0 {
			ApiResponse.StatusCode = 403
			return ApiResponse, nil
		}

		if tokenHead[0] != token {
			ApiResponse.StatusCode = 403
			return ApiResponse, nil
		}

		res := response.Result{}
		filters := service.NewFilter("ID")

		if errDb != nil {
			res := response.NewFailResult(err, http.StatusInternalServerError)
			ApiResponse.Body = res.Body
			ApiResponse.StatusCode = res.Status
			return ApiResponse, nil
		}

		if err != nil {
			res := response.NewFailResult(err, http.StatusInternalServerError)
			ApiResponse.Body = res.Body
			ApiResponse.StatusCode = res.Status
			return ApiResponse, nil
		}

		repo := handleRepository(db, method, request.Path)
		serv := service.NewService(filters, repo)

		switch request.HTTPMethod {
		case "GET":
			res = serv.Get(request.QueryStringParameters)
		case "POST":
			res = serv.Insert(request.Body)
		case "PUT":
			res = serv.Update(request.Body)
		case "DELETE":
			res = serv.Delete(request.Body)
		default:
			res = response.NewFailResult(
				fmt.Errorf("metodo %q no existente", request.HTTPMethod),
				http.StatusBadRequest,
			)
		}

		ApiResponse.Body = res.Body
		ApiResponse.StatusCode = res.Status

		return ApiResponse, nil
	}
}
