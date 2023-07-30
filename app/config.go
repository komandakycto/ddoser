package main

// Opts is a struct that represents the command line options.
type Opts struct {
	LogPath             string `short:"f" long:"logpath" env:"DDOSER_LOG_PATH" description:"Path to log file to read"`
	ReadIntervalSeconds int    `short:"ri" long:"readinterval" env:"DDOSER_READ_INTERVAL" default:"60" description:"Interval in seconds to read the log file"`
	BytesToRead         int64  `short:"br" long:"numberbytestoread" env:"DDOSER_NUMBER_BYTES_TO_READ" default:"1024" description:"Number of bytes to read from end of the log file"`
	IpNumbersThreshold  int    `short:"t" long:"ipnumbersthreshold" env:"DDOSER_IP_NUMBERS_THRESHOLD" default:"10" description:"Number of requests from the same IP to be considered as a threat"`
	TimeWindow          int    `short:"tw" long:"timewindow" env:"DDOSER_TIME_WINDOW" default:"60" description:"Time window in seconds to count requests from the same IP address"`
	UrlPattern          string `short:"up" long:"urlpattern" env:"DDOSER_URL_PATTERN" default:"" description:"Pattern to match the URL"`
	LinesInGroup        int    `short:"lg" long:"linesingroup" env:"DDOSER_LINES_IN_GROUP" default:"100" description:"Parsed line will be grouped into groups of this size to process them concurrently"`
	OutputPath          string `short:"o" long:"outputpath" env:"DDOSER_OUTPUT_PATH" description:"Path to output file with attackers IP addresses"`
	AverageLogBytes     int64  `short:"alb" long:"averagelogbytes" env:"DDOSER_AVERAGE_LOG_BYTES" default:"157" description:"Average log line size in bytes"`
	JsonLogTimeLayout   string `long:"jsonlogtimelayout" env:"DDOSER_JSON_LOG_TIME_LAYOUT" default:"" description:"Time layout for JSON log format"`
	JsonLogFormat       bool   `long:"jsonlogformat" env:"DDOSER_JSON_LOG_FORMAT" description:"Use JSON log format"`
	OutputOverwrite     bool   `long:"outputoverwrite" env:"DDOSER_OUTPUT_OVERWRITE" description:"Overwrite the output file. Append by default."`
	OnlyIpv4            bool   `long:"onlyipv4" env:"DDOSER_ONLY_IPV4" description:"Process only IPv4 addresses"`
}
