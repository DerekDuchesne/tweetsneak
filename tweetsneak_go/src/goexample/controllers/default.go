package controllers

import (
    "goexample"
	"github.com/astaxie/beegae"
    "github.com/ChimeraCoder/anaconda"
    "io/ioutil"
    "encoding/json"
    "math"
    "appengine"
    "appengine/urlfetch"
    "appengine/datastore"
    "goexample/models"
    "time"
    "net/url"
    "strings"
    "strconv"
)

type MainController struct {
	beegae.Controller
}

func (c *MainController) Index() {
	c.TplNames = "index.html"
}

func (c *MainController) Search(){
    const(
        MaxTweets int = 1000
        Rpp int = 50
    )
    type Cred struct{
        CONSUMER_KEY string `json:"CONSUMER_KEY"`
        CONSUMER_SECRET string `json:"CONSUMER_SECRET"`
        ACCESS_TOKEN string `json:"ACCESS_TOKEN"`
        ACCESS_TOKEN_SECRET string `json:"ACCESS_TOKEN_SECRET"`
    }
    file, err := ioutil.ReadFile("credentials.json")
    if err != nil{
        panic(err)
    }
    var cred Cred
    json.Unmarshal(file, &cred)
    
    anaconda.SetConsumerKey(cred.CONSUMER_KEY)
    anaconda.SetConsumerSecret(cred.CONSUMER_SECRET)
    api := anaconda.NewTwitterApi(cred.ACCESS_TOKEN, cred.ACCESS_TOKEN_SECRET)
    ctx := appengine.NewContext(c.Ctx.Request)
    api.HttpClient.Transport = &urlfetch.Transport{Context: ctx}
    
    q := c.Input().Get("q")
    v := url.Values{}
    v.Set("count", strconv.Itoa(MaxTweets))
    var (
        tweets      []anaconda.Tweet
        word_list   map[string]int      = make(map[string]int)
        max_id      int64               = 0
    )  
    for len(tweets) < MaxTweets{
        v.Set("max_id", strconv.Itoa(int(max_id)))
        search_result, _ := api.GetSearch(q, v)
        for _, tweet := range search_result.Statuses{
            tweets = append(tweets, tweet)
            text := goexample.StripChars(tweet.Text, "-=+|!@#$%^&*()`~[]{};:'\",<.>\\/?")
            for _, word := range strings.Fields(text){
                word_list[strings.ToLower(word)] += 1
            }
        }
        if len(tweets) == 0{
            break
        }
        if max_id == tweets[len(tweets)-1].Id{
            break
        }
        max_id = tweets[len(tweets)-1].Id
    }
    
    num_pages := int(math.Ceil(float64(len(tweets)) / float64(Rpp)))
    if num_pages == 0 {
        num_pages = 1
    }
    pages := make([]int, num_pages)
    for i := 0; i < len(pages); i++{
        pages[i] = i+1
    }
    
    var (
        most_common         [10]map[string]int          
        most_common_words   []string
        most_common_counts  [10]int    
    ) 
    if len(word_list) >= 10{
        most_common_words = goexample.SortWords(word_list)[:10]
        for i, word := range most_common_words{
            most_common[i] = make(map[string]int)
            most_common[i][word] = word_list[word]
            most_common_counts[i] = word_list[word]
        }
    }
    
    most_common_json, _ := json.Marshal(most_common)
    k := datastore.NewIncompleteKey(ctx, "Query", nil)
    e := &models.Query{q, string(most_common_json), time.Now().UTC()}
    if _, err := datastore.Put(ctx, k, e); err!= nil{
        panic(err)
    }
    
    c.Data["Q"] = q
    c.Data["Tweets"] = tweets
    c.Data["MostCommonWords"] = most_common_words
    c.Data["MostCommonCounts"] = most_common_counts
    c.Data["Rpp"] = Rpp
    c.Data["Pages"] = pages
    c.TplNames = "search.html"
}
