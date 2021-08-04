import time
import sys

import redis
from flask import Flask
from flask import request

app = Flask(__name__)
cache = redis.Redis(host='redis', port=6379)

def get_hit_count():
    retries = 5
    while True:
        try:
            return cache.incr('hits')
        except redis.exceptions.ConnectionError as exc:
            if retries == 0:
                raise exc
            retries -= 1
            time.sleep(0.5)

@app.route('/')
def hello():
    count = get_hit_count()
    print('nuf', file=sys.stderr)
    for k, v in request.headers:
        print('app %s: %s' % (k, v), file=sys.stderr)
    return 'Hello World! I have been seen {} times.\n'.format(count)
