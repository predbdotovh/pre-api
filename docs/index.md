# Predb.ovh API documentation

## Description

Disclaimer : This API has been built for my personal usage, and may not fit your needs. You may fill issues in the Github project, but don't expect too much.

Disclaimer 2 : Once again, this API does **not** host any content, only metadata associated to scene releases.

## Status

This API is currently usable, though not stable and may break at any point.

## API versions

Current version is **v1**

Current base URL is : [https://predb.ovh/api/v1/](https://predb.ovh/api/v1/)

## Responses

The API will always return HTTP 200 with application/json content.

On success :
```
{
    "status": "success",
    "message": "",
    "data": "... See endpoints responses ..."
}
```

On failure :
```
{
    "status": "error",
    "message": "Human readable error message",
    "data": null
}
```

## Endpoints

### GET /
List releases matching a set of filters given in parameters

* Cache : 60 seconds
* Usage : Get releases info

#### Parameters
TODO

#### Response
TODO

#### Example
* [https://predb.ovh/api/v1/?q=bdrip](https://predb.ovh/api/v1/?q=bdrip)

```
{
    "status": "success",
    "message": "",
    "data": {
        "rowCount": 20,
        "rows": [
            {
                "id": 7813747,
                "name": "Fortitude.S02E10.German.DUBBED.BDRip.x264-AIDA",
                "team": "AIDA",
                "cat": "TV-DVDRIP",
                "genre": "",
                "url": "http://imdb.com/title/tt3498622/",
                "size": 0,
                "files": 0,
                "preAt": 1493815560,
                "nuke": null
            },
            {
                "snip": "Content snipped for lisibility"
            },
            {
                "id": 7812994,
                "name": "Fortitude.S02E07.BDRip.x264-HAGGiS",
                "team": "HAGGiS",
                "cat": "TV-DVDRIP",
                "genre": "",
                "url": "http://imdb.com/title/tt3498622",
                "size": 258,
                "files": 19,
                "preAt": 1493762423,
                "nuke": null
            }
        ],
        "offset": 0,
        "reqCount": 20,
        "total": 70729,
        "time": 0.108095174
    }
}
```


### GET /live
This method is the exact clone of [GET /](#get), without any HTTP cache.

* Cache : None
* Usage : Get fresh data before listening to websocket updates

**To avoid abuse, this method is severly rate limited**

#### Example
* [https://predb.ovh/api/v1/live](https://predb.ovh/api/v1/live)


### GET /stats
Generate database statistics

* Cache : 60 seconds
* Usage : Keep track of current database status and response times

#### Parameters
TODO

#### Response
TODO

#### Example
* [https://predb.ovh/api/v1/stats](https://predb.ovh/api/v1/stats)

TODO


### GET /ws
This method binds to a websocket, sending near realtime updates.

* Cache : None

#### Parameters
There is no parameter, and any input on the websocket is discarded

#### Response
TODO

#### Example
* [https://predb.ovh/api/v1/ws](https://predb.ovh/api/v1/ws)

TODO

