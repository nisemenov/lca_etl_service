package main

// Возможно, нужно будет вот так как-то делать, чтобы нормально протестировать именно сборку с правильными
// исходными; в частности, что создается у PaymentProducer правильный HTTPClient,
// со всеми WithHeaders, WithMiddleware, etc.
//
//
// type App struct {
//     PaymentProducer producer.PaymentProducer
// }
//
// func BuildApp(cfg *config.Config, logger *slog.Logger) *App {
//     rawHTTP := &http.Client{}
//
//     paymentsHTTP := httpclient.New(
//         cfg.APIBaseURL,
//         rawHTTP,
//         httpclient.WithHeaders(map[string]string{
//             "X-Internal-Token": cfg.InternalToken,
//         }),
//     )
//
//     paymentProducer := producer.NewPaymentProducer(
//         paymentsHTTP,
//         logger,
//     )
//
//     return &App{
//         PaymentProducer: paymentProducer,
//     }
// }
//
// func TestBuildApp_ConfiguresPaymentHTTPClient(t *testing.T) {
//     var got *http.Request
//
//     ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         got = r
//         w.WriteHeader(200)
//     }))
//     defer ts.Close()
//
//     cfg := &config.Config{
//         APIBaseURL:   ts.URL,
//         InternalToken: "secret123",
//     }
//
//     logger := slog.New(slog.NewTextHandler(io.Discard, nil))
//
//     app := BuildApp(cfg, logger)
//
//     err := app.PaymentProducer.SendPayment(context.Background(), Payment{ID: "1"})
//     require.NoError(t, err)
//
//     require.Equal(t, "secret123", got.Header.Get("X-Internal-Token"))
// }
