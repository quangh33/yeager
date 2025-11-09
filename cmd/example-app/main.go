package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"yeager/pkg/tracer"
)

var (
	checkoutTracer  *tracer.Tracer
	paymentTracer   *tracer.Tracer
	inventoryTracer *tracer.Tracer
	shippingTracer  *tracer.Tracer
)

func init() {
	checkoutTracer = tracer.NewTracer("http://localhost:8080", "checkout-service")
	paymentTracer = tracer.NewTracer("http://localhost:8080", "payment-service")
	inventoryTracer = tracer.NewTracer("http://localhost:8080", "inventory-service")
	shippingTracer = tracer.NewTracer("http://localhost:8080", "shipping-service")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	// Simulate 3 concurrent user requests
	go func() {
		defer wg.Done()
		performCheckout(context.Background(), "cart-101", 99.99)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		performCheckout(context.Background(), "cart-202", 250.00)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
		performCheckout(context.Background(), "cart-303", 5.00)
	}()

	wg.Wait()
	time.Sleep(1 * time.Second)
	checkoutTracer.Close()
	paymentTracer.Close()
	inventoryTracer.Close()
	shippingTracer.Close()
}

func performCheckout(ctx context.Context, cartID string, total float64) {
	ctx, span := checkoutTracer.StartSpan(ctx, "checkout")
	defer span.Finish()
	fmt.Printf("trace id: %s\n", span.TraceID)
	span.SetTag("cart_id", cartID)
	span.SetTag("total", fmt.Sprintf("%.2f", total))
	log.Printf("[%s] Checkout started", cartID)

	sleep(50, 100)

	callPaymentService(ctx, total)
	callInventoryService(ctx, cartID)
	callShippingService(ctx)
}

func callPaymentService(ctx context.Context, amount float64) {
	_, span := paymentTracer.StartSpan(ctx, "charge_credit_card")
	defer span.Finish()

	span.SetTag("amount", fmt.Sprintf("%.2f", amount))
	log.Println("  -> Payment service charging card...")
	sleep(200, 400)
}

func callInventoryService(ctx context.Context, cartID string) {
	ctx, span := inventoryTracer.StartSpan(ctx, "update_stock")
	defer span.Finish()

	log.Println("  -> Inventory service updating stock...")
	sleep(50, 100)

	dbQuery(ctx, "a random query")
}

func dbQuery(ctx context.Context, query string) {
	_, span := inventoryTracer.StartSpan(ctx, "db:exec")
	defer span.Finish()
	span.SetTag("db.statement", query)
	sleep(10, 20)
}

func callShippingService(ctx context.Context) {
	_, span := shippingTracer.StartSpan(ctx, "schedule_dispatch")
	defer span.Finish()

	log.Println("  -> Shipping service scheduling dispatch...")
	sleep(50, 80)
}

func sleep(min, max int) {
	time.Sleep(time.Duration(rand.Intn(max-min)+min) * time.Millisecond)
}
