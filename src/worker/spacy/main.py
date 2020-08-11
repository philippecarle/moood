#!/usr/bin/env python

from bus.connection import connection
from consumer.callback import on_message

connection = connection()
channel = connection.channel()
channel.basic_consume('entries', on_message)

try:
    channel.start_consuming()
except KeyboardInterrupt:
    channel.stop_consuming()
connection.close()
