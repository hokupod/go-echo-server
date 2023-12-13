go-echo-server
==

Simple Echo Web Server

Usage
--
Start echo server
```
go run main.go
```

Send a web request
```
curl "http://localhost:8080/test?hoge=aa&hoge=bb&fuga=cc"

{
  "client_ip": "[::1]",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/8.4.0"
    ]
  },
  "method": "GET",
  "host": "localhost:8080",
  "path": "/test",
  "params": {
    "hoge": [
      "aa",
      "bb"
    ],
    "fuga": [
      "cc"
    ]
  }
}
```

Options
--
- `-port`: Select the port to listen on (default: 8080)
- `-authKey`: Use the GET parameter authKey to set the authentication key (default: not used)

Environments
--
- `PORT`: Select the port to listen on (default: 8080)
- `AUTH_KEY`: Use the GET parameter authKey to set the authentication key (default: not used)
- `SLACK_WEBHOOK_URL`: Specify the slack incoming webhook url. This will be echoed to slack as well (default: not used)

Docker Support
--

### Usage

Start echo server
```
docker run -p 8080:8080 hokupod/go-echo-server:latest
```

### Environments
- `PORT`: Select the port to listen on (default: 8080)
- `AUTH_KEY`: Use the GET parameter authKey to set the authentication key (default: not used)
- `SLACK_WEBHOOK_URL`: Specify the slack incoming webhook url. This will be echoed to slack as well (default: not used)
