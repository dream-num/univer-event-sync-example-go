# univer-event-sync-example-go
Go example for illustrating usage of Univer event sync.

## Quick start

1. Firstly, you should [set up Univer server](https://univer.ai/guides/sheet/server/docker#quick-start).

2. Test if RabbitMQ server is reachable.

By default, RabbitMQ server is not reachable from outside if you set up Univer server by docker compose. Modify the `univer-rabbitmq` section in docker-compose.yaml to make it reachable from outside.

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

4. Run the example.


```bash
# export RABBITMQ_URL=${THE_RABBITMQ_URL}
go run main.go
```
