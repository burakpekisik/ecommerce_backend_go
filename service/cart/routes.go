package cart

import (
	"fmt"
	"net/http"

	"github.com/burakpekisik/ecommerce_backend_go/service/auth"
	"github.com/burakpekisik/ecommerce_backend_go/types"
	"github.com/burakpekisik/ecommerce_backend_go/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/orders/cancel", auth.WithJWTAuth(h.handleCancelOrder, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	// get products
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userID, cart.Address)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}

func (h *Handler) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	// Kullanıcıdan gelen JSON verisini parse et
	var cart types.CancelOrderPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Kullanıcının siparişlerini getir
	orders, err := h.userStore.GetOrdersByUserID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// OrderID'yi kontrol et
	var orderFound bool
	for _, order := range orders {
		if order.ID == cart.OrderID {
			orderFound = true
			break
		}
	}

	if !orderFound {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("couldn't find order with id %d", cart.OrderID))
		return
	}

	// Sipariş iptalini gerçekleştir
	err = h.store.CancelOrder(cart.OrderID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Güncel sipariş listesini döndür
	orders, err = h.userStore.GetOrdersByUserID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}
