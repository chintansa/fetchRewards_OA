# Fetch Rewards API

This Go API provides functionality for processing receipts and calculating points based on specific rules. It includes two main endpoints:

1. **POST /receipts/process** - For submitting a receipt.
2. **GET /receipts/{id}/points** - For retrieving the points awarded to a processed receipt.

## Problem Statement 
[Fetch Rewards: receipt-processor-challenge](https://github.com/fetch-rewards/receipt-processor-challenge/tree/main)
## Getting Started

### Prerequisites

- Go (Golang) installed on your system. [Download Go](https://go.dev/dl/)

**All the imports used in this code are part of the Go standard library and do not require any external packages to be installed.**

### Installation and Running the Server

1. **Clone the Repository or Copy the Code**:
   - Clone this repository or copy the provided Go code into a file, say `fetchRewards.go`.

2. **Navigate to Your Project Directory**:
   - Open a terminal and navigate to the directory where your `fetchRewards.go` is located.

3. **Run the Server**:
   - Execute the command: `go run fetchRewards.go`.
   - You should see the output: `"Server starting on port 8080..."`.

### Using the API

#### **POST /receipts/process**

- **Description**: This endpoint accepts a JSON object representing a receipt and returns a unique ID for the receipt.

- **Request Format**:
  ```json
    {
        "retailer": "Walgreens",
        "purchaseDate": "2022-01-02",
        "purchaseTime": "08:13",
        "total": "2.65",
        "items": [
            {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
            {"shortDescription": "Dasani", "price": "1.40"}
        ]
    }
    ```
- **Example Request**
    ```bash
    curl -X POST http://localhost:8080/receipts/process \
     -H "Content-Type: application/json" \
     -d '{
        "retailer": "Walgreens",
        "purchaseDate": "2022-01-02",
        "purchaseTime": "08:13",
        "total": "2.65",
        "items": [
            {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
            {"shortDescription": "Dasani", "price": "1.40"}
        ]
    }'
    ```

- **Example Reponse**: 
    ```json
        {
            "id": "aed8cb35-288a-4ea5-9290-f45514dc50da"
        }
    ```

#### **GET /receipts/{id}/points**

- **Description**:Retrieves the number of points awarded for a given receipt, identified by its unique ID.


- **Example Usage**:
    ```bash
    curl http://localhost:8080/receipts/aed8cb35-288a-4ea5-9290-f45514dc50da/points
    ```

- **Response Output**:
    ```json
    {
        "points": 15
    }
    ```


## How It Works

### Receipt Processing

- The POST endpoint decodes the JSON payload into a `Receipt` struct.
- Generates a UUID for each receipt.
- Stores the receipt in an in-memory map (`receiptsMap`).

### UUID Generation

- A custom `generateUUID` function is used to create a unique identifier for each receipt.
- This utilizes the `crypto/rand` package for secure random number generation.

### Points Calculation

- The GET endpoint retrieves the receipt using its ID.
- Rules

    These rules collectively define how many points should be awarded to a receipt.

    - One point for every alphanumeric character in the retailer name.
    - 50 points if the total is a round dollar amount with no cents.
    - 25 points if the total is a multiple of 0.25.
    - 5 points for every two items on the receipt.
    - If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
    - 6 points if the day in the purchase date is odd.
    - 10 points if the time of purchase is after 2:00pm and before 4:00pm.
- `calculatePoints` function computes points based on defined rules
- The calculated points are returned in the response.

