package api

import (
	"net/http"
	"strconv"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{db: db}
}

type PublicOrder struct {
	ID              uint              `json:"id"`
	TotalAmount     float64           `json:"total_amount"`
	Status          string            `json:"status"`
	ShippingAddress string            `json:"shipping_address"`
	UserID          uint              `json:"user_id"`
	OrderItems      []PublicOrderItem `json:"order_items"`
}

type PublicOrderItem struct {
	ID        uint    `json:"id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	ProductID uint    `json:"product_id"`
}

// Additional methods for OrderController can be added here.

func (c *OrderController) handleGetOrders(ctx *gin.Context) {
	var orders []models.Order

	result := c.db.Preload("OrderItems").Find(&orders)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	publicOrders := make([]PublicOrder, len(orders))
	for i, order := range orders {
		publicOrders[i] = PublicOrder{
			ID:              order.ID,
			TotalAmount:     order.TotalAmount,
			Status:          order.Status,
			ShippingAddress: order.ShippingAddress,
			UserID:          order.UserID,
			OrderItems:      make([]PublicOrderItem, len(order.OrderItems)),
		}
		for j, item := range order.OrderItems {
			publicOrders[i].OrderItems[j] = PublicOrderItem{
				ID:        item.ID,
				Quantity:  item.Quantity,
				Price:     item.Price,
				ProductID: item.ProductID,
			}
		}
	}

	ctx.JSON(http.StatusOK, publicOrders)
}

func (c *OrderController) handleGetOrder(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	result := c.db.Preload("OrderItems").First(&order, orderID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	publicOrder := PublicOrder{
		ID:              order.ID,
		TotalAmount:     order.TotalAmount,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		UserID:          order.UserID,
		OrderItems:      make([]PublicOrderItem, len(order.OrderItems)),
	}
	for i, item := range order.OrderItems {
		publicOrder.OrderItems[i] = PublicOrderItem{
			ID:        item.ID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			ProductID: item.ProductID,
		}
	}
	ctx.JSON(http.StatusOK, publicOrder)
}

func (c *OrderController) handleCreateOrder(ctx *gin.Context) {
	var order models.OrderPayload
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order payload"})
		return
	}

	newOrder := models.Order{
		TotalAmount:     order.TotalAmount,
		Status:          "cart", // Default status
		ShippingAddress: order.ShippingAddress,
		UserID:          order.UserID,
	}

	for _, item := range order.OrderItems {
		newOrder.OrderItems = append(newOrder.OrderItems, models.OrderItem{
			Quantity:  item.Quantity,
			Price:     item.Price,
			ProductID: item.ProductID,
		})
	}

	if err := c.db.Create(&newOrder).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order_id": newOrder.ID})
}

func (c *OrderController) handleUpdateOrder(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var payload models.UpdateOrderPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update payload"})
		return
	}

	var order models.Order
	if err := c.db.First(&order, orderID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = payload.Status
	order.ShippingAddress = payload.ShippingAddress

	if err := c.db.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

func (c *OrderController) handleUpdateOrderItem(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var payload models.UpdateOrderItemPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update payload"})
		return
	}

	var order models.Order
	if err := c.db.First(&order, orderID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var item models.OrderItem
	if err := c.db.First(&item, itemID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	if item.OrderID != order.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item does not belong to this order"})
		return
	}

	item.Quantity = payload.Quantity

	if err := c.db.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order item"})
		return
	}

	order.RefreshPrice()
	c.db.Save(&order)

	ctx.JSON(http.StatusOK, gin.H{"message": "Order item updated successfully"})
}

func (c *OrderController) handleAddItem(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var payload models.OrderItemPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item payload"})
		return
	}

	var order models.Order
	if result := c.db.First(&order, orderID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	newItem := models.OrderItem{
		Quantity:  payload.Quantity,
		Price:     payload.Price,
		ProductID: payload.ProductID,
		OrderID:   order.ID,
	}

	if err := c.db.Create(&newItem).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to order"})
		return
	}

	order.RefreshPrice()
	if err := c.db.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order total amount"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Item added to order successfully", "item_id": newItem.ID})
}

func (c *OrderController) handleRemoveItem(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var order models.Order
	if result := c.db.First(&order, orderID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var item models.OrderItem
	if result := c.db.First(&item, itemID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	if item.OrderID != order.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item does not belong to this order"})
		return
	}

	if err := c.db.Delete(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from order"})
		return
	}

	order.RefreshPrice()
	if err := c.db.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order total amount"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item removed from order successfully"})
}

func (c *OrderController) handleDeleteOrder(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if result := c.db.First(&order, orderID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := c.db.Delete(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
