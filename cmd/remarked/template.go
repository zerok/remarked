package main

import (
	"html/template"
	"io/ioutil"

	"github.com/pkg/errors"
)

var outputTemplate = `<!DOCTYPE html>
<html>
  <head>
	<title>{{ .Title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta charset="utf-8">
	{{ if .StyleSheetURL }}
	<link rel="stylesheet" href="{{ .StyleSheetURL }}">
	{{ end }}
  </head>
  <body>
	<textarea id="source">{{.Source}}</textarea>
    <script src="{{ .RemarkJS }}"></script>
    <script>
      var slideshow = remark.create({
		highlightLines: true
	  });
	  {{ if or .IsGuided }}
	  function connect() {
	  var socket = new WebSocket((window.location.protocol === "https:" ? "wss://" : "ws://") + window.location.host + "/ws/guide{{ if not .IsGuide }}d{{ end }}");
	  window.addEventListener('close', function() {
	    socket.close();
	  });
	  socket.onclose = function() {
	    window.setTimeout(function() {connect();}, 2000);
	  };
		{{ if .IsGuide }}
			socket.onopen = function() {
				socket.send(JSON.stringify({
				  type: 'auth',
				  token: '{{.Token}}'
				}));
				socket.send(JSON.stringify({
				  type: 'goto',
				  slideIndex: slideshow.getCurrentSlideIndex()
				}));
			};
			slideshow.on('showSlide', function(slide) {
				socket.send(JSON.stringify({
				type: 'goto',
				slideIndex: slide.getSlideIndex()
				}));
			});
		{{ else }}
			socket.onmessage = function(evt) {
				var cmd = JSON.parse(evt.data);
				switch (cmd.type) {
					case 'next': slideshow.gotoNextSlide(); break;
					case 'prev': slideshow.gotoPreviousSlide(); break;
					case 'goto': slideshow.gotoSlide(cmd.slideIndex + 1); break;
				}
			}
		{{ end }}
	  }
	  connect();
	  {{ end }}
    </script>
  </body>
</html>
`

func loadOutputTemplate(path string) (*template.Template, error) {
	var rawTemplate string
	if path == "" {
		rawTemplate = outputTemplate
	} else {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read template file %s", path)
		}
		rawTemplate = string(data)
	}
	tmpl, err := template.New("ROOT").Parse(rawTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template file %s", path)
	}
	return tmpl, nil
}
