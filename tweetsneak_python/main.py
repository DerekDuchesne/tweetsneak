from flask import Flask, render_template, url_for, request, redirect
from collections import Counter
from twitter import *
import json, math, sys, datetime
from google.appengine.ext import db

class Query(db.Model):
    query = db.StringProperty()
    result = db.StringProperty(required = True)
    timestamp = db.DateTimeProperty(required = True)

app = Flask(__name__)
app.config.from_pyfile("tweetsneak.py")

@app.route("/")
def index():
    return render_template("index.html")

@app.route("/search")
def search():
    q = request.args.get("q")
    if q is None:
        q = ""
    
    BEARER_TOKEN = oauth2_dance(app.config["CONSUMER_KEY"], app.config["CONSUMER_SECRET"])
    t = Twitter(auth=OAuth2(bearer_token = BEARER_TOKEN))
    
    if q == "":
        tweets = []
    else:
        tweets = t.search.tweets(q=q, count=app.config["MAX_TWEETS"])["statuses"]
    
    transtab = dict((ord(char), None) for char in u"-=+|!@#$%^&*()`~[]{};:'\",<.>\\/?") #trans table for removing punctuation
    word_list = Counter()
    
    for tweet in tweets:
        for word in tweet["text"].translate(transtab).split():
            word_list[word.lower()] += 1
    
    most_common = word_list.most_common(10)
    most_common_json = json.dumps(most_common)
        
    entity = Query(query=q, result=most_common_json, timestamp=datetime.datetime.utcnow())
    entity.put()
        
    most_common = enumerate(most_common)
    num_pages = int(math.ceil(float(len(tweets)) / app.config["RPP"]))
    if num_pages == 0:
        num_pages = 1
            
    return render_template("search.html", q = q, tweets = map(json.dumps, tweets), most_common = most_common, rpp = app.config["RPP"], num_pages = num_pages)
    
@app.errorhandler(404)
def notfound(e):
    return render_template("404.html"), 404