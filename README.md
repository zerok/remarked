# Remarked

[![Build Status](https://dev.azure.com/hgutmann/remarked/_apis/build/status/zerok.remarked)](https://dev.azure.com/hgutmann/remarked/_build/latest?definitionId=2)

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
- `templateFile`: remarked generates a simple HTML output in which remarkJS is
  included. The template that should be used for that HTML output can be
  customized with this flag. You can find the default template on
  [GitHub](https://github.com/zerok/remarked/blob/master/cmd/remarked/template.go).
- `stylesheet`: If you need any custom styling, specify your CSS file here.
- `title`: The title as it is rendered inside the browser's title bar.
- `remarkJS`: If you prefer a modified version of Remark.JS, specify it here.
- `staticFolder`: This folder will be made available under `/static` by the
  built-in webserver.
- `markdownAsTemplate`: If you set this to  `true` then the Markdown file
  will be treated as a template file for Go's [html/template](https://golang.org/pkg/html/template/)
  package.
- `leftActionDelimiter`: Used within `html/template` (Default: `{{`)
- `rightActionDelimiter`: Used within `html/template` (Default: `}}`)

All of these can be overriden with command-line flags.


## Markdown as a template

If you set `markdownAsTemplate` to `true` inside the remarked.yml file, 
remarked will try to parse the specified Markdown file using Go's
[html/template](https://golang.org/pkg/html/template/) package. This allows
you to do things like if branching or loops. For now, the following functions
are provided:

- `loadCode PATH`: Loads the content of the given PATH and renders includes it
  into the content.

- `markLines RANGES CONTENT`: Parses the given content and adds a `*` in front
  of every line matching the given ranges. Ranges can be specified as a 
  comma-separated list of either positive numbers or `START-STOP` ranges.

  E.g.: `2-5,8` would highlight lines 2, 3, 4, 5, and 8.

- `counter START END STEP` is a helper for creating range-loops of a custom
  size:

  ```
  {{ range (counter 0 2 1) }}
    {{ .Counter }}
  {{ end }}
  ```

  The output of this snippet is `0 1 2` as output.

