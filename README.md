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

Before run the exmple, put the correct RabbitMQ url to env `RABBITMQ_URL`.

Otherwise, the program will use `amqp://guest:guest@localhost:5672/` by default.

```bash
# export RABBITMQ_URL=${THE_RABBITMQ_URL}
go run main.go
```
