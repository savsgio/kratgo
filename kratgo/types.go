package kratgo

import (
	"os"
)

// Kratgo ...
type Kratgo struct {
	Proxy Server
	Admin Server

	logFile *os.File
}

// ###### INTERFACES ######

// Server ...
type Server interface {
	ListenAndServe() error
}
