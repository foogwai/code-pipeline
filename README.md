# Example Data Pipeline 

## Overview

The Example Data Pipeline is a service that accepts POST requests and produces the data to kafka. 
In this example I have also set up a kafka consumer which reads from the topic and stores the data in an Aerospike database to simulate a data pipeline. 
In a real life scenario, the kafka consumer/Aerospike db would likely be a [snowflake data warehouse with a kafka streaming connector](https://docs.snowflake.com/en/user-guide/data-load-snowpipe-streaming-kafka?gad_source=1)

## Architecture

The code structure follows hexagonal architecture principles, also known as *ports and adapters*. This decouples the different parts of the service which allows for plug and playability with the adapters.  

## Components

- **HTTP Server**: Handles HTTP requests and validates incoming data.
- **Kafka Producer**: Sends POST data to Kafka topics.
- **Kafka Consumer**: Consumes messages from Kafka and processes them.
- **Aerospike Repository**: Stores the processed data in an Aerospike database

## Running the Service
```console
docker-compose up -d --build
```
Once the app container is up and running, POST requests can be sent to localhost:8080

## Endpoints

### POST /submit 
Submits data to the service. The request body should be a JSON object with the following fields:

- ip_address (string): IP address of the user (must be a valid IPv4 address).
- user_agent (string): User agent string from the HTTP request.
- referring_url (string): URL from which the user was referred (must be a valid URL).
- advertiser_id (string): Unique identifier for the advertiser.
- metadata (object): Additional metadata associated with the post data.
Example Request
```json
{
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0",
  "referring_url": "https://example.com",
  "advertiser_id": "12345",
  "metadata": {
    "key": "value"
  }
}
```
Responses:

`204 No Content`: If the data is processed successfully.

`400 Bad Request`: If the request body is invalid.

`500 Internal Server Error`: If an error occurs during processing.

## Performance
Design decisions were made to try and make the server as efficient as possible. Echo was chosen as the web framework due to it's high performance and minimal design. I'm also producing to kafka in batches and consuming the messages in a go routine. I ran load testing on the server use the Vegeta load testing tool and the server is quite performant. With only one kafka topic/partition and the infra running in limited memory local containers, I recieved the following results:
```
vegeta attack -rate 0 -duration 10s -targets=target.txt -max-workers 2000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         173150, 17315.90, 16486.92
Duration      [total, attack, wait]             10.502s, 9.999s, 502.78ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.251ms, 108.941ms, 100.345ms, 181.812ms, 220.749ms, 335.312ms, 652.996ms
Bytes In      [total, mean]                     0, 0.00
Bytes Out     [total, mean]                     56793200, 328.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      204:173150
```
17315 requests per second with very limited hardware. In a scaled production environment this could handle many many times that. It's also possible that tweaking vegeta settings could produce even better results. 

## Work TODO Before Production and Considerations Made

- Obviously switching out the kafka consumer/aerospike db for the real data warehouse back end would need to be done  
- More structured logging using something like logorus. 
- I'm currently only testing the handlers, but we would want unit tests for all the functions before push to prod.
- CICD integration
- Security would need to be considered, for example the CORS config currently allows all origins, we'd want to change that. 
