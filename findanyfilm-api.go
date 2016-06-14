package main

import (
  "log"
  "fmt"
  "bytes"
  // "encoding/json"
  "net/http"
  "regexp"
  "strings"
  "io/ioutil"
  "github.com/PuerkitoBio/goquery"
  "github.com/gin-gonic/gin"
)

type Movie struct {
  Film_id string
  Link string
  Label string
  Category string
}

func main() {
  r := gin.Default()

  r.GET("/", func(c *gin.Context) {
    c.String(200, "FindAnyFilm API powered by Impero.")
  })

  r.GET("/movies/find-by-name/:name", func(c *gin.Context) {
    var url bytes.Buffer
    name := c.Param("name")
    url.WriteString("http://www.findanyfilm.com/search/live-film?term=")
    url.WriteString(strings.Replace(name, " ", "+", -1))

    fmt.Printf(url.String())

    resp, err := http.Get(url.String())
    if err != nil {
    	log.Fatal(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	log.Fatal(err)
    }

    var jResp bytes.Buffer
    jResp.Write(body)
    c.String(200, jResp.String())

    // movies := make([]Movie, 0)
    // json.Unmarshal(body, &movies)
    //
    // if len(movies) == 0 {
    //   c.JSON(404, gin.H{
    //     "status": 404,
    //     "error": "No movie found",
    //   })
    // } else {
    //   var jResp bytes.Buffer
    //   jResp.Write(body)
    //   c.String(200, jResp.String())
    // }
  })

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

  r.GET("/cinemas/find-by-movie-date-postcode/:movie/:date/:postcode", func(c *gin.Context) {
    var url bytes.Buffer

    movie := c.Param("movie")
    date := c.Param("date")
    postcode := c.Param("postcode")

    url.WriteString("http://www.findanyfilm.com/api/screenings/by_film_id/film_id/")
    url.WriteString(movie)
    url.WriteString("/date_from/")
    url.WriteString(date)
    url.WriteString("/townpostcode/")
    url.WriteString(strings.Replace(postcode, " ", "%20", -1))

    fmt.Printf(url.String())

    resp, err := http.Get(url.String())
    if err != nil {
      log.Fatal(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Fatal(err)
    }

    var jResp bytes.Buffer
    jResp.Write(body)
    c.String(200, jResp.String())
  })

  r.Run() // listen and server on 0.0.0.0:8080
}
