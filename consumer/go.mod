module github.com/go-rabbitmq-sample/consumer

go 1.19

replace github.com/go-rabbitmq-sample/shared => ../shared

require github.com/go-rabbitmq-sample/shared v0.0.0-00010101000000-000000000000

require github.com/rabbitmq/amqp091-go v1.5.0 // indirect
