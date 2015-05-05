from flask import Flask, render_template, url_for, request, redirect
from collections import Counter
from twitter import *
import json
import math
import sys

app = Flask(__name__)
app.config.from_pyfile("tweetsneak.py")

@app.route("/")
def index():
    return render_template("index.html")

@app.route("/search")
def search():
    #TODO: add token to datastore
    token = request.args.get("token")
    if token is None:
        token = ""
    
    BEARER_TOKEN = oauth2_dance(app.config["CONSUMER_KEY"], app.config["CONSUMER_SECRET"])
    t = Twitter(auth=OAuth2(bearer_token = BEARER_TOKEN))
    
    if token == "":
        tweets = []
    else:
        tweets = t.search.tweets(q=token, count=app.config["MAX_TWEETS"])["statuses"]
    
    transtab = dict((ord(char), None) for char in u"-=+|!@#$%^&*()`~[]{};:'\",<.>\\/?")
    word_list = Counter()
    
    for tweet in tweets:
        for word in tweet["text"].translate(transtab).split():
            word_list[word.lower()] += 1
            
    most_common = enumerate(word_list.most_common(10))
    num_pages = int(math.ceil(float(len(tweets)) / app.config["RPP"]))
    if num_pages == 0:
        num_pages = 1
            
    return render_template("search.html", token = token, tweets = map(json.dumps, tweets), most_common = most_common, rpp = app.config["RPP"], num_pages = num_pages)
    
@app.errorhandler(404)
def notfound(e):
    return render_template("404.html"), 404