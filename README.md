# Remarked

Remarked is a simple wrapper around remark.js that makes creating new
presentations and holding them a bit easier. It sports the following features
on top of the original framework:

- A simple HTTP server for which you can specify the Markdown and CSS file to
  use.
- A guide/guided mode for being able to remote-control a presentation through a
  websocket connection.


## Getting started

Once you have remarked installed, you can create a sample project file using
`remarked --init`. See that file for descriptions on all the available
settings.


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

- `markdownFile`
- `stylesheet`
- `staticFolder` (default: `./static`)

All of these can be overriden with commandline flags.
