import json

import pika
import spacy
from vaderSentiment import vaderSentiment


def on_message(channel: pika.spec.Channel,
               method: pika.spec.Basic.Deliver,
               properties: pika.spec.BasicProperties,
               body: bytes):
    nlp = spacy.load("en_core_web_sm")
    doc = nlp(body.decode("utf-8"))
    sentences = [str(s) for s in doc.sents]
    analyzer = vaderSentiment.SentimentIntensityAnalyzer()
    data = {'sentences': []}
    for s in sentences:
        sentence = str(s)
        data['sentences'].append({
            'sentence': sentence,
            'score': analyzer.polarity_scores(sentence)
        })
        sentence = nlp(body.decode("utf-8"))

    print(json.dumps(data))
