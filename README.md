# findanyfilm-api
An API for http://findanyfilm.com written in Go. A directory of films and UK Cinema Listings.

### Endpoints

#### Movies
- /movies/find-by-name/:name
- /movies/find-by-cinema-date/:cinema/:date

##### Examples
- /movies/find-by-name/The+Conjuring+2
- /movies/find-by-cinema-date/9358/2016-06-22

#### Cinemas
- /cinemas/find-by-postcode/:postcode
- /cinemas/find-by-movie-date-postcode/:movie/:date/:postcode
