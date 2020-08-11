import pika
import json

connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
channel = connection.channel()

channel.queue_declare(queue='queue1', durable=True)
channel.queue_declare(queue='queue2', durable=True)

payload = {
    'value': 5
}

channel.basic_publish(
    exchange='',
    routing_key='queue1',
    body=json.dumps(payload)
)

print('message sent')

def callback(ch, method, properties, body):
    print("Received message", body)

channel.basic_consume(
    queue='queue2',
    auto_ack=True,
    on_message_callback=callback
)

channel.start_consuming()
