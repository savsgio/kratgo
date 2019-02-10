package kratgo

import (
	"os"

	"github.com/savsgio/kratgo/internal/admin"
	"github.com/savsgio/kratgo/internal/proxy"
)

// Kratgo ...
type Kratgo struct {
	Proxy *proxy.Proxy
	Admin *admin.Admin

	logFile *os.File
}
