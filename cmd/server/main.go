package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	job "github.com/Oso404/distributed-queue/internal/job"
	internal "github.com/Oso404/distributed-queue/internal/queue"
	worker "github.com/Oso404/distributed-queue/internal/worker"
)

var dq *internal.Queue

func main() {
	dq := internal.Create_Queue("Queue1")
	//create pool of workers
	for i := 0; i < 10; i++ {
		w := worker.Create_Worker()
		go w.Start(dq)
	}
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
		//we only accept json fomat for body
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
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
		w.Write([]byte("Job recieved"))
	})

	fmt.Println("Server listening on :8080")
	//ListenAndServe runs until we stop it
	log.Fatal(http.ListenAndServe(":8080", nil))

}
