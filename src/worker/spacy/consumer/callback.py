import json

import pika.channel
import spacy
from vaderSentiment import vaderSentiment


def on_message(channel: pika.channel.Channel,
               method: pika.spec.Basic.Deliver,
               properties: pika.spec.BasicProperties,
               body: bytes):
    nlp = spacy.load("en_core_web_sm")
    entry = json.loads(body.decode("utf-8"))
    
    doc = nlp(entry['content'])
    sentences = [str(s) for s in doc.sents]

    analyzer = vaderSentiment.SentimentIntensityAnalyzer()
    entry['score'] = analyzer.polarity_scores(entry['content'])

    entry['sentences'] = []
    for s in sentences:
        entry['sentences'].append({
            'sentence': s,
            'score': analyzer.polarity_scores(s)
        })

    channel.basic_publish('', 'processed', bytearray(json.dumps(entry), encoding='utf8'))
    print("Processed entry " + entry['id'])
