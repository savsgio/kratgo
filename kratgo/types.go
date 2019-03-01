package kratgo

import (
	"os"
)

// Kratgo ...
type Kratgo struct {
	Proxy ProxyServer
	Admin AdminServer

	logFile *os.File
}

// ###### INTERFACES ######

// ProxyServer ...
type ProxyServer interface {
	ListenAndServe() error
}

// AdminServer ...
type AdminServer interface {
	ListenAndServe() error
}
