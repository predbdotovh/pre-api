# Predb.ovh API documentation

Most pre databases websites do not allow data scapping, nor do they offer an API.

Having a easily accessible API was a requirement for me, so I built this one using public pre sources.

## Disclaimers

Disclaimer : This API has been built for my personal usage, and may not fit your needs.
[Github issues](https://github.com/predbdotovh/pre-api/issues) can be filled, but don't expect too much.

Disclaimer 2 : Once again, this API does **NOT** host any content, only metadata associated to scene releases.

Disclaimer 3 : Don't know what scene releases are ? You're probably at the wrong place.

## Status

This API is currently usable, and used by [demo website](https://predb.ovh/) which is
also [open-sourced](https://github.com/predbdotovh/website-vuejs).

As the API is fed by Sphinx, results are hard limited at 1000, and I don't expect to modify this behaviour in the
future.

A monthly-ish sql dump is available here : [https://predb.ovh/download/](https://predb.ovh/download/)

## API versions

Current version is **v1**. This version is expected to be updated on each breaking change.

Current base URL is : [https://predb.ovh/api/v1/](https://predb.ovh/api/v1/)

## Responses

The API will always return HTTP 200 with application/json content (except for [RSS](#get-rss)).

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

- Cache : 60 seconds
- Usage : Get releases info
- Rate limit : 30/60s

#### Parameters

All parameters are optional

| param key | type   | content                              |
| --------- | ------ | ------------------------------------ |
| count     | int    | Maximum releases count expected      |
| page      | int    | Page offset                          |
| offset    | int    | Row offset (overwrites page param)   |
| q         | string | Query                                |
| id        | int    | Specific pre ID (overwrites q param) |

Query is directly fed to a SphinxSearch engine,
allowing [specific syntax](http://sphinxsearch.com/docs/current/extended-syntax.html). Note: cat and team are indexed,
allowing fast queries like
[https://predb.ovh/api/v1/?q=@cat%20EBOOK](https://predb.ovh/api/v1/?q=@cat%20EBOOK)

#### Response

##### Data

| json key | type      | content                       |
| -------- | --------- | ----------------------------- |
| rowCount | int       | Count of rows returned        |
| offset   | int       | Row count offset requested    |
| reqCount | int       | Row count requested           |
| total    | int       | Total matching rows           |
| time     | float     | Request internal duration     |
| rows     | []release | Array of [releases](#release) |

##### Release

| json key | type      | content                           |
| -------- | --------- | --------------------------------- |
| id       | int       | Internal unique pre ID            |
| name     | string    | Release name                      |
| team     | string    | Release group extracted from name |
| cat      | string    | Category                          |
| genre    | string    | Genre                             |
| url      | string    | Info link                         |
| size     | float     | Release size in kb                |
| files    | int       | Original file count               |
| preAt    | int       | Release pre timestamp             |
| nuke     | nuke/null | [Nuke](#nuke) info if available   |

##### Nuke

| json key | type   | content                     |
| -------- | ------ | --------------------------- |
| id       | int    | Internal unique nuke ID     |
| typeId   | int    | [Nuke type](#nuke-types) ID |
| type     | string | [Nuke type](#nuke-types)    |
| preId    | int    | Nuked pre ID                |
| reason   | string | Nuke reason                 |
| net      | string | Nuke source net             |
| nukeAt   | int    | Nuke timestamp              |

##### Nuke types

Known nuke types and type ids

| nuke type ID | nuke type |
| ------------ | --------- |
| 1            | nuke      |
| 2            | unnuke    |
| 3            | modnuke   |
| 4            | delpre    |
| 5            | undelpre  |

#### Example

- [https://predb.ovh/api/v1/?q=bdrip](https://predb.ovh/api/v1/?q=bdrip)

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

- Cache : None
- Usage : Get fresh data before listening to websocket updates
- Rate limit : 2/20s

**To avoid abuse, this method is severly rate limited**

#### Example

- [https://predb.ovh/api/v1/live](https://predb.ovh/api/v1/live)

```
{
    "status": "success",
    "message": "",
    "data": {
        "rowCount": 20,
        "rows": [
            {
                "id": 7814699,
                "name": "AllOver30.com_17.05.03.Drew.Jones.XXX.iMAGESET-YAPG",
                "team": "YAPG",
                "cat": "XXX-iMGSET",
                "genre": "Mature.Women",
                "url": "http://HTTP://AllOver30.com",
                "size": 211.9,
                "files": 45,
                "preAt": 1493881734,
                "nuke": null
            },
            {
                "snip": "Content snipped for lisibility"
            },
            {
                "id": 7814680,
                "name": "Robert_Abigail_Feat_13_Amps_-_Living_On_The_Right_Side-(AG_010)-WEB-2017-ZzZz",
                "team": "ZzZz",
                "cat": "MP3",
                "genre": "Dance",
                "url": "http://junodownload.com",
                "size": 7,
                "files": 1,
                "preAt": 1493880783,
                "nuke": null
            }
        ],
        "offset": 0,
        "reqCount": 20,
        "total": 7751247,
        "time": 0.898239015
    }
}
```

### GET /rss

This method is the exact clone of [GET /](#get), formatted using RSS2.0 spec.

- Cache : 60 seconds
- Usage : Get releases info using a RSS reader
- Rate limit : 30/60s

#### Example

- [https://predb.ovh/api/v1/rss](https://predb.ovh/api/v1/rss)

```
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
    <channel>
        <title>PreDB</title>
        <link>https://predb.ovh/</link>
        <description></description>
        <pubDate>Sun, 15 Dec 2019 21:23:59 +0100</pubDate>
        <item>
            <title>Le.Steppe.Dell.Asia.La.Piu.Grande.Steppa.Della.Terra.iTALiAN.HDTV.x264-iDiB</title>
            <link>https://predb.ovh/?id=9643108</link>
            <description>Cat:TV-SD-X264 | Genre: | Size:0MB | Files:0 | ID:9643108</description>
            <pubDate>Sun, 15 Dec 2019 21:23:49 +0100</pubDate>
        </item>
        <item>
            <title>Le.Steppe.Dell.Asia.La.Piu.Grande.Steppa.Della.Terra.iTALiAN.720p.HDTV.x264-iDiB</title>
            <link>https://predb.ovh/?id=9643107</link>
            <description>Cat:TV-HD-X264 | Genre: | Size:0MB | Files:0 | ID:9643107</description>
            <pubDate>Sun, 15 Dec 2019 21:23:48 +0100</pubDate>
        </item>
    </channel>
</rss>
```

### GET /teams

Teams stats with first and latest pre, and total recorded pre.

This is currently hard limited to 1000 results, may be subject to change if necessary.

- Cache : 3600 seconds
- Usage : Teams statistics

#### Parameters

None

#### Response

##### Data

| json key | type      | content                       |
| -------- | --------- | ----------------------------- |
| rowCount | int       | Count of rows returned        |
| offset   | int       | Row count offset requested    |
| reqCount | int       | Row count requested           |
| total    | int       | Total matching rows           |
| time     | float     | Request internal duration     |
| rows     | []team    | Array of [teams](#team)       |

##### Team

| json key  | type      | content                               |
| --------- | --------- | ------------------------------------- |
| team      | string    | Release group name                    |
| firstPre  | int       | First recorded release pre timestamp  |
| latestPre | int       | Latest recorded release pre timestamp |
| count     | int       | Total team pre count                  |

#### Example

- [https://predb.ovh/api/v1/teams](https://predb.ovh/api/v1/teams)

```
{
    "status": "success",
    "message": "",
    "data": {
        "rowCount": 1000,
        "rows": [
            {
                "team": "KTR",
                "firstPre": 1203561700,
                "latestPre": 1598341396,
                "count": 307933
            },
            {
                "snip": "Content snipped for lisibility"
            },
        ],
        "time": 5.465514844
    }
}
```

### GET /stats

Basic stats about internal database health.

- Cache : 60 seconds
- Usage : Keep track of current database status and response times

#### Parameters

None

#### Response

##### Data

| json key | type   | content                             |
| -------- | ------ | ----------------------------------- |
| total    | int    | Total indexed releases count        |
| date     | string | Current RFC3339 timestamp of server |
| time     | int    | Full index scan duration            |

#### Example

- [https://predb.ovh/api/v1/stats](https://predb.ovh/api/v1/stats)

```
{
    "status": "success",
    "message": "",
    "data": {
        "total": 7751280,
        "date": "2017-05-04T09:49:42.527504724+02:00",
        "time": 0.893255737
    }
}
```

### GET /ws

A websocket endpoint, sending near realtime updates. To use this, you need to bind using
a [Websocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_client_applications).

- Cache : None
- Usage : Realtime updates

#### Parameters

There is no parameter, and any input on the websocket is discarded

#### Response

There is no response per se, but a range of frames sent over time

##### Frame

Each frame represents an action towards a specific release.

| json key | type    | content             |
| -------- | ------- | ------------------- |
| action   | string  | Action type         |
| row      | release | [Release](#release) |

##### Action types

Known action types

| action   | context                                 |
| -------- | --------------------------------------- |
| insert   | First release pre                       |
| update   | Any release field update                |
| delete   | Erroneous release (should never happen) |
| nuke     | Release nuked by net                    |
| unnuke   | Release unnuked by net                  |
| modnuke  | Nuke reason modified by net             |
| delpre   | Pre deleted by net                      |
| undelpre | Pre undeleted by net                    |

#### Example

- [https://predb.ovh/api/v1/ws](https://predb.ovh/api/v1/ws)

```
{
    "action": "insert",
    "row": {
        "id": 7814727,
        "name": "Doom_Squad-Countdown_To_Doomsday_II-WEB-2016-ESG",
        "team": "ESG",
        "cat": "MP3",
        "genre": "",
        "url": "",
        "size": 0,
        "files": 0,
        "preAt": 1493884429,
        "nuke": null
    }
}
```
