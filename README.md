# MQ

*Warning: bydlokod*

RabbitMQ, and a small helper service that invokes http endpoints on new messages in specific queues.

To run the example:

```
docker run --name rabbitmq -d -p 5672:5672 rabbitmq
go run pusher.go
uvicron service:app --reload
python source.py
```

## Config format

Pusher searches for config.yaml in current directory. You can set up triggers like this:

```
...
triggers:
- source: source-queue
  endpoint: url-to-call
  sink: sink-queue-or-nothing
- ...
- ...
```