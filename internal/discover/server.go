// TCP servers and identifies what's running them
package discover 

// AppType broad program family, uses listenig socket 
type AppType string

const (
  Node      AppType = "Node"
  PHP       AppType = "PHP"
  Spring    AppType = "Spring"
  Java      AppType = "Java"
  Rails     AppType = "Rails"
  Ruby      AppType = "Ruby"
  Python    AppType = "Python"
  WebServer AppType = "Web Server"
)

// Info about the targeted Process
type Server struct {
  Port    uint16   // 
  Addr    string   // "127.0.0.1", etc.
  PID     int      // 0 for another user's
  PGID    int      // process group ID
  SID     int      // session ID
  Type    AppType  // Node, PHP, Rails, etc 
  Label   string   // "Vite', "Laravel", Puma"
  Dir     string   // working directory of server
  Cmdline []string // argv
}

// Display returns "Type (Label)" or just "Type".
func (s Server) Display() string {
  if s.Label == "" {
    return string(s.Type)
  }
  return string(s.Type) + " (" + s.Label + ")"
}

