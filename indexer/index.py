#!/usr/bin/env python3

from collections import defaultdict
import json
from io import StringIO
from os.path import join
import re
import sys

import elasticsearch


NUM_PARTS = re.compile(r'(\d+)')


def natural_key(s):
    '''Key function for natural sort.'''
    return [int(part) if part.isdigit() else part
            for part in NUM_PARTS.split(s)]


es = elasticsearch.Elasticsearch(*sys.argv[3:])

es.indices.delete('_all')

# Documents are lists of snippets.
es.indices.create('documents', body={
    "mappings": {
        "document": {
            "properties": {
                "sub": {
                    "type": "keyword",
                },
            },
        },
    },
})

# Snippets contain the actual text, in fields text and lemmata.
# The analyzer is "simple" because the texts are pre-tokenized.
es.indices.create('snippets', body={
    "mappings": {
        "snippet": {
            "properties": {
                "text": {
                    "type": "text",
                    "analyzer": "simple",
                },
                "lemma": {
                    "type": "text",
                    "analyzer": "simple",
                },
                "document": {
                    "type": "keyword",
                },
            },
        },
    },
})


CHUNKSIZE = 100


snippets = defaultdict(set)

print('Indexing snippets...')

textdir = sys.argv[2]
with open(sys.argv[1]) as idfile:
    data = StringIO()

    n = 0
    for ident in idfile:
        ident = ident.strip()
        json.dump({'index': {'_id': ident}}, data)
        data.write('\n')

        doc, part = ident.split('_paragraph_', 1)
        snippets[doc].add(ident)

        text = open(join(textdir, '%s/paragraph_%s.txt' % (doc, part))).read()

        json.dump({'text': text, 'document': doc}, data)
        data.write('\n')

        n += 1
        if n >= CHUNKSIZE:
            es.bulk(index='snippets', doc_type='snippet', body=data.getvalue())
            data = StringIO()
            n = 0

    if n > 0:
        es.bulk(index='snippets', doc_type='snippet', body=data.getvalue())

print('Indexing documents...')
data = StringIO()
for doc, snippetset in snippets.items():
    json.dump({'index': {'_id': doc}}, data)
    data.write('\n')
    json.dump({'sub': sorted(snippetset, key=natural_key)}, data)
    data.write('\n')

es.bulk(index='documents', doc_type='document', body=data.getvalue())

print('Indexing done.')
