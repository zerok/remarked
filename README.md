# Remarked

Remarked is a simple wrapper around remark.js that makes creating new
presentations and holding them a bit easier. It sports the following features
on top of the original framework:

- A simple HTTP server for which you can specify the Markdown and CSS file to
  use.
- A guide/guided mode for being able to remote-control a presentation through a
  websocket connection.

## Installation

You can either pick one of the [binaries from Github](https://github.com/zerok/remarked/releases)
or, if you have Go 1.9+ installed, by using go-get:

```
$ go get github.com/zerok/remarked/cmd/remarked
```


## Getting started

Once you have installed remarked, you can create a sample project file with
`remarked --init`. See that file for descriptions on all the available
settings.

Then simply write your Markdown as you would in preparation for using it with
RemarkJS and then start `remarked` to launch a small webserver with your
presentation on.


## Remote 

If you start remarked with the `--guide` flag, you can access the `/guide`
endpoint to control other clients that just opened the normal `/` endpoint.
When you access the guide-endpoint for the first time, you will be asked for
a token which was printed in the terminal you used to start remarked.

This feature is using websockets in the background to send commands from the 
guide-instance to the guided-instance.

**Note:** When you first log into the guide-mode, the token will be sent
directly through the websocket. In order to keep it secure, please access
remarked through an HTTPS connection.


## Styling

Remarked will use the default remark.js styling which you can extend or
override using the `stylesheet` setting in the configuration file or the
`--stylesheet` command-line flag.

If you set a local file as stylesheet, it will be served by the HTTP server as
`/style/_.css`.


## Configuration

By default, remarked will look for a `remarked.yml` file within the current
working directory with the following options:

- `markdownFile`: This file contains your presentation content. See the
  Remark.JS documentation for details on how this file has to be formatted.
- `stylesheet`: If you need any custom styling, specify your CSS file here.
- `title`: The title as it is rendered inside the browser's title bar.
- `remarkJS`: If you prefer a modified version of Remark.JS, specify it here.
- `staticFolder`: This folder will be made available under `/static` by the
  built-in webserver.

All of these can be overriden with command-line flags.
