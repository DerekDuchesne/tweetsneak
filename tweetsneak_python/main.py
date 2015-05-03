from flask import Flask, render_template, url_for, request, redirect
import tweepy
import sys

app = Flask(__name__)
app.config.from_pyfile("tweetsneak.py")

@app.route("/", methods = ["GET", "POST"])
def index():
    if request.method == "POST":
        token = request.form["token"]
        return redirect(url_for("search", token=token))
    return render_template("index.html")

@app.route("/search/", defaults={"token": ""})
@app.route("/search/<token>")    
def search(token):
    #TODO: add token to datastore
    #TODO: get top 10 words
    
    auth = tweepy.OAuthHandler(app.config["CONSUMER_KEY"], app.config["CONSUMER_SECRET"])
    auth.set_access_token(app.config["ACCESS_TOKEN"], app.config["ACCESS_TOKEN_SECRET"])
    api = tweepy.API(auth)
    if token == "":
        tweets = []
    else:
        tweets = api.search(q=token)
    return render_template("search.html", token = token, tweets = tweets)
    
@app.errorhandler(404)
def notfound(e):
    return render_template("404.html"), 404