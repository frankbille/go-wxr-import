package wxr

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
	"time"
)

func TestParseWxr(t *testing.T) {
	Convey("When parsing a valid file", t, func() {
		data, err := ioutil.ReadFile("example.xml")

		So(err, ShouldBeNil)

		wxr := ParseWxr(data)

		Convey("ParseWxr should return Wxr struct", func() {
			So(wxr, ShouldNotBeNil)

			Convey("Wxr should only contain 1 channel", func() {
				So(len(wxr.Channels), ShouldEqual, 1)
			})

			Convey("First channel should have correctly parsed fields", func() {
				channel := wxr.Channels[0]

				dateTime := time.Date(2014, time.July, 6, 23, 32, 34, 0, time.UTC).Unix()

				So(channel.Title, ShouldEqual, "Example")
				So(channel.Link, ShouldEqual, "http://example.com")
				So(channel.Description, ShouldEqual, "Example Site")
				So(channel.PubDate.Unix(), ShouldEqual, dateTime)
				So(channel.Language, ShouldEqual, "en-US")
				So(channel.WxrVersion, ShouldEqual, "1.2")
				So(channel.BaseSiteUrl, ShouldEqual, "http://example.com")
				So(channel.BaseBlogUrl, ShouldEqual, "http://example.com")
				So(channel.Generator, ShouldEqual, "http://wordpress.org/?v=3.9.1")

				Convey("There should be 2 authors, named John and Jane Doe", func() {
					authors := channel.Authors

					So(len(authors), ShouldEqual, 2)

					author := authors[0]
					So(author.Id, ShouldEqual, 1)
					So(author.Login, ShouldEqual, "johndoe")
					So(author.Email, ShouldEqual, "johndoe@example.com")
					So(author.DisplayName, ShouldEqual, "John Doe")
					So(author.FirstName, ShouldEqual, "John")
					So(author.LastName, ShouldEqual, "Doe")

					author = authors[1]
					So(author.Id, ShouldEqual, 2)
					So(author.Login, ShouldEqual, "janedoe")
					So(author.Email, ShouldEqual, "janedoe@example.com")
					So(author.DisplayName, ShouldEqual, "Jane Doe")
					So(author.FirstName, ShouldEqual, "Jane")
					So(author.LastName, ShouldEqual, "Doe")
				})

				Convey("There should be 2 categories, cat1 and cat2", func() {
					categories := channel.Categories

					So(len(categories), ShouldEqual, 2)

					category := categories[0]
					So(category.Id, ShouldEqual, 1)
					So(category.NiceName, ShouldEqual, "cat1")
					So(category.Parent, ShouldBeBlank)
					So(category.Name, ShouldEqual, "Category 1")

					category = categories[1]
					So(category.Id, ShouldEqual, 2)
					So(category.NiceName, ShouldEqual, "cat2")
					So(category.Parent, ShouldEqual, "cat1")
					So(category.Name, ShouldEqual, "Category 2")
				})

				Convey("There should be 2 tags, tag1 and tag2", func() {
					tags := channel.Tags

					So(len(tags), ShouldEqual, 2)

					tag := tags[0]
					So(tag.Id, ShouldEqual, 3)
					So(tag.Slug, ShouldEqual, "tag1")
					So(tag.Name, ShouldEqual, "Tag 1")

					tag = tags[1]
					So(tag.Id, ShouldEqual, 4)
					So(tag.Slug, ShouldEqual, "tag2")
					So(tag.Name, ShouldEqual, "Tag 2")
				})

				Convey("There should be 3 items, 1 post with 2 attachments", func() {
					items := channel.Items

					So(len(items), ShouldEqual, 3)

					item := items[0]
					So(item.Title, ShouldEqual, "Attachment Title 1")
					So(item.Link, ShouldEqual, "http://example.com/2014/07/06/post1/image1.jpg")
					So(item.PubDate.Unix(), ShouldEqual, dateTime)
					So(item.Creator, ShouldEqual, "johndoe")
					So(item.Guid, ShouldEqual, "http://example.com/wp-content/uploads/2014/07/image1.jpg")
					So(item.Description, ShouldBeBlank)
					So(item.Content, ShouldBeBlank)
					So(item.Excerpt, ShouldBeBlank)

					item = items[1]
					So(item.Title, ShouldEqual, "Attachment Title 2")
					So(item.Link, ShouldEqual, "http://example.com/2014/07/06/post1/image2.jpg")
					So(item.PubDate.Unix(), ShouldEqual, dateTime)
					So(item.Creator, ShouldEqual, "janedoe")
					So(item.Guid, ShouldEqual, "http://example.com/wp-content/uploads/2014/07/image2.jpg")
					So(item.Description, ShouldBeBlank)
					So(item.Content, ShouldBeBlank)
					So(item.Excerpt, ShouldBeBlank)

					item = items[2]
					So(item.Title, ShouldEqual, "Post 1")
					So(item.Link, ShouldEqual, "http://example.com/2014/07/06/post1/")
					So(item.PubDate.Unix(), ShouldEqual, dateTime)
					So(item.Creator, ShouldEqual, "johndoe")
					So(item.Guid, ShouldEqual, "http://example.com/?p=1")
					So(item.Description, ShouldBeBlank)
					So(item.Content, ShouldEqual, "Content of the post\nIncluding line breaks, which is translated to \\n character.")
					So(item.Excerpt, ShouldBeBlank)
				})
			})
		})
	})
}
