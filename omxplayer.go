package omxplayer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/guelfey/go.dbus"
	"os"
)

const (
	prefixOmxDbusFiles = "/tmp/omxplayerdbus."
	suffixOmxDbusPid   = ".pid"
	ifaceMpris         = "org.mpris.MediaPlayer2"
)

var (
	user            string
	home            string
	fileOmxDbusPath string
	fileOmxDbusPid  string
)

func init() {
	SetUser(os.Getenv("USER"), os.Getenv("HOME"))
}

func SetUser(u, h string) {
	user = u
	home = h
	fileOmxDbusPath = prefixOmxDbusFiles + user
	fileOmxDbusPid = prefixOmxDbusFiles + user + suffixOmxDbusPid
}

// New returns a new Player instance that can be used to control an OMXPlayer
// instance that is playing the video located at the specified URL.
func New(url string) (player *Player, err error) {
	removeDbusFiles()
	return
}

// removeDbusFiles removes the files that OMXPlayer creates containing the D-Bus
// path and PID. This ensures that when the path and PID are read in, the new
// files are read instead of the old ones.
func removeDbusFiles() {
	removeFile(fileOmxDbusPath)
	removeFile(fileOmxDbusPid)
}

// getDbusPath reads the D-Bus path from the file OMXPlayer writes it's path to.
// If the file cannot be read, it returns an error, otherwise it returns the
// path as a string.
func getDbusPath() (string, error) {
	if err := waitForFile(fileOmxDbusPath); err != nil {
		return "", err
	}
	return readFile(fileOmxDbusPath)
}

// getDbusPath reads the D-Bus PID from the file OMXPlayer writes it's PID to.
// If the file cannot be read, it returns an error, otherwise it returns the
// PID as a string.
func getDbusPid() (string, error) {
	if err := waitForFile(fileOmxDbusPid); err != nil {
		return "", err
	}
	return readFile(fileOmxDbusPid)
}

// getDbusConnection establishes and returns a D-Bus connection to the specified
// D-Bus path with the specified given D-Bus PID. Since the connection's `Auth`
// method attempts to use Go's `os/user` package to get the current user's name
// and home directory, and `os/user` is not implemented for Linux-ARM, the
// `authMethods` parameter is specified explicitly rather than passing `nil`.
func getDbusConnection(path, pid string) (conn *dbus.Conn, err error) {
	authMethods := []dbus.Auth{
		dbus.AuthExternal(user),
		dbus.AuthCookieSha1(user, home),
	}

	log.Debug("omxplayer: opening dbus session")
	if conn, err = dbus.SessionBusPrivate(); err != nil {
		return
	}

	log.Debug("omxplayer: authenticating dbus session")
	if err = conn.Auth(authMethods); err != nil {
		return
	}

	log.Debug("omxplayer: initializing dbus session")
	err = conn.Hello()
	return
}
