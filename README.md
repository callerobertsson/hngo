# Hacker News Command Line

Simple command line program for listing and viewing Hacker News (https://news.ycombinator.com/).

## Status

`hngo` - fetches the top stories

`hngo <INDEX>` - opens the story with index <INDEX>

## Configuration

When run first time the config file (~/.hngorc) will be created with default values.

    {
        "ApiBaseUrl": "https://hacker-news.firebaseio.com/v0/",
        "ItemsLimit": 10,
        "CacheFilePath": "/tmp/hngocache",
        "OpenCommand": "echo"
    }

You would like to change OpenCommand to something better, like "open" for OSX.

Note that the `echo` command doesn't really print anything in the console.
It's run in a different process without access to stdout of the current console.

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

### More general open command

Make OpenCommand into an array so it can accept options for the command.

Perhaps add a tag for the Storys URL so it can be anywhere in the command line.

### Improve warning

Print a warning if the cache file is stale, i.e older than a defined age.

