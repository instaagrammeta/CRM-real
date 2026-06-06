package auth

import "time"

// TTL of issued tokens.
func (m *Manager) TTL() time.Duration { return m.ttl }
