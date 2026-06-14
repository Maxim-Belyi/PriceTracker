package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	amqp "github.com/rabbitmq/amqp091-go"

	"pricetracker/api/internal/repository"
	"pricetracker/api/internal/service"
	"pricetracker/api/internal/handler"
	"pricetracker/api/internal/middleware"
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
	
	mux.HandleFunc("/items", itemHandler.GetAll)
	mux.HandleFunc("/track", itemHandler.Track)

	handlerWithCORS := middleware.CORS(mux)
	log.Println("сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",  handlerWithCORS))
}
