package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"sync"

	"pricetracker/worker/internal/parser"
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

	err = ch.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		log.Fatalf("Не удалось настроить QoS: %v", err)
	}

	log.Printf("Очередь объявлена! Имя: %s, Сообщений: %d", q.Name, q.Messages)

	itemRepo := repository.NewItemRepository(db)
	citilinkParser := parser.NewCitilinkParser()
	itemService := service.NewItemService(itemRepo, citilinkParser)

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

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		for msg := range msgs {
			wg.Add(1)
			log.Printf("Получено сообщение: %s", msg.Body)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			if err := itemService.ProcessTask(ctx, msg.Body); err != nil {
				log.Printf("Ошибка декодирования Json: %v", err)
				msg.Nack(false, false)
			}else {
				msg.Ack(false)
			}
			cancel()
			wg.Done()
		}
	}()
	log.Println("Worker запущен! Ожидание сообщений...")

	<-sigs
	log.Println("Получаем сигнал завершения, воркер выключается")
	wg.Wait()
	log.Println("Все задачи завершены")
}
