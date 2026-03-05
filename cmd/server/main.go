package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	job "github.com/Oso404/distributed-queue/internal/job"
	internal "github.com/Oso404/distributed-queue/internal/queue"
)

var dq *internal.Queue

func main() {
	dq := internal.Create_Queue("Queue1")

	/*
		HandleFunc registers a handler function and runs whenever someone visits URL path (e.g."/")
		HandleFunc takes in URL pattern and function to run when visited
		ResponseWriter allows to write a response back to client
		Request contains metadata about incoming HTTP request
	*/
	//sample handler function for demonstration
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world! The server is running.")
	})

	http.HandleFunc("/enqueue", func(w http.ResponseWriter, r *http.Request) {
		//check to see that method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}
		/*
			r.Body is stream of bytes as http request comes in
			r.Body is accessible only once thus we can only read it once
			io.ReadAll(r.Body) takes stream of bytes and converts it into byte slice for us to use later
		*/
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		job := job.Create_job(body)
		dq.Enqueue(job)
		fmt.Println("Pending queue", dq.PendingQueue)
		fmt.Println("All jobs", dq.Jobs)
		// /*
		// 	data will contain all of information passed in request (type and payload)
		// 	type(string) : (string)
		// 	payload(string) : json (string: string)
		// */
		// //why string:interface? interface represents any data type

		// var data map[string]interface{}
		// //Unmarshal([]byte, v any) tries to convert stream of bytes into json
		// //Unmarshal requires stream of json bytes and pointer to variable (must be pointer to modify variable!)
		// err = json.Unmarshal(body, &data)
		// if err != nil {
		// 	fmt.Println("Invalid JSON:", err)
		// 	return
		// }
		// //check if payload field exists
		// rawPayload, ok := data["payload"]
		// if !ok {
		// 	http.Error(w, "Missing payload field!", http.StatusBadRequest)
		// 	return
		// }
		// //check if payload field contains json object
		// payload, ok := rawPayload.(map[string]interface{})
		// if !ok {
		// 	http.Error(w, "Payload field not JSON object", http.StatusBadRequest)
		// 	return
		// }
		// fmt.Println("Payload map:", payload)
		w.Write([]byte("Job recieved"))
	})

	fmt.Println("Server listening on :8080")
	//ListenAndServe runs until we stop it
	log.Fatal(http.ListenAndServe(":8080", nil))
}
