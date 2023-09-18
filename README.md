# receipt-processor-api
API that processes receipts and calculates rewarded points according to the established rules. This is done through 2 endpoints:

## Process receipt
This endpoint will receive receipt information in a json format and will generate and return a UUID, saving the receipt data in-memory.

`POST /recept-processor/v1/receipts/process`

### Example request
`http://localhost:8000/receipt-processor/v1/receipts/process`

```json
{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
}
```

### Example response
```json
{
    "id": "47c5b092-5f6e-41b3-a622-f91c383d76aa"
}
```

## Get points
This endpoint will receive the generated UUID and will look for the corresponding receipt data to calculate the corresponding points according to the following rules:

- One point for every alphanumeric character in the retailer name.
- 50 points if the total is a round dollar amount with no cents.
- 25 points if the total is a multiple of 0.25.
- 5 points for every two items on the receipt.
- If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
- 6 points if the day in the purchase date is odd.
- 10 points if the time of purchase is after 2:00pm and before 4:00pm.

This will return the number of points.

`GET /recept-processor/v1/receipts/{id}/points`

### Example request
`http://localhost:8000/receipt-processor/v1/receipts/47c5b092-5f6e-41b3-a622-f91c383d76aa/points`

### Example response
```json
{
    "points": 31
}
```

## How to run the service

The most useful make targets for working locally are:

* `make build`: Builds the service.
* `make run`: Starts the service locally running on port `8000`.
* `make clean`: Clean temporary files.

## Dependencies
This project has dependencies on:
* go (`1.18.10`)