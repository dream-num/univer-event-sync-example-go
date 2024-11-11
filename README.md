# univer-event-sync-example-go
Go example for illustrating usage of Univer event sync.

## Quick start

1. Firstly, you should [set up Univer server](https://univer.ai/guides/sheet/server/docker#quick-start).

2. Test if RabbitMQ server is reachable.
> If you have configured RabbitMQ server by your own, please processed to next section.

By default, RabbitMQ server is not reachable from outside if you set up Univer server by docker compose because port 5672 is not mapped to the host machine.

Modify the `univer-rabbitmq` section in docker-compose.yaml to make it reachable from outside.

For example, add the following lines to the `univer-rabbitmq` section:

```yaml
univer-rabbitmq:
  ports:
    - 5672:5672
```

3. Enable Data Sync feature in Univer server.

By default, Data Sync feature is not enabled in Univer server.

Edit .env file, set `EVENT_SYNC` to `true`.
```
EVENT_SYNC=true
```

Now you should restart the servers to take effect:
```bash
bash run.sh
```

4. Run the example.

Before run the example, put the correct RabbitMQ url to env `RABBITMQ_URL`.

Otherwise, the program will use `amqp://guest:guest@localhost:5672/` by default.

Basic Consumer

This is a basic RabbitMQ consumer example that listens to the univer-event-sync.changeset topic and prints the events to the console.

To run the basic consumer:
```bash
# export RABBITMQ_URL=${THE_RABBITMQ_URL}
cd basic-consumer
go run main.go
```

Persistent Consumer

This is a RabbitMQ consumer example with message persistence. The queue is set with a maximum length of 100000 messages.

To run the persistent consumer:
```bash
# export RABBITMQ_URL=${THE_RABBITMQ_URL}
cd persistent-consumer
go run main.go
```

Persistent multi-queue Consumer

This is a RabbitMQ example with multiple consumers listening to different queues. Each queue can consume the full amount of data in the exchange.
At the same time, it also supports message persistence.

To run the multi-queue consumer:
```bash
# export RABBITMQ_URL=${THE_RABBITMQ_URL}
cd persistent-multi-queue-consumer
go run main.go
```