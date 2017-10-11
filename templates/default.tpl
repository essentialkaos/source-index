<!DOCTYPE html>
<html>
  <head>
    <meta charset='utf-8'>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta http-equiv="Content-Language" content="en">

    <link rel="shortcut icon" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAQAAADZc7J/AAAA6UlEQVR4Ae2TX9HEIAzE6wAJlQBkBVQCEpCABCRUQiVUAhJOQiUg4W7myjdD6J/wdA/fNPsGyW9mNzD823rKajt9pfm5cUhF4WKQFiS8marW6i6djNNcLlu9oDoA8HyMKXUAaLsBFCPGUdwFf3Bft9OGxLOgLOU+sfa4m6L4JwTBgtWta4rGtYusO44Z5HP3lBGgegCRD3KI1QJARGwioES5XiGMEwAM4yli5alQvAUgVK0LVHMmA2hpkk9IzV78GYA9JEGKA+S/yOU7tnC1AcrGSQ+plHEthDLNUOy1yp/LjnYqGoenflEfVoabpmprvccAAAAASUVORK5CYII=" />

    <title>Source's Index</title>

    <link href='https://fonts.googleapis.com/css?family=Roboto:400,300,700' rel='stylesheet' type='text/css'>

    <style type="text/css">
      html,body { color: #222; font-family: Roboto, Verdana, sans-serif; height: 100%; margin: 0; min-width: 820px; padding: 0 }
      div,section { position: relative }
      h1,h2,h3,h4,h5,h6 { color: #666; font-weight: 100; padding-top: 16px }
      h1 { border-bottom: 1px #DDD solid; padding-bottom: 8px }
      h2.project { display: inline-block }
      a { text-decoration: none }
      a.anchor { color: #CCC; font-size: 1.2em; font-weight: 700; opacity: 0 }
      a.anchor:hover { color: #999 }
      section:hover a.anchor { opacity: 1 }
      div.content { display: block; margin-left: auto; margin-right: auto; padding-top: 32px; width: 600px }
      div.release { background-color: #FFF; border-radius: 4px; margin: 0 0 4px -8px; padding: 4px 0 4px 8px; }
      div.release:hover { background-color: #F7F7F7 }
      div.release::before { color: #CCC; content: attr(data-date); font-size: 80%; margin-right: 12px; margin-top: 6px; position: absolute; right: 100% }
      span.version { font-size: 1.3em }
      span.date { color: #999; font-size: .8em }
      span.latest { font-weight: 700 }
      span.badge { background-color: #BBB; border-radius: 4px; color: #FFF; font-size: 60%; font-weight: 700; padding: 2px 4px; transition: background-color .25s ease-out; vertical-align: middle }
      span.format0:hover { background-color: #DD7A58 }
      span.format1:hover { background-color: #DDC458 }
      span.format2:hover { background-color: #B0DD58 }
      span.format3:hover { background-color: #58A3DD }
      span.format4:hover { background-color: #9158DD }
      dl { margin: 0; overflow: hidden; padding: 0 }
      dt { float: left; width: 96px }
      p.footer { color: #999; font-size: 80%; padding: 64px 0 40px; text-align: center }
      p.footer a { color: #666; border-bottom: 1px solid #666 }
    </style>
  </head>
  <body>
    <div class="content">
      <h1>Source's Index</h1>
      {{ range .Projects }}
        <section id="{{ .Name }}">
          <h2 class="project">{{ .Name }}</h2>
          <a href="#{{ .Name }}" class="anchor">#</a>
            {{ range .Releases }}
              <div data-date="{{ .Date }}" class="release">
                <dl>
                  <dt><span class="version{{ if .Latest }} latest{{ end }}">{{ .Version }}</span></dt>
                  <dd>
                  {{ range $i, $s := .Sources }}
                    <a href="{{ $s.File }}"><span class="badge format{{ $i }}">{{ $s.Ext}}</span></a>
                  {{ end }}
                  </dd>
                </dl>
              </div>
            {{ end }}
        </section>
      {{ end }}
    </div>
    <p class="footer">Genereated with ❤ by <a href="https://github.com/essentialkaos/source-index">SourceIndex</a></p>
  </body>
</html>
