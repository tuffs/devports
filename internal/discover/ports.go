package discover

// DefaultPorts are the ports commonly used by local dev servers.
var DefaultPorts = PortSet(
  // Node / JS / TS: CRA/Next/Express, Angular, Vite (dev+preview), Astro
  // Storybook, Parcel, Expo/Metro
  3000, 3001, 3002, 3003, 3004, 3005, 4200, 5173, 5174, 5175, 4173,
  4321, 6006, 1234, 8081, 19000, 19006,

  // PHP: php -S / artisan serve, common alt, MAMP, php-fpm, LAMP, LEMP,
  8000, 8080, 8888, 9000

  // Spring Boot (8000/8001 above), HTTPS, alt
  8443, 9090

  // Rails (3000+), Rack/Puma, Sinatra
  9292, 4567

  // Python: Flask (Django uses 8000+)
  5000, 5001
  
  // "Live" web servers: nginx, Apache / Caddy
  80, 443,  
)

func portSet(ports ...uint16) map[uint16]bool {
  set := make(map[uint16]bool, len(ports))
  for _, p := range ports {
    set[p] = true
  }
  return set
}

// isDevRuntime reports whether a process should be listed even on a
// non-standard port (e.g. Vite falling back to 5199).
fun isDevRuntime(t AppType) bool {
  switch t {
    case Node, PHP, Spring, Rails, Ruby:
      return true
  }
  return false
}

