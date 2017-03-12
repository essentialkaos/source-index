<!DOCTYPE html>
<html>
  <head>
    <meta charset='utf-8'>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta http-equiv="Content-Language" content="en">

    <link rel="shortcut icon" href="favicon.ico">
    <link rel="icon" type="image/x-icon" href="favicon.ico">

    <title>Index</title>

    <link href='https://fonts.googleapis.com/css?family=Roboto:400,300,700' rel='stylesheet' type='text/css'>

    <style type="text/css">
      html,body { color: #222; font-family: Roboto, Verdana, sans-serif; height: 100%; margin: 0; padding: 0 }
      h1,h2,h3,h4,h5,h6 { color: #666; font-weight: 100; padding-top: 16px }
      h1 { border-bottom: 1px #DDD solid; padding-bottom: 8px }
      h2 { padding-top: 16px }
      div { position: relative }
      div a { text-decoration: none }
      div.content { display: block; margin-left: auto; margin-right: auto; padding-top: 32px; width: 600px }
      div.release { background-color: #FFF; border-radius: 2px; margin: 0 0 4px -4px; padding: 4px 0 4px 4px; }
      div.release:hover { background-color: #F7F7F7 }
      div.release::before { color: #CCC; content: attr(data-date); font-size: 80%; margin-right: 12px; margin-top: 6px; position: absolute; right: 100% }
      span.version { font-size: 1.3em }
      span.date { color: #999; font-size: .8em }
      span.latest { font-weight: 700 }
      span.badge { background-color: #BBB; border-radius: 4px; color: #FFF; font-size: 60%; font-weight: 700; padding: 2px 4px; transition: background-color .5s ease-out; vertical-align: middle }
      span.format0:hover { background-color: #DD7A58 }
      span.format1:hover { background-color: #DDC458 }
      span.format2:hover { background-color: #B0DD58 }
      span.format3:hover { background-color: #58A3DD }
      span.format4:hover { background-color: #9158DD }
      table { border-spacing: 0 }
      td.version { width: 80px }
      p.footer { color: #999; font-size: 80%; padding: 64px 0 40px; text-align: center }
      p.footer a { color: #666 }
    </style>
  </head>
  <body>
    <div class="content">
      <h1>Source's Index</h1>
      {{ range .Projects }}
        <div>
          <h2>{{ .Name }}</h2>
           {{ range .Releases }}
              <div data-date="{{ .Date }}" class="release">
                <table><tr>
                <td class="version">
                  <span class="version{{ if .Latest }} latest{{ end }}">{{ .Version }}</span>
                </td>
                <td>
                {{ range $i, $s := .Sources }}
                  <a href="{{ $s.File }}"><span class="badge format{{ $i }}">{{ $s.Ext}}</span></a>
                {{ end }}
                </td>
                </tr></table>
              </div>
            {{ end }}
        </div>
      {{ end }}
    </div>
    <p class="footer">Genereated with ❤ by <a href="https://github.com/essentialkaos/source-index">SourceIndex</a><p>
  </body>
</html>
