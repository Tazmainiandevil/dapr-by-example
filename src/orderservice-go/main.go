package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
)

type Order struct {
	OrderID string `json:"orderId"`
	Amount  int    `json:"amount"`
}

func main() {
	// Create Dapr client
	daprClient, err := client.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}
	defer daprClient.Close()

	http.HandleFunc("/orders", createOrderHandler(daprClient))
	http.HandleFunc("/orders/", getOrderHandler(daprClient))
	http.HandleFunc("/dapr/subscribe", subscribeHandler)
	http.HandleFunc("/healthz", healthHandler(daprClient))

	log.Println("Order service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createOrderHandler(daprClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			log.Printf("Failed to decode order: %v", err)
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		// Basic validation
		if order.OrderID == "" {
			http.Error(w, "orderId is required", http.StatusBadRequest)
			return
		}
		if order.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}

		ctx := context.Background()

		if err := saveOrder(ctx, daprClient, order); err != nil {
			log.Printf("saveOrder error: %v", err)
			http.Error(w, "failed to save order", http.StatusInternalServerError)
			return
		}

		if err := publishOrder(ctx, daprClient, order); err != nil {
			log.Printf("publishOrder error: %v", err)
			// Note: Order is already saved, so we return success but log the error
		}

		if err := storeReceipt(ctx, daprClient, order); err != nil {
			log.Printf("storeReceipt error: %v", err)
			// Non-critical error, continue
		}

		log.Printf("Order %s created successfully", order.OrderID)
		w.WriteHeader(http.StatusAccepted)
	}
}

func getOrderHandler(daprClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract order ID from URL path (e.g., /orders/123)
		orderID := r.URL.Path[len("/orders/"):]
		if orderID == "" {
			http.Error(w, "order ID is required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		order, err := getOrder(ctx, daprClient, orderID)
		if err != nil {
			log.Printf("getOrder error: %v", err)
			http.Error(w, "failed to get order", http.StatusInternalServerError)
			return
		}

		if order == nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}
}

func saveOrder(ctx context.Context, daprClient client.Client, order Order) error {
	orderData, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return daprClient.SaveState(ctx, "statestore", order.OrderID, orderData, nil)
}

func getOrder(ctx context.Context, daprClient client.Client, orderID string) (*Order, error) {
	result, err := daprClient.GetState(ctx, "statestore", orderID, nil)
	if err != nil {
		return nil, err
	}

	if result.Value == nil {
		return nil, nil // Order not found
	}

	var order Order
	if err := json.Unmarshal(result.Value, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func publishOrder(ctx context.Context, daprClient client.Client, order Order) error {
	return daprClient.PublishEvent(ctx, "pubsub", "orders", order)
}

func storeReceipt(ctx context.Context, daprClient client.Client, order Order) error {
	data := []byte("Order receipt for " + order.OrderID)
	req := &client.InvokeBindingRequest{
		Name:      "storage",
		Operation: "create",
		Data:      []byte(base64.StdEncoding.EncodeToString(data)),
		Metadata: map[string]string{
			"blobName": order.OrderID + ".txt",
			"key":      order.OrderID + ".txt",
			"fileName": order.OrderID + ".txt",
		},
	}

	_, err := daprClient.InvokeBinding(ctx, req)
	return err
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// Return empty array since this service only publishes, doesn't subscribe
	subs := []interface{}{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

func healthHandler(daprClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple health check - verify Dapr client is available
		ctx := context.Background()
		_, err := daprClient.GetMetadata(ctx)
		if err != nil {
			log.Printf("Health check failed: %v", err)
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}
}
