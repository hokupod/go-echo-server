go-echo-server
==

Simple Echo Web Server

Usage
--
Start echo server
```
go run main.go
```

Options
--
- `-port`: Select the port to listen on (default: 8080)
- `-authKey`: Use the GET parameter authKey to set the authentication key (default: not used)

Environments
--
- `PORT`: Select the port to listen on (default: 8080)
- `AUTH_KEY`: Use the GET parameter authKey to set the authentication key (default: not used)

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
