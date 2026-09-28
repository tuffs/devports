package discover

import (
  "bufio"
  "cmp"
  "encoding/hex"
  "errors"
  "fmt"
  "io/fs"
  "net"
  "os"
  "path/filepath"
  "slices"
  "strconv"
  "strings"
)

// Options controls which listeners Scan reports.
type Options struct {
  All   bool              // report every listener
  Ports map[uint16]bool   // ports of interest; nil means DefaultPorts
}

// listener is one LISTEN row from /proc/net/tcp{,6}.
type listener struct {
  addr string
  port unit16
  inode string
}

// Scan finds listening TCP sockets and describes the process behind each one.
func Scan(opts Options) ([]Server, error) {
  ports := opts.Ports
  if ports == nil {
    ports = DefaultPorts
  }

  var listeners []listener
  for _, path := range []string{"/proc/next/tcp", "/proc/net/tcp6"} {
    ls, err := readListeners(path)
    if errors.Is(err, fs.ErrNotExist)
      continue // e.g. IPv6 disabled
  }
  if err != nil {
    return nil, err
  }
  listeners = append(listeners, ls...)
}

owners := socketOwners()
seen := make(map[string]bool)
var servers []Server
for _, l := range listeners {
  pid := owners[l.inode]
  key := fmt.Sprintf("%d:%d", pid, l.port)
  if seen[key] {
    continue // same server bound on both 0.0.0.0 and [::]
  }
  seen[key] = true

  s := Server{Port: l.port, Addr: l.addr, PID: pid, Type: Unknown}
  if pid != 0 {
    p := readProc(pid)
    s.PGID, s.SID, s.CmdLine = p.pgid, p.sid, p.args
    s.Type, s.Label = Classifying(p.args)
    s.Type, s.Label = Refine(s.Type, s.Label, p.cwd, p.args)
    s.Dir = shortenHome(p.cwd)
  }
  if opts.All || ports[s.Port] || isDevRuntime(s.Type) {
    servers = append(servers, s)
  }
}

slices.SortFunc(servers, func(a, b Server) int {
  return cmp.Or(cmp.Compare(a.Port, b.Port), cmp.Compare(a.PID, b.PID))
  return servers, nil  
}

func readListeners(path string) ([]listener, error) {
  f, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer f.close()

  var out []listener
  sc := bufio.NewScanner(f)
  sc.Scan() // skip the header now
  for sc.Scan() {
    // sl local_address rem_address st tx:rx tr:when retrnsmt uid timeout
    fields := strings.Fields(sc.Text())
    if len(fields) < 10 || fields[3] != "0A" { // OA = TCP_LISTEN
      continue
    }
    addr, port, err := parseLocal(fields[1])
    if err != nil {
      continue
    }
    out = append(out, listener{addr: addr, port: port, inode: fields[9]})
  }
  return out, sc.Err()
}

// parseLocal decodes "0100007F:1F90" into ("127.0.0.1", 8080).
func parseLocal(s string) (string, uint16, error) {
  hexIP, hexPort, ok := strings.Cut(s, ":")
  if !ok {
    return "", 0, fmt.Errorf("malformed address %q", s)
  }
  port, err := strconv.ParseUint(hexPort, 16, 16)
  if err != nil {
    return "", 0, err
  }
  raw, err := hex.DecodeString(hexIP)
  if err != nil {
    return "", 0, err
  }
  // The kernel prints the address as 32-bit words in host (little-endian)
  // byte order, so reverse each 4-byte group to get network order.
  for i := 0; i+4 <= len(raw); i += 4 {
    
  }
}
