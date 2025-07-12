package main

import (
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"net/http"
)

func main() {
	inventoryService := service.NewInventoryService()
	menuService := service.NewMenuService()
	orderService := service.NewOrderService()

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
	mux.HandleFunc("/orders/update-status", orderHandler.UpdateOrderStatus)
	mux.HandleFunc("/orders/delete", orderHandler.DeleteOrder)

	http.ListenAndServe(":8080", mux)
}
