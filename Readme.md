# DDoser [![Build Status](https://github.com/komandakycto/ddoser/workflows/build/badge.svg)](https://github.com/komandakycto/ddoser/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/komandakycto/ddoser)](https://goreportcard.com/report/github.com/komandakycto/ddoser) [![Coverage Status](https://coveralls.io/repos/github/komandakycto/ddoser/badge.svg)](https://coveralls.io/github/komandakycto/ddoser)

DDoser is a simple tool to find attackers IP addresses in nginx log and dump it to file.

## How it works

The DDoser reads k bytes from the end of nginx log file each n seconds. After this it parses the log and finds IP
addresses of attackers.
Ip address is marked as attacker if it has more than m requests in last t seconds.

## Usage

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

```
make build
```

## How to run

```
./ddoser --logpath=/var/log/nginx/access.log --readinterval=60 --numberlinestoread=1024 --ipnumbersthreshold=10 --timewindow=60 --urlpattern=/api/v1/health --linesingroup=100 --outputpath=/tmp/ddoser.txt --jsonlogformat --outputoverwrite
```