package main

import (
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/Roverr/hotstreak"
	"github.com/natefinch/lumberjack"
	"github.com/riltech/streamer"
)

type IProcess interface {
	Spawn(path, URI string) *exec.Cmd
}
type CORS struct {
	Enabled          bool     `envconfig:"CORS_ENABLED" default:"false"`           // Indicates if cors should be handled as configured or as default
	AllowedOrigins   []string `envconfig:"CORS_ALLOWED_ORIGINS" default:""`        // A list of origins a cross-domain request can be executed from.
	AllowCredentials bool     `envconfig:"CORS_ALLOW_CREDENTIALS" default:"false"` // Indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
	MaxAge           int      `envconfig:"CORS_MAX_AGE" default:"0"`               // Indicates how long (in seconds) the results of a preflight request can be cached.
}
type ProcessLoggingOpts struct {
	Enabled    bool   // Option to set logging for transcoding processes
	Directory  string // Directory for the logs
	MaxSize    int    // Maximum size of kept logging files in megabytes
	MaxBackups int    // Maximum number of old log files to retain
	MaxAge     int    // Maximum number of days to retain an old log file.
	Compress   bool   // Indicates if the log rotation should compress the log files
}
type Stream struct {
	ID          string               `json:"id"`
	Path        string               `json:"path"`
	Running     bool                 `json:"running"`
	CMD         *exec.Cmd            `json:"-"`
	Process     IProcess             `json:"-"`
	Mux         *sync.Mutex          `json:"-"`
	Streak      *hotstreak.Hotstreak `json:"-"`
	OriginalURI string               `json:"-"`
	StorePath   string               `json:"-"`
	KeepFiles   bool                 `json:"-"`
	LoggingOpts *ProcessLoggingOpts  `json:"-"`
	Logger      *lumberjack.Logger   `json:"-"`
	WaitTimeOut time.Duration        `json:"-"`
}

type StreamDTO struct {
	URI   string `json:"uri"`
	Alias string `json:"alias"`
}
type ApiServer struct {
	PortNumber string
	Stream     []streamer.Stream
	Server     *http.Server
	Spec       *Specification
}

type Specification struct {
	Debug bool `envconfig:"DEBUG" default:"false"` // Indicates if debug log should be enabled or not
	Port  int  `envconfig:"PORT" default:"8080"`   // Port that the application listens on
	EndpointYML
}
type EndpointSetting struct {
	Enabled bool   `yml:"enabled"`
	Secret  string `yml:"secret"`
}
type ListenSetting struct {
	Enabled    bool   `yaml:"enabled"`
	Uri        string `yaml:"uri"`
	Alias      string `yaml:"alias"`
	MacAddress string `yaml:"macAddress"`
}
type EndpointYML struct {
	Version string `yaml:"version"`

	Endpoints struct {
		Start  EndpointSetting `yaml:"start"`
		Stop   EndpointSetting `yaml:"stop"`
		List   EndpointSetting `yaml:"list"`
		Static EndpointSetting `yaml:"static"`
	} `yaml:"endpoints"`
	Listen []ListenSetting `yaml:"listen"`
}

type FolderInfo struct {
	FolderName string `json:"foldername"`
	MacAdd     string `json:"macadress"`
	Count      string `json:"count"`
	Url        string `json:"url"`
}
