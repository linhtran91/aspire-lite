# aspire-lite

## Install

It has 2 ways to run the application:
- Run via docker-compose
- Install go and other components

### User docker-compose

The application is using postgres as a database, can use the docker-compose to run it via command
```sh
./start.sh
```

### Install Go and another components

The application is written with Go 1.21.4 and uses Postgresql as a database. Should ensure that both are installed to run the application manually. Go to this link to [download Go](https://go.dev/doc/install)
And also go-migrate to run the data migration as below

Mac installation : 
```sh
brew install golang-migrate
```

Pre-built binary (Windows, MacOS, or Linux) : 
```sh
curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz
```
More detail in [this link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) 

```sh
migrate -path migrations -database "postgres://user:password@host:5432/database?sslmode=disable" up
```

After installed all of components, update the config in the file config/config.yml and can run the application via command
```sh
go run cmd/main.go
```

## API Interface

### Customer create a loan
API request
```sh
curl --location 'localhost:8080/api/customers/1/loans' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 10000,
    "term": 3,
    "date": "2022-02-08"
}'
```

Response
```sh
{
    "data": {
        "loan_id": 9,
        "repayments": [
            {
                "id": "gGcVBuayDyp0OJR1",
                "loan_id": 9,
                "schedule_date": 3333.33,
                "actual_amount": 0,
                "status": 1,
                "paid_at": "0001-01-01T00:00:00Z",
                "schedule_pay_at": "2022-02-15T00:00:00Z",
                "created_at": "2023-12-29T11:03:50.035851Z",
                "updated_at": "2023-12-29T11:03:50.035851Z"
            },
            {
                "id": "5faoZkjpsdulNd6G",
                "loan_id": 9,
                "schedule_date": 3333.33,
                "actual_amount": 0,
                "status": 1,
                "paid_at": "0001-01-01T00:00:00Z",
                "schedule_pay_at": "2022-02-22T00:00:00Z",
                "created_at": "2023-12-29T11:03:50.035851Z",
                "updated_at": "2023-12-29T11:03:50.035851Z"
            },
            {
                "id": "O4fS3P5OjKj6zqpF",
                "loan_id": 9,
                "schedule_date": 3333.34,
                "actual_amount": 0,
                "status": 1,
                "paid_at": "0001-01-01T00:00:00Z",
                "schedule_pay_at": "2022-03-01T00:00:00Z",
                "created_at": "2023-12-29T11:03:50.035851Z",
                "updated_at": "2023-12-29T11:03:50.035851Z"
            }
        ]
    }
}
```

### Admin approve the loan

```sh
curl --location --request PUT 'localhost:8080/api/loans/4/approve'
```

Response
```sh
{
    "data": 4
}
```

### Customer can view loan belong to him/her

```sh
curl --location 'localhost:8080/api/customers/1/loans?page=1&size=10' \
--header 'Authorization: Bearer Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9pZCI6MSwiZXhwIjoxNzAzODU3MzAxLCJuYmYiOjE3MDM4NTM3MDEsImlhdCI6MTcwMzg1MzcwMX0.OSYBIeZDwad2WEyh4ErFChodFn08te9yu5l3Jc_f_gc'
```

Response
```sh
{
    "data": [
        {
            "id": 3,
            "amount": 10000,
            "term": 3,
            "customer_id": 1,
            "status": 2,
            "schedule_date": "2022-02-07T00:00:00Z",
            "created_at": "2023-12-28T16:00:01.825125Z",
            "updated_at": "2023-12-29T09:10:08.173804Z"
        },
        {
            "id": 5,
            "amount": 10000,
            "term": 3,
            "customer_id": 1,
            "status": 1,
            "schedule_date": "2022-02-08T00:00:00Z",
            "created_at": "2023-12-29T09:25:18.77688Z",
            "updated_at": "2023-12-29T09:25:18.77688Z"
        },
        {
            "id": 6,
            "amount": 10000,
            "term": 3,
            "customer_id": 1,
            "status": 1,
            "schedule_date": "2022-02-08T00:00:00Z",
            "created_at": "2023-12-29T09:26:53.011348Z",
            "updated_at": "2023-12-29T09:26:53.011348Z"
        },
        {
            "id": 10,
            "amount": 10000,
            "term": 3,
            "customer_id": 1,
            "status": 1,
            "schedule_date": "2022-02-08T00:00:00Z",
            "created_at": "2023-12-29T12:41:37.142264Z",
            "updated_at": "2023-12-29T12:41:37.142264Z"
        }
    ]
}
```

To check the policy of request, I decided to use the JWT token to check if the request comes from the customer or not
Can use this API to create the JWT token

```sh
curl --location 'localhost:8080/api/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "tester01",
    "password": "123456"
}'
```

Response
```sh
{
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9pZCI6MSwiZXhwIjoxNzAzODU3MzAxLCJuYmYiOjE3MDM4NTM3MDEsImlhdCI6MTcwMzg1MzcwMX0.OSYBIeZDwad2WEyh4ErFChodFn08te9yu5l3Jc_f_gc"
    }
}
```

### Customer add a repayments
```sh
curl --location --request PUT 'localhost:8080/api/repayments/gGcVBuayDyp0OJR1' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 2500
}'
```

Failed Response
```sh
{
    "error": {
        "status": 400,
        "title": "Amount should be greater or equal to the scheduled repayment"
    }
}
```

200 OK Response
```sh
{
    "data": {
        "repayment_id": "pw73kxnCToxy2CuH"
    }
}
```
