package main

import (
	"flag"
	"fmt"
	"hot-coffee/help"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"log"
	"net/http"
	"strings"
)

func main() {
	helpFlag := flag.Bool("help", false, "Prints help information")
	port := flag.Int("port", 8080, "Port number for the server")
	flag.Parse()

	if *helpFlag {
		help.PrintInfo()
		return
	}

	inventoryRepo := dal.NewJSONInventoryManager("data/inventory.json")
	menuRepo := dal.NewJSONMenuManager("data/menu_items.json")
	orderRepo := dal.NewJSONOrderManager("data/orders.json")

	inventoryService := service.NewInventoryService(inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo)

	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)

	mux := http.NewServeMux()

	mux.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			inventoryHandler.GetAllInventoryItems(w, r)
		case http.MethodPost:
			inventoryHandler.AddNewInventoryItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/inventory/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		if id == "" {
			http.Error(w, "Missing inventory ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			inventoryHandler.GetInventoryItem(w, r, id)
		case http.MethodPut:
			inventoryHandler.UpdateInventoryItem(w, r, id)
		case http.MethodDelete:
			inventoryHandler.DeleteInventoryItem(w, r, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			menuHandler.GetAllMenuItems(w, r)
		case http.MethodPost:
			menuHandler.AddNewMenuItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/menu/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		if id == "" {
			http.Error(w, "Missing menu item ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			menuHandler.GetMenuItem(w, r, id)
		case http.MethodPut:
			menuHandler.UpdateMenuItem(w, r, id)
		case http.MethodDelete:
			menuHandler.DeleteMenuItem(w, r, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetAllOrders(w, r)
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		if strings.HasSuffix(path, "/close") && r.Method == http.MethodPost {
			orderID := strings.TrimSuffix(path, "/close")
			orderID = strings.TrimSuffix(orderID, "/")
			orderHandler.CloseOrder(w, r, orderID)
			return
		}

		orderID := strings.TrimSuffix(path, "/")
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetOrderByID(w, r, orderID)
		case http.MethodPut:
			orderHandler.UpdateOrder(w, r, orderID)
		case http.MethodDelete:
			orderHandler.DeleteOrder(w, r, orderID)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	if *port < 1 || *port > 65535 {
		log.Fatalf("Invalid port number: %d. Must be between 1 and 65535.", *port)
	}
	address := fmt.Sprintf(":%d", *port)
	log.Printf("Server started at http://localhost%s\n", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
