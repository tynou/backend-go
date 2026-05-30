package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	authServiceURL, _ := url.Parse("http://localhost:8081")
	paymentServiceURL, _ := url.Parse("http://localhost:8082")
	billingServiceURL, _ := url.Parse("http://localhost:8083")

	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)
	paymentProxy := httputil.NewSingleHostReverseProxy(paymentServiceURL)
	billingProxy := httputil.NewSingleHostReverseProxy(billingServiceURL)

	mux := http.NewServeMux()

	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", authProxy))
	mux.Handle("/api/payment/", http.StripPrefix("/api/payment", paymentProxy))
	mux.Handle("/api/billing/", http.StripPrefix("/api/billing", billingProxy))

	log.Fatal(http.ListenAndServe(":8084", mux))
}
