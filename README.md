# findanyfilm-api
An API for http://findanyfilm.com written in Go. A directory of films and UK Cinema Listings.

### Endpoints

#### Movies
- /movies/find-by-name/:name
- /movies/find-by-cinema-date/:cinema/:date

#### Cinemas
- /cinemas/find-by-postcode/:postcode
- /cinemas/find-by-movie-date-postcode/:movie/:date/:postcode
