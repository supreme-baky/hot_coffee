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
	"strconv"
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

	mux.HandleFunc("/inventory", inventoryHandler.GetAllInventoryItems)
	mux.HandleFunc("/inventory/add", inventoryHandler.AddNewInventoryItem)
	mux.HandleFunc("/inventory/update", inventoryHandler.UpdateInventoryItem)
	mux.HandleFunc("/inventory/delete", inventoryHandler.DeleteInventoryItem)

	mux.HandleFunc("/menu", menuHandler.GetAllMenuItems)
	mux.HandleFunc("/menu/add", menuHandler.AddNewMenuItem)
	mux.HandleFunc("/menu/update", menuHandler.UpdateMenuItem)
	mux.HandleFunc("/menu/delete", menuHandler.DeleteMenuItem)

	mux.HandleFunc("/orders", orderHandler.GetAllOrders)
	mux.HandleFunc("/orders/create", orderHandler.CreateOrder)
	mux.HandleFunc("/orders/update-status", orderHandler.UpdateOrder)
	mux.HandleFunc("/orders/delete", orderHandler.DeleteOrder)

	log.Println("Starting server on :", *port)
	newPort := strconv.Itoa(*port)
	if err := http.ListenAndServe(newPort, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
