package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	amqp "github.com/rabbitmq/amqp091-go"

	"pricetracker/api/internal/handler"
	"pricetracker/api/internal/middleware"
	"pricetracker/api/internal/repository"
	"pricetracker/api/internal/service"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://admin:qwerty@localhost:5432/pricetracker"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("База не отвечает: %v", err)
	}
	log.Println("Успешное подключение к бд")

	rmqUrl := os.Getenv("RMQ_URL")
	if rmqUrl == "" {
		rmqUrl = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rmqUrl)
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMq: %v", err)
	}
	defer conn.Close()
	log.Println("Успешное подключение к RabbitMq!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %v", err)
	}
	defer ch.Close()
	log.Println("Успешно подключились к каналу!")

	q, err := ch.QueueDeclare(
		"parsing_tasks",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %v", err)
	}

	log.Printf("Очередь объявлена! Имя: %s, Сообщений: %d", q.Name, q.Messages)

	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo, ch)
	itemHandler := handler.NewItemHandler(itemService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET/items", itemHandler.GetAll)
	mux.HandleFunc("POST/track", itemHandler.Track)
	mux.HandleFunc("GET /history/{id}", itemHandler.GetHistory)

	handlerWithCORS := middleware.CORS(mux)
	
	srv := &http.Server{
		Addr: ":8080",
		Handler: handlerWithCORS,
	}

	go func() {
	log.Println("сервер запущен на http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Ошибка сервера: %v", err)
	}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Получен сигнал на завершение, останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при выкл сервера: %v", err)
	}
	


	log.Println("HTTP сервер остановлен, зыкрываем соединение с Бд и Rabbit")
}
