package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly"
)

const base = "https://api.tumblr.com"

func StartScraping(run <-chan bool, Mdb MongoDb, blog string, continueLastRun bool, progressBar *widget.ProgressBar) {

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:96.0) Gecko/20100101 Firefox/96.0")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Referer", "https://badbitchjuice.tumblr.com/")
		r.Headers.Set("Content-Type", "text/plain;charset=UTF-8")
		r.Headers.Set("Origin", "https://badbitchjuice.tumblr.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Sec-Fetch-Dest", "empty")
		r.Headers.Set("Sec-Fetch-Mode", "cors")
		r.Headers.Set("Sec-Fetch-Site", "cross-site")
		r.Headers.Set("AlexaToolbar-ALX_NS_PH", "AlexaToolbar/alx-4.0.2")
		r.Headers.Set("TE", "trailers")
		r.Headers.Set("Authorization", "Bearer aIcXSOoTtqrzR8L8YEIOmBeW94c3FmbSNSWAUbxsny9KKx5VFh")
		r.Headers.Set("Cookie", "language=%2Cfr_FR; tmgioct=61f9402270ccb00684737330; _ga=GA1.1.93384512.1639324347; _ga_KPWKWLWW25=GS1.1.1641207892.1.0.1641208053.0; logged_in=1; pfg=1a38a0851adf19a8d982526e17a1d3e36e6d6cc6cef19bd4d19b0a4d6b2d97af%23%7B%22gdpr_is_acceptable_age%22%3A1%2C%22exp%22%3A1675785312%2C%22vc%22%3A%22%22%7D%234011591233")
	})
	c.OnResponse(func(r *colly.Response) {
		select {
		case i := <-run:
			if !i {
				log.Println("stopped")
				return
			}
		default:
			var resData ResponseData
			json.Unmarshal(r.Body, &resData)
			//posts := resData.Response.Posts
			log.Println(resData.Response.Blog.TotalPosts)
			Mdb.SetBlogInfo(blog, bloginfo{TotalPosts: resData.Response.Blog.TotalPosts})
			//
			postsOnPage := ParsingTumblrPostResult(resData)
			step := float64(1) / float64(len(postsOnPage))
			var StepNow float64
			for _, post := range postsOnPage {
				StepNow = step + StepNow
				progressBar.SetValue(StepNow)
				if len(post.Id) == 0 {

					log.Println(post)
				}
				Mdb.AddPost(blog, post)

			}
			time.Sleep(5 * time.Second)

			nextLinks := resData.Response.Links.Next
			if len(nextLinks.Href) != 0 {
				log.Println("next")
				Mdb.setLastPage(blog, base+nextLinks.Href)
				progressBar.SetValue(0)

				c.Visit(base + nextLinks.Href)
			} else {

			}
		}

	})
	c.OnError(func(r *colly.Response, e error) {
		log.Println(e)
	})
	if continueLastRun {

		nextUrl := Mdb.GetLastPage(blog)
		if len(nextUrl) != 0 {
			log.Println("Continuing page")

			c.Visit(nextUrl)

		} else {
			page := fmt.Sprintf("https://api.tumblr.com/v2/blog/%v/posts?fields[blogs]=name,avatar,title,url,is_adult,?is_member,description_npf,uuid,can_be_followed,?followed,?advertiser_name,theme,?primary,?is_paywall_on,?paywall_access,?subscription_plan,share_likes,share_following,can_subscribe,subscribed,ask,?can_submit,?is_blocked_from_primary,?is_blogless_advertiser,?tweet,updated,first_post_timestamp,posts,description,?top_tags_all&npf=true&reblog_info=true&type=photo", blog)
			c.Visit(page)
		}
	}
}
func ParsingTumblrPostResult(resData ResponseData) []PostMeta {
	var Posts []PostMeta
	for _, post := range resData.Response.Posts {
		if post.Type == "post" {
			id, timestamp := post.Id, post.Timestamp
			var medias []PostMedia
			var Pcontent content
			if len(post.RebloggedFromName) != 0 {
				Pcontent = post.Trail[0].Content
			} else {
				Pcontent = post.Content
			}
			for _, media := range Pcontent {
				mediaType := media.Type
				filterBool := mediaType == "image" || mediaType == "video"
				if len(media.Media) != 0 && filterBool {
					mediaRessourceType := media.Media[0].Type
					mediaUrl := media.Media[0].Url
					medias = append(medias, PostMedia{Type: mediaRessourceType, Url: mediaUrl})

				}
			}
			Posts = append(Posts, PostMeta{Id: id, Timestamp: timestamp, Medias: medias, Type: "post"})
		}
	}
	return Posts
}

type PostMedia struct {
	Type string
	Url  string
}
type PostMeta struct {
	Type      string      `bson:"type"`
	Id        string      `bson:"id"`
	Timestamp int         `bson:"timestamp"`
	Medias    []PostMedia `bson:"medias"`
}
type ResponseData struct {
	Response struct {
		Blog struct {
			TotalPosts int `json:"posts"`
		} `json:"blog"`

		Posts []Post `json:"posts"`
		Links struct {
			Next struct {
				Href   string `json:"href"`
				Params struct {
					PageNumber string `json:"page_number"`
				} `json:"query_params"`
			} `json:"next"`
		} `json:"_links"`
	} `json:"response"`
}
type Post struct {
	Type              string `json:"object_type"`
	Isnsfw            bool   `json:"isnsfw"`
	Id                string `json:"id"`
	Timestamp         int    `json:"timestamp"`
	RebloggedFromName string `json:"reblogged_from_name"`
	Trail             []struct {
		Content content `json:"content"`
	} `json:"trail"`
	Content content `json:"content"`
}
type content []struct {
	Type  string `json:"type"`
	Media []struct {
		Url                   string `json:"url"`
		Type                  string `json:"type"`
		HasOriginalDimensions bool   `json:"hasOriginalDimensions"`
	} `json:"media"`
}

func BuildPageUrl(blog string, page_number string) {

}
