from flask import Flask
app = Flask(__name__)

@app.route('/')
def hello_world():
    print("incoming request")
    return 'Hello, World from Flask!\n'