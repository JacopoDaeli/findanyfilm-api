package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "FindAnyFilm API: An API for http://findanyfilm.com written in Go. A directory of films and UK Cinema Listings.")
	})

	r.GET("/movies/find-by-name/:name", func(c *gin.Context) {
		var url bytes.Buffer
		name := c.Param("name")
		url.WriteString("http://www.findanyfilm.com/search/live-film?term=")
		url.WriteString(strings.Replace(name, " ", "+", -1))

		resp, err := http.Get(url.String())

		defer resp.Body.Close()

		if err != nil {
			c.JSON(500, err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, err)
			return
		}

		if resp.StatusCode >= 400 {
			var errorMessage bytes.Buffer
			errorMessage.WriteString("FindAnyFilm returns: ")
			errorMessage.WriteString(resp.Status)
			c.JSON(500, gin.H{
		    "status": 500,
		    "error": errorMessage.String(),
		  })
			return
		}

		movies := make([]Movie, 0)
		json.Unmarshal(body, &movies)

		if len(movies) == 0 {
		  c.JSON(404, gin.H{
		    "status": 404,
		    "error": "Movie not found",
		  })
		} else {
		  c.JSON(resp.StatusCode, movies)
		}
	})

	r.GET("/movies/find-by-cinema-date/:cinema/:date", func(c *gin.Context) {
		var url bytes.Buffer

		cinema := c.Param("cinema")
		date := c.Param("date")

		url.WriteString("http://www.findanyfilm.com/api/screenings/by_venue_id/venue_id/")
		url.WriteString(cinema)
		url.WriteString("/date_from/")
		url.WriteString(date)

		fmt.Printf(url.String())

		resp, err := http.Get(url.String())

		defer resp.Body.Close()

		if err != nil {
			c.JSON(500, err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, err)
			return
		}

		var jResp bytes.Buffer
		jResp.Write(body)
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.String(resp.StatusCode, jResp.String())
	})

	r.GET("/cinemas/find-by-postcode/:postcode", func(c *gin.Context) {
		var url bytes.Buffer
		postcode := c.Param("postcode")
		url.WriteString("http://www.findanyfilm.com/find-cinema-tickets?townpostcode=")
		url.WriteString(strings.Replace(postcode, " ", "%20", -1))

		fmt.Printf(url.String())

		doc, err := goquery.NewDocument(url.String())
		if err != nil {
			c.JSON(500, err)
			return
		}

		cinemas := make([]Cinema, 10)
		last := -1

		doc.Find(".cinemaResult").EachWithBreak(func(i int, s *goquery.Selection) bool {
			attr, exists := s.Attr("rel")

			if i > 9 || exists == false {
				return false
			}

			re := regexp.MustCompile("[0-9]+")
			tmpId := re.FindAllString(attr, -1)

			cinemaName := s.Find("span").Text()
			cinemas[i] = Cinema{tmpId[0], cinemaName}

			last = i

			return true
		})

		c.JSON(200, cinemas)
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

		defer resp.Body.Close()

		if err != nil {
			c.JSON(500, err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, err)
			return
		}

		var jResp bytes.Buffer
		jResp.Write(body)
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.String(resp.StatusCode, jResp.String())
	})

	r.Run() // listen and server on 0.0.0.0:8080
}
