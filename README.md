# Short hash

This a URL shortener service: given a URL, it returns a shorter version
of it.

When the short URL is used, you get redirected to the original location.

## Endpoints

### Shorten URL
```http request
POST /urls
```

Returns the short version of the URL

### Get original URl
```http request
GET /urls/{short_url}
```

Redirects you to the original URL. Returns an HTTP 307 TEMPORARY REDIRECT status code