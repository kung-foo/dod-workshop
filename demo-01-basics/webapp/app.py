#!/usr/bin/env python

import os
from flask import Flask, Response
import netifaces as ni

app = Flask(__name__)

container_addr = ni.ifaddresses('eth0')[ni.AF_INET][0]['addr']

@app.route('/')
def index():
    return 'Hello World!'

@app.route('/ping')
def ping():
    return 'OK'

@app.route('/hosts')
def hosts():
    hosts = []
    with open('/etc/hosts', 'rb') as f:
        for line in f:
            if 'ip6' not in line and 'localhost' not in line:
                hosts.append(line.strip())
    return Response(
        '\n'.join(hosts),
        mimetype='text/plain'
    )

@app.route('/env')
def env():
    exclude = set(["PATH", "TERM", "HOSTNAME", "HOME", "WERKZEUG_RUN_MAIN"])
    return Response(
        "\n".join(['{}={}'.format(k, v) for k, v in os.environ.items() if k not in exclude]),
        mimetype='text/plain')

if __name__ == '__main__':
    app.run(host=container_addr, debug=True)
