###python script to send x number of jobs (same job for now to test out system)


import requests
import time
import random

url = "http://localhost:8080/enqueue"

def send_job(payload):
    try:
        response = requests.post(url, json=payload) #sends POST/enqueue to http://localhost:8080/{looks for enqueue}
        print(f"Sent job: {payload}, Status: {response.status_code}")
    except Exception as e:
        print(f"Failed to send job: {payload}, Error: {e}")

def main():
    num_of_jobs = 100
    for i in range(num_of_jobs):
        #im sending the same job every time just for testing purposes
        job_payload = {
            "payload": {
                "type": "email",
                "sender": "joules",
                "recipient": "oreo"
            },
            "random-note": "helloworld"
        }
        send_job(job_payload)
        time.sleep(random.uniform(0.1, .2)) #delay when sending jobs to simulate real life world 

if __name__ == "__main__":
    main()