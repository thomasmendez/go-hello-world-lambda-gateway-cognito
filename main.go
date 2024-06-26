package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Hello World...")

	reqBytes, err := json.Marshal(request)

	if err != nil {
		errMessage := fmt.Sprintf("Error in marshal request: %v", err)
		log.Print(errMessage)
		return events.APIGatewayProxyResponse{
			Headers:    addProxyHeaders("dev"),
			StatusCode: http.StatusInternalServerError,
			Body:       errMessage,
		}, err
	}

	log.Print(string(reqBytes))

	log.Print("Extract token claims...")

	tokenRequest := request.Headers["Authorization"]

	bearerToken := strings.Split(tokenRequest, "Bearer ")

	// TODO: Validate token signature

	claims := jwt.MapClaims{}

	_, _, err = new(jwt.Parser).ParseUnverified(bearerToken[1], claims)

	if err != nil {
		errMessage := fmt.Sprintf("Error in parsing token: %v", err)
		log.Print(errMessage)
		return events.APIGatewayProxyResponse{
			Headers:    addProxyHeaders("dev"),
			StatusCode: http.StatusInternalServerError,
			Body:       errMessage,
		}, err
	}

	for key, value := range claims {
		keyVal := fmt.Sprintf("%s: %v\n", key, value)
		log.Print(keyVal)
	}

	var user User

	if username, ok := claims["cognito:username"].(string); ok {
		fmt.Println("cognito:username:", username)
		user.Username = &username
	} else {
		fmt.Println("cognito:username: claim not found or not a string")
	}

	if email, ok := claims["email"].(string); ok {
		fmt.Println("email:", email)
		user.Email = &email
	} else {
		fmt.Println("email claim not found or not a string")
	}

	responseJson, err := json.Marshal(user)

	if err != nil {
		errMessage := fmt.Sprintf("Error in marshalling user response: %v", err)
		return events.APIGatewayProxyResponse{
			Headers:    addProxyHeaders("dev"),
			StatusCode: http.StatusInternalServerError,
			Body:       errMessage,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Headers:    addProxyHeaders("dev"),
		StatusCode: http.StatusOK,
		Body:       string(responseJson),
	}, nil
}

func addProxyHeaders(env string) map[string]string {
	switch env {
	case "dev":
		return map[string]string{
			"Access-Control-Allow-Origin":  "http://localhost:5173",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		}
	default:
		return map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		}
	}
}

type User struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}
