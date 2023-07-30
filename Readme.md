# DDoser [![Build Status](https://github.com/komandakycto/ddoser/workflows/build/badge.svg)](https://github.com/komandakycto/ddoser/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/komandakycto/ddoser)](https://goreportcard.com/report/github.com/komandakycto/ddoser) [![Coverage Status](https://coveralls.io/repos/github/komandakycto/ddoser/badge.svg)](https://coveralls.io/github/komandakycto/ddoser)

`DDoser` is a simple tool to find attackers IP addresses in an access log if you under DDoS attack.

## Supported log formats

By default `DDoser` supports nginx access log format.

## How it works

The DDoser reads `k bytes` from the end of log file each `n seconds`. After this it parses the log and finds IP
addresses of attackers.
Ip address is marked as attacker if it has more than `m` requests in last `t` seconds.

## Use cases

* You are under DDoS attack you want to find IP addresses of attackers to block them. Attackers are using
  different IP addresses for each request. It hard to find them in the log file just looking at the log.
* You want to find IP addresses which are sending a lot of requests to your server. It can be a sign of
  DDoS attack or just a bug in an application. You can set up rate limiting for these IP addresses.
* You want to analyze the geo distribution of your users. You can find IP addresses which are sending a lot of
  requests to your server and analyze their geo distribution.

## Usage

See full list of options in the help or in
the [config.go](https://github.com/komandakycto/ddoser/blob/master/app/config.go) file:

```
Usage of ./ddoser:
  --logpath string
    	Path to nginx access log
  --readinterval int
    	Interval in seconds to read the log file (default: 60)
  --numberlinestoread int
    	Number of bytes to read from end of the log file (default: 1024)
  --ipnumbersthreshold int
    	Number of requests from an IP to be considered as a threat (default: 10)
  --timewindow int
    	Time window in seconds to consider requests from an IP (default: 60)
  --urlpattern string
    	Pattern to match the URL. (default: all allowed)
  --linesingroup int  	
        Number of lines to group together (default: 100) 
  --outputpath string
    	Path to output file   
  --jsonlogformat bool
    	Is log in json format (default: false) 
  --outputoverwrite bool
    	Overwrite output file (default: false)   
  --DDOSER_ONLY_IPV4 bool
    	Only IPv4 addresses (default: false)  	
```

## How to build

build for current system

```
make build
```

or build for linux system

```
make build-linux
```

## How to run

```
./ddoser --logpath=/var/log/nginx/access.log --readinterval=60 --numberlinestoread=1024 --ipnumbersthreshold=10 --timewindow=60 --urlpattern=/api/v1/health --linesingroup=100 --outputpath=/tmp/ddoser.txt --jsonlogformat --outputoverwrite
```

## Contributing

If you have any ideas, just open an issue and tell what you think.