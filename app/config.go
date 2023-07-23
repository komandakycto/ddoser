package main

type Opts struct {
	LogPath             string `long:"logpath" env:"DDOSER_LOG_PATH" description:"Path to nginx access log file"`
	ReadIntervalSeconds int    `long:"readinterval" env:"DDOSER_READ_INTERVAL" default:"60" description:"Interval in seconds to read the log file"`
	BytesToRead         int64  `long:"numberlinestoread" env:"DDOSER_NUMBER_LINES_TO_READ" default:"1000" description:"Number of lines to read from end of the log file"`
	IpNumbersThreshold  int    `long:"ipnumbersthreshold" env:"DDOSER_IP_NUMBERS_THRESHOLD" default:"10" description:"Number of requests from an IP to be considered as a threat"`
	TimeWindow          int    `long:"timewindow" env:"DDOSER_TIME_WINDOW" default:"60" description:"Time window in seconds to consider requests from an IP"`
	UrlPattern          string `long:"urlpattern" env:"DDOSER_URL_PATTERN" default:"/" description:"Pattern to match the URL"`
	LinesInGroup        int    `long:"linesingroup" env:"DDOSER_LINES_IN_GROUP" default:"100" description:"Number of lines to group together"`
	OutputPath          string `long:"outputpath" env:"DDOSER_OUTPUT_PATH" description:"Path to output file"`

	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}
