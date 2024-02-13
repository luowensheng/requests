package main

import "fmt"

func main() {

		// Example GET request
		resp, err := Fetch("https://jsonplaceholder.typicode.com/posts/1").Execute()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	
		// Print response body as string
		bodyString, err := resp.IntoString()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Response body:", bodyString)
	
		// Example POST request with JSON payload
		data := map[string]interface{}{
			"title":  "foo",
			"body":   "bar",
			"userId": 1,
		}
		resp, err = NewRequest("POST", "https://jsonplaceholder.typicode.com/posts").
			JSON(data).
			Execute()
			
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	
		// Print status code of the response
		fmt.Println("Status Code:", resp.StatusCode)
}
