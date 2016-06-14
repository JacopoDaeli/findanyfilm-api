package main

import (
  "log"
  "fmt"
  "bytes"
  "regexp"
  "strings"
  "github.com/PuerkitoBio/goquery"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/cinemas/find-by-postcode/:postcode", func(c *gin.Context) {
    var url bytes.Buffer
    postcode := c.Param("postcode")
    url.WriteString("http://www.findanyfilm.com/find-cinema-tickets?townpostcode=")
    url.WriteString(strings.Replace(postcode, " ", "%20", -1))

    fmt.Printf(url.String())

    doc, err := goquery.NewDocument(url.String())
    if err != nil {
      log.Fatal(err)
    }

    ids := make([]string, 10)
    last := -1

    doc.Find(".cinemaResult").EachWithBreak(func(i int, s *goquery.Selection) bool {
      attr, exists := s.Attr("rel")

      if i > 9 || exists == false {
        return false
      }

      re := regexp.MustCompile("[0-9]+")
      tmpId := re.FindAllString(attr, -1)
      ids[i] = tmpId[0]
      last = i

      return true
    })

    c.String(200, ids[last])
  })
  r.Run() // listen and server on 0.0.0.0:8080
}
