from flask import Flask, render_template, url_for, request, redirect
import sys

app = Flask(__name__)
app.config['DEBUG'] = True

@app.route("/", methods = ["GET", "POST"])
def index():
    if request.method == "POST":
        token = request.form["token"]
        return redirect(url_for("search", token=token))
    return render_template("index.html")

@app.route("/search", defaults={"token": ""})
@app.route("/search/<token>")    
def search(token):
    return render_template("search.html", token=token)
    
@app.errorhandler(404)
def notfound(e):
    return render_template("404.html"), 404