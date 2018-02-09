# Emoji URL Shortener
URL Shortener is server application which is used to create short URLs that can be easily shared, tweeted, or emailed to friends.

*Note*: please don't consider this project seriously. It was made just for fun and for experimental purposes.

## Instalation
```
$ go get -u github.com/zitryss/url-shortener
```

## Usage

On the server side:
```
$ url-shortener -domain example.com -port 8080
```

On the client side:
```
$ curl -F "url=https://wwww.website.com/extremely/long/url/" example.com:8080
http://example.com:8080/🗃🏇😈☃️📁⛳️🌧🦌
```
