package middleware

import "sync"


var(
  blacklist = make(map[string]bool)
  mu sync.Mutex
)


func AddToBlacklist(token string) {
  mu.Lock()
  defer mu.Unlock()
  blacklist[token] = true 
}

func IsTokenBlacklisted(token string) bool {
  mu.Lock()
  defer mu.Unlock()
  return blacklist[token]
}
