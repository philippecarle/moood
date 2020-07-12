import pika


def connection():
    credentials = pika.PlainCredentials('guest', 'guest')
    return pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', credentials=credentials))
