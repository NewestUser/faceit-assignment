## A simple User Rest Service for FACEIT

#### Owner: Marian Zlatev

### Things to consider when reviewing this solution

- I haven't had time to implement all the requirements. Functionality like listing/paging is ommited.
- I haven't handled cases like hashing sensitive data like user passwords.
- I haven't enforced uniqueness for username or email.
- I have omitted stuff like connection pooling, connection timeouts etc...
- Logging and error handling can be greatly improved.
- I haven't introduce another layer between HTTP Handles and Database interactions.
- Health Check.
- Paging API with filtering.
- Fine tuning kafka.
- Connection pooling and properly closing db sessions.

### Getting started

The application uses:
- Mongo as a persistence layer
- Kafka for messaging

**IMPORTANT:** I have come across a cross compilation issue with the package `github.com/confluentinc/confluent-kafka-go/kafka`
This issue manifest particularly when building the project inside a Docker container.
I haven't had time to fix it but you can see here [issues 119](https://github.com/confluentinc/confluent-kafka-go/issues/119)
that it is still open.  
Depending on the OS you are running you might experience docker build issues.

#### Build and run tests
Note that the tests are **End 2 End**. The application will be started in order to run the test.

``` shell script
docker-compose up faceitapi-tets
```

#### Run application for development
```shell script
# Start dependencies
docker-compose up zookeeper kafka mongo

# start the app from terminal (or you can run it in your IDE)
go mod download && cd cmd && go build -o faceit-app && ./faceit-app
```

### API

```shell script
# Create user
curl -X POST localhost:8080/users -d '{"firstName":"Vincent", "lastName":"Furnier", "nickName":"Vincent", "password":"qwerty", "email":"alice.cooper@mail.com", "country":"US"}'

# Fetch created user (id should be taken from previous request)
curl localhost:8080/users/5f9339159f1dd812c9eab2a6

# Update existing user (id should be taken from previous request)
curl -X PUT localhost:8080/users/5f9339159f1dd812c9eab2a6 -d '{"firstName":"Alice", "lastName":"Cooper", "nickName":"Vincent", "password":"qwerty", "email":"alice.cooper@mail.com", "country":"US"}'
```
