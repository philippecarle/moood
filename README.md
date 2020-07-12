# Mood (temporary name)

## Main principles

The idea is to develop an app on which the user can simply write text, ideally on a daily basis. Some random thoughts, some feelings, but something meaningful. 
The main Go API would receive each paragraph (aka entry) as it's being written. 
The entry is persisted by the API in a MongoDB storage. A message is sent to RabbitMQ. A Python worker using spacy is in charge of extracting the sentences, the words, extracting the sentiments and generating another message is sent to the queue. Another worker, the store, will process those messages to enrich the initial entry with the analysis results. 

## Example

A POST HTTP request is sent to the API

```json
{
    "content": "He took a sip of the drink. He wasn't sure whether he liked it or not, but at this moment it didn't matter. She had made it especially for him so he would have forced it down even if he had absolutely hated it. That's simply the way things worked. She made him a new-fangled drink each day and he took a sip of it and smiled, saying it was excellent."
}
```

The document is stored and the response is returned:

```json
{
    "id": "5f0b97388a1afee2a6d6150b",
    "content": "He took a sip of the drink. He wasn't sure whether he liked it or not, but at this moment it didn't matter. She had made it especially for him so he would have forced it down even if he had absolutely hated it. That's simply the way things worked. She made him a new-fangled drink each day and he took a sip of it and smiled, saying it was excellent."
}
```

The vaspacyder worker processed the text and created this metadata which should now be attached to the document:

```json
{
    "sentences": [
        {
            "sentence": "He took a sip of the drink.",
            "score": {
                "neg": 0.0,
                "neu": 1.0,
                "pos": 0.0,
                "compound": 0.0
            }
        },
        {
            "sentence": "He wasn't sure whether he liked it or not, but at this moment it didn't matter.",
            "score": {
                "neg": 0.148,
                "neu": 0.743,
                "pos": 0.109,
                "compound": 0.0793
            }
        },
        {
            "sentence": "She had made it especially for him so he would have forced it down even if he had absolutely hated it.",
            "score": {
                "neg": 0.283,
                "neu": 0.717,
                "pos": 0.0,
                "compound": -0.8173
            }
        },
        {
            "sentence": "That's simply the way things worked.",
            "score": {
                "neg": 0.0,
                "neu": 1.0,
                "pos": 0.0,
                "compound": 0.0
            }
        },
        {
            "sentence": "She made him a new-fangled drink each day and he took a sip of it and smiled, saying it was excellent.",
            "score": {
                "neg": 0.0,
                "neu": 0.725,
                "pos": 0.275,
                "compound": 0.802
            }
        }
    ]
}
```

## Architecture Diagram

![Architecture Diagram](doc/architecture.png)