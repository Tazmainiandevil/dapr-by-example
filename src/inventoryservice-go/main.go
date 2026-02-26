package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprhttp "github.com/dapr/go-sdk/service/http"
)

type Order struct {
	OrderID string `json:"orderId"`
	Amount  int    `json:"amount"`
}

func main() {
	s := daprhttp.NewService(":8081")

	// Add topic subscription using the proper Dapr Go SDK approach
	subscription := &common.Subscription{
		PubsubName: "pubsub",
		Topic:      "orders",
		Route:      "/orders",
	}

	if err := s.AddTopicEventHandler(subscription, ordersHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// Add health check endpoint
	http.HandleFunc("/healthz", healthHandler())

	log.Println("Inventory service listening on :8081")
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting server: %v", err)
	}
}

func ordersHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	// The Dapr Go SDK provides e.Data which can be []byte or map[string]interface{}
	// depending on the content type and how Dapr processes it
	var order Order

	switch data := e.Data.(type) {
	case []byte:
		// Data is raw bytes, unmarshal directly
		if err := json.Unmarshal(data, &order); err != nil {
			log.Printf("error unmarshaling order from bytes: %v", err)
			return false, err
		}
	case map[string]interface{}:
		// Data is already parsed as a map, marshal back to JSON then unmarshal to struct
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("error marshaling map data: %v", err)
			return false, err
		}
		if err := json.Unmarshal(dataBytes, &order); err != nil {
			log.Printf("error unmarshaling order from map: %v", err)
			return false, err
		}
	default:
		log.Printf("error: unexpected data type %T", e.Data)
		return false, nil
	}

	// Basic validation
	if order.OrderID == "" {
		log.Printf("error: received order with empty OrderID")
		return false, nil
	}
	if order.Amount <= 0 {
		log.Printf("error: received order %s with invalid amount: %d", order.OrderID, order.Amount)
		return false, nil
	}

	log.Printf("Inventory service received order: %s (amount: %d)", order.OrderID, order.Amount)

	log.Printf("[INVENTORY] 📦 Received order: %s (amount: $%d) - %s", order.OrderID, order.Amount, time.Now().Format("2006-01-02 15:04:05"))
	log.Printf("[INVENTORY] ✅ Order %s processed successfully", order.OrderID)

	return false, nil
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple health check - verify Dapr client is available
		ctx := context.Background()
		daprClient, err := client.NewClient()
		if err != nil {
			http.Error(w, "dapr unavailable", http.StatusServiceUnavailable)
			return
		}
		defer daprClient.Close()

		if _, err := daprClient.GetMetadata(ctx); err != nil {
			http.Error(w, "dapr unavailable", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	}
}
