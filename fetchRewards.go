package main

import (
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "math"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// Receipt and Item structs
type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

type Receipt struct {
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []Item `json:"items"`
    Total        string `json:"total"`
}

// Map to store receipts with their generated UUIDs
var receiptsMap = make(map[string]Receipt)

// Main function
func main() {
    http.HandleFunc("/receipts/process", processReceiptHandler)
    http.HandleFunc("/receipts/", getPointsHandler)

    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler for processing receipts
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var receipt Receipt
    err := json.NewDecoder(r.Body).Decode(&receipt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    id, err := generateUUID()
    if err != nil {
        http.Error(w, "Failed to generate UUID", http.StatusInternalServerError)
        return
    }

    receiptsMap[id] = receipt

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// Handler for getting points
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    id := strings.TrimPrefix(r.URL.Path, "/receipts/")
    id = strings.TrimSuffix(id, "/points")

    receipt, exists := receiptsMap[id]
    if !exists {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    points := calculatePoints(receipt)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"points": points})
}


func calculatePoints(receipt Receipt) int {
    totalPoints := 0

    // One point for every alphanumeric character in the retailer name
    re := regexp.MustCompile(`[A-Za-z0-9]`)
    totalPoints += len(re.FindAllString(receipt.Retailer, -1))

    // 50 points if the total is a round dollar amount with no cents
    if matched, _ := regexp.MatchString(`\.00$`, receipt.Total); matched {
        totalPoints += 50
    }

    // 25 points if the total is a multiple of 0.25
    totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
    if math.Mod(totalFloat*100, 25) == 0 {
        totalPoints += 25
    }

    // 5 points for every two items on the receipt
    totalPoints += (len(receipt.Items) / 2) * 5

    // Points for item descriptions and prices
    for _, item := range receipt.Items {
        trimmedDescription := strings.TrimSpace(item.ShortDescription)
        if len(trimmedDescription)%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64)
            totalPoints += int(math.Ceil(price * 0.2))
        }
    }

    // 6 points if the day in the purchase date is odd
    purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
    if purchaseDate.Day()%2 != 0 {
        totalPoints += 6
    }

    // 10 points if the time of purchase is between 2:00 PM and 4:00 PM
    purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
    if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
        totalPoints += 10
    }

    return totalPoints
}

// Function to generate a UUID
func generateUUID() (string, error) {
    uuid := make([]byte, 16)
    n, err := io.ReadFull(rand.Reader, uuid)
    if n != len(uuid) || err != nil {
        return "", err
    }
    uuid[8] = uuid[8]&^0xc0 | 0x80
    uuid[6] = uuid[6]&^0xf0 | 0x40
    return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}