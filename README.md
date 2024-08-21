# Telegram-bot-testing-kafka

is a telegram bot controlled order generator, that can send generated orders to Kafka producer.

Bot goes in two deployable versions, with databases(postgres) and with in-memory storage.

#### deployment:
`docker-compose -f docker-compose-nodb.yml up -d` for in-memory
`docker-compose up -d` to deploy with two databases(not recommended)

