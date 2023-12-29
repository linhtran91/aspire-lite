# aspire-lite

## Install

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

```

### Customer can view loan belong to him

### Customer add a repayments
