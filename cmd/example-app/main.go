package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"yeager/pkg/tracer"
)

var t *tracer.Tracer

func main() {
	t = tracer.NewTracer("http://localhost:8080", "checkout-service")

	log.Println("Starting example shop app...")

	performCheckout(context.Background(), "cart-abc-123", 150.00)
	time.Sleep(time.Second)
	performCheckout(context.Background(), "cart-xyz-987", 9.99)
}

func performCheckout(ctx context.Context, cartID string, total float64) {
	ctx, span := t.StartSpan(ctx, "checkout")
	defer span.Finish()

	log.Printf("[%s] Trace ID: %s", cartID, span.TraceID)
	span.SetTag("cart_id", cartID)
	span.SetTag("total", fmt.Sprintf("%.2f", total))

	log.Printf("[%s] Checkout started", cartID)
	sleep(50, 100)

	chargeCreditCard(ctx)
	updateInventoryDB(ctx)
	dispatchShipping(ctx)
}

func chargeCreditCard(ctx context.Context) {
	// Child Span 1
	_, span := t.StartSpan(ctx, "charge_credit_card")
	defer span.Finish()
	span.SetTag("payment.provider", "visa")

	sleep(200, 400)
}

func updateInventoryDB(ctx context.Context) {
	// Child Span 2
	_, span := t.StartSpan(ctx, "update_inventory")
	defer span.Finish()
	span.SetTag("db.system", "postgres")

	sleep(50, 100)
}

func dispatchShipping(ctx context.Context) {
	// Child Span 3
	_, span := t.StartSpan(ctx, "dispatch_shipping")
	defer span.Finish()

	sleep(20, 50)
}

func sleep(min, max int) {
	time.Sleep(time.Duration(rand.Intn(max-min)+min) * time.Millisecond)
}
