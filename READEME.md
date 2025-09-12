api-rabbitmq/
├── main.go
├── go.mod
├── models/
│   └── models.go
├── services/
│   └── rabbitmq_service.go
└── handlers/
    └── message_handler.go

### RabbitMQ message model:
```json
{
    "name": "Roger",
    "zipCode": "13086656",
    "document_number": "251478526"
}
```