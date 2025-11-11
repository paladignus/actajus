// Package main
package main

func main() {
	// 	cfg := config.Load()
	// 	nc, js, err := nats.ConnectAndSetup(cfg.NATS)
	// 	if err != nil {
	// 		log.Fatal("NATS setup failed:", err)
	// 	}
	// 	defer nc.Close()
	//
	// 	subscriber := nats.NewSubscriber(js, cfg.NATS.DLQSubject)
	//
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	defer cancel()
	//
	// 	// Handler simples para debug
	// 	handler := func(data []byte) error {
	// 		log.Printf("📧 EMAIL HANDLER RECEIVED: %s", string(data))
	// 		return nil
	// 	}
	//
	// 	// Nome único para evitar conflito
	// 	durableName := fmt.Sprintf("debug-consumer-%d", time.Now().Unix())
	//
	// 	log.Printf("🔧 Subscribing with durable: %s", durableName)
	//
	// 	err = subscriber.Subscribe(ctx, "user.created", handler)
	// 	if err != nil {
	// 		log.Fatal("❌ Subscribe failed:", err)
	// 	}
	//
	// 	log.Println("✅ Consumer ready. Waiting for messages...")
	// 	// Manter vivo
	// 	<-ctx.Done()
	// }

	// config := config.Load()
	// nc, js, err := nats.ConnectAndSetup(config.NATS)
	// if err != nil {
	// 	log.Fatal("NATS connection failed:", err)
	// }
	// defer nc.Close()
	// subscriber := nats.NewSubscriber(js, config.NATS.DLQSubject)
	// if err != nil {
	// 	log.Fatal("NATS connection failed:", err)
	// }
	//
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// smtp := adapter.NewSMTPEmail(config.SMTP)
	//
	// err = subscriber.Subscribe(
	// 	ctx,
	// 	"user.created",
	// 	"email-service",
	// 	func(msg []byte) error {
	// 		event := event.RecoveredPassword{}
	// 		evt, err := event.Deserialize(msg)
	// 		if err != nil {
	// 			log.Printf("❌ Failed to deserialize event: %v", err)
	// 		}
	// 		log.Printf("📧 EMAIL HANDLER RECEIVED: %s", evt.Email)
	// 		return smtp.SendEmail(ctx, evt.Email, evt.URL)
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal("Failed to subscribe:", err)
	// }
	//
	// log.Println("📬 NATS consumer ready. Listening to 'user.created'...")
	//
	// // Esperar sinal para encerrar
	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// <-sig
	// log.Println("Shutting down...")
}
