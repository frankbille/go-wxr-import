/*
Parser for WordPress backup format WXR version 1.2
All structs also contains json marshal options, for easy serialization
to JSON format.

Example:

	import (
		"github.com/frankbille/go-wxr-import"
		"io/ioutil"
		"log"
	)

	func main() {
		wxrXmlData, _ := ioutil.ReadFile("wordpress-backup.xml")
		wxr := ParseWxr(wxrXmlData)
		log.Print(wxr)
	}
*/
package wxr

import (
	"encoding/xml"
	"time"
)

// Root tag, which is an <rss> tag in the XML.
type Wxr struct {
	Channels []Channel `xml:"channel" json:"channels"`
}

type Channel struct {
	Title         string     `xml:"title" json:"title"`
	Link          string     `xml:"link" json:"link"`
	Description   string     `xml:"description" json:"description"`
	PubDateString string     `xml:"pubDate" json:"-"`
	PubDate       time.Time  `xml:"-" json:"pub-date"`
	Language      string     `xml:"language" json:"language"`
	WxrVersion    string     `xml:"wxr_version" json:"wxr-version"`
	BaseSiteUrl   string     `xml:"base_site_url" json:"base-site-url"`
	BaseBlogUrl   string     `xml:"base_blog_url" json:"base-blog-url"`
	Generator     string     `xml:"generator" json:"generator"`
	Authors       []Author   `xml:"author" json:"authors"`
	Categories    []Category `xml:"category" json:"categories"`
	Tags          []Tag      `xml:"tag" json:"tags"`
	Items         []Item     `xml:"item" json:"items"`
}

type Author struct {
	Id          int32  `xml:"author_id" json:"id"`
	Login       string `xml:"author_login" json:"login"`
	Email       string `xml:"author_email" json:"email"`
	DisplayName string `xml:"author_display_name" json:"display-name"`
	FirstName   string `xml:"author_first_name" json:"first-name"`
	LastName    string `xml:"author_last_name" json:"last-name"`
}

type Category struct {
	Id       int32  `xml:"term_id" json:"id"`
	NiceName string `xml:"category_nicename" json:"nice-name"`
	Parent   string `xml:"category_parent" json:"parent"`
	Name     string `xml:"cat_name" json:"name"`
}

type Tag struct {
	Id   int32  `xml:"term_id" json:"id"`
	Slug string `xml:"tag_slug" json:"slug"`
	Name string `xml:"tag_name" json:"name"`
}

type Item struct {
	Title             string         `xml:"title" json:"title"`
	Link              string         `xml:"link" json:"link"`
	PubDateString     string         `xml:"pubDate" json:"-"`
	PubDate           time.Time      `xml:"-" json:"pub-date"`
	Creator           string         `xml:"creator" json:"creator"`
	Guid              string         `xml:"guid" json:"guid"`
	Description       string         `xml:"description" json:"description"`
	Content           string         `xml:"http://purl.org/rss/1.0/modules/content/ encoded" json:"content"`
	Excerpt           string         `xml:"http://wordpress.org/export/1.2/excerpt/ encoded" json:"excerpt"`
	PostId            int32          `xml:"post_id" json:"post-id"`
	PostDateString    string         `xml:"post_date" json:"-"`
	PostDate          time.Time      `xml:"-" json:"post-date"`
	PostDateGmtString string         `xml:"post_date_gmt" json:"-"`
	PostDateGmt       time.Time      `xml:"-" json:"post-date-gmt"`
	CommentStatus     string         `xml:"comment_status" json:"comment-status"`
	PingStatus        string         `xml:"ping_status" json:"ping-status"`
	PostName          string         `xml:"post_name" json:"post-name"`
	Status            string         `xml:"status" json:"status"`
	PostParent        int32          `xml:"post_parent" json:"post-parent"`
	MenuOrder         int32          `xml:"menu_order" json:"menu-order"`
	PostType          string         `xml:"post_type" json:"post-type"`
	PostPassword      string         `xml:"post_password" json:"post-password"`
	IsSticky          bool           `xml:"is_sticky" json:"sticky"`
	AttachmentUrl     string         `xml:"attachment_url" json:"attachment-url"`
	Categories        []ItemCategory `xml:"category" json:"categories"`
	PostMetas         []PostMeta     `xml:"postmeta" json:"post-meta"`
	Comments          []Comment      `xml:"comment" json:"comments"`
}

type ItemCategory struct {
	Domain      string `xml:"domain,attr" json:"domain"`
	NiceName    string `xml:"nicename,attr" json:"nice-name"`
	DisplayName string `xml:",chardata" json:"display-name"`
}

type PostMeta struct {
	Key   string `xml:"meta_key" json:"key"`
	Value string `xml:"meta_value" json:"value"`
}

type Comment struct {
	Id            int32     `xml:"comment_id" json:"id"`
	Author        string    `xml:"comment_author" json:"author-name"`
	AuthorEmail   string    `xml:"comment_author_email" json:"author-email"`
	AuthorUrl     string    `xml:"comment_author_url" json:"author-url"`
	AuthorIp      string    `xml:"comment_author_ip" json:"author-ip"`
	DateString    string    `xml:"comment_date" json:"-"`
	Date          time.Time `xml:"-" json:"date"`
	DateGmtString string    `xml:"comment_date_gmt" json:"-"`
	DateGmt       time.Time `xml:"-" json:"date-gmt"`
	Content       string    `xml:"comment_content" json:"content"`
	Approved      bool      `xml:"comment_approved" json:"approved"`
	Type          string    `xml:"comment_type" json:"type"`
	Parent        int32     `xml:"comment_parent" json:"parent"`
	UserId        int32     `xml:"comment_user_id" json:"user-id"`
}

const (
	ISODATEFORMAT = "2006-01-02 15:04:05"
)

// Parse an incoming byte array representing the XML file in WXR format.
// Returns a Wxr struct representing all data in the file.
func ParseWxr(data []byte) Wxr {
	var wxr Wxr
	xml.Unmarshal(data, &wxr)

	// Parse date strings
	for channelIndex, channel := range wxr.Channels {
		pubDate, _ := time.Parse(time.RFC1123Z, channel.PubDateString)
		wxr.Channels[channelIndex].PubDate = pubDate

		for itemIndex, item := range channel.Items {
			pubDate, _ := time.Parse(time.RFC1123Z, item.PubDateString)
			wxr.Channels[channelIndex].Items[itemIndex].PubDate = pubDate

			postDate, _ := time.Parse(ISODATEFORMAT, item.PostDateString)
			wxr.Channels[channelIndex].Items[itemIndex].PostDate = postDate

			postDateGmt, _ := time.Parse(ISODATEFORMAT, item.PostDateGmtString)
			wxr.Channels[channelIndex].Items[itemIndex].PostDateGmt = postDateGmt

			for commentIndex, comment := range item.Comments {
				commentDate, _ := time.Parse(ISODATEFORMAT, comment.DateString)
				wxr.Channels[channelIndex].Items[itemIndex].Comments[commentIndex].Date = commentDate

				commentDateGmt, _ := time.Parse(ISODATEFORMAT, comment.DateGmtString)
				wxr.Channels[channelIndex].Items[itemIndex].Comments[commentIndex].DateGmt = commentDateGmt
			}
		}
	}

	return wxr
}
