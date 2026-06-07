package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pricetracker/worker/internal/repository"
	"pricetracker/worker/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	amqp "github.com/rabbitmq/amqp091-go"
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
		log.Fatalf("БД не отвечает: %v", err)
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
		log.Fatalf("Не удалось открыть канал, %v", err)
	}
	defer ch.Close()
	log.Printf("Успешно подключились к каналу!")

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
	itemService := service.NewItemService(itemRepo)

	msgs, err := ch.Consume(
		"parsing_tasks",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка регистрации консумера: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		for msg := range msgs {
			log.Printf("Получено сообщение: %s", msg.Body)
			if err := itemService.ProcessTask(msg.Body); err != nil {
				log.Printf("Ошибка декодирования Json: %v", err)
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}()
	log.Println("Worker запущен! Ожидание сообщений...")

	<-sigs
	log.Println("Получаем сигнал завершения, воркер выключается")
}
