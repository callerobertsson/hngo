# Hacker News Command Line

Simple command line program for listing and viewing Hacker News (https://news.ycombinator.com/).

Tested on OSX and Linux. 

Not tested on Windows where it might be some issues with creating the config file.

## Status

`hngo` - fetches the top stories

`hngo <INDEX>` - opens the story with index <INDEX>

## Configuration

When run first time the config file (~/.hngorc) will be created with default values.

    {
        "ApiBaseUrl": "https://hacker-news.firebaseio.com/v0/",
        "ItemsLimit": 10,
        "CacheFilePath": "/tmp/hngocache",
        "OpenCommand": ["echo", "-n"],
        "ShowCommandOutput", true
    }

You would like to change OpenCommand to something better, like "open" for OSX.

OpenCommand is an array of strings forming the open command and its options. 
The Story URL is appended to the end as the last option.
A future improvement will be to add a placeholder for where to insert the URL.

## Wishlist

_Future improvements_

### Show story status

Add a symbol indicating status for each story in the story list.

     0 * A story that's been viewed
     1 + A new story since last fetch
     2 ^ A story that has been moved up since last fetch
     3 v A story that has been moved down since last fetch
     4 - A story that is in the same position since last fetch
     5 E A story that has some kind of error

### Improve OpenCommand

Add a tag for the Storys URL so it can be anywhere in the command line.

### Improve warning

Print a warning if the cache file is stale, i.e older than a defined age.

