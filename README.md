# Telegram-bot-testing-kafka
<h2> 🐒 ABOUT: </h2>
This is a pet project, I had the need to code in order to learn more about kafka.
I felt the need to have automated order generator, that will send the generated orders to kafka producer, so I can work on the consumers and processing the messages received.
So here it is, it works as intended, but still require some fixes and adjustments.
I learned a lot while doing this project and hope to learn even more in the future.
The project is for study purpose only and is not meant to be used by multiple users flooding your kafka endpoint.
Project is written in **Golang**, deploy in **Docker**, **Kafka**, ZooKeeper and Kafka-UI will be deployed as well.

<h2> 🙈 INFO: </h2>
🔧 Important❗project is still in development 🔧

short: 
is a telegram bot controlled order generator, that can send generated orders to Kafka producer.

long:
This is a Telegram bot controlled microservice application.
services:
- telegram bot: controls the application 
- order service: generates random orders based on order structure.(on bot command)
The orders are sent to kafka producer and over to kafka.(on bot command)
- kafka manager: receive and forwards orders to kafka producer and then to kafka


---
<h2>🤖 DEPLOYMENT:</h2>
Goes in two deployable versions, with databases(postgres) and with in-memory storage.
Important: before deployment, update your .env files.
```
//for in-memory storage (recommended)
docker-compose -f docker-compose-nodb.yml up -d
```
```
//to deploy with two databases(not recommended)
docker-compose up -d
```
<h2> 😮 PLANNED CHANGES: </h2>

- [ ] Extract order generation logic into separate service
- [ ] make order service only store and forward orders
- [ ] complete test coverage 

