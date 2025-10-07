package app

import (
	"errors"
	"fmt"

	uqdefines "github.com/Yeah114/FunShuttler/uqholder/defines"
	"github.com/Yeah114/FunShuttler/utils/sync_wrapper"
)

var (
	// ErrPlayerUUIDUnknown indicates the UQ holder has no UUID information for the player.
	ErrPlayerUUIDUnknown = errors.New("player uuid unavailable")
)

// PlayerRegistry keeps a bind of UUID strings to player readers.
type PlayerRegistry struct {
	store *sync_wrapper.SyncKVMap[string, uqdefines.PlayerUQReader]
}

// NewPlayerRegistry constructs an empty registry.
func NewPlayerRegistry() *PlayerRegistry {
	return &PlayerRegistry{
		store: sync_wrapper.NewSyncKVMap[string, uqdefines.PlayerUQReader](),
	}
}

// Bind registers the provided reader and returns its UUID string.
func (r *PlayerRegistry) Bind(reader uqdefines.PlayerUQReader) (string, error) {
	if reader == nil {
		return "", fmt.Errorf("player reader is nil")
	}
	uuidStr, ok := reader.GetUUIDString()
	if !ok || uuidStr == "" {
		return "", ErrPlayerUUIDUnknown
	}
	r.store.Set(uuidStr, reader)
	return uuidStr, nil
}

// Rebind stores the reader using the supplied uuidStr.
func (r *PlayerRegistry) Rebind(uuidStr string, reader uqdefines.PlayerUQReader) {
	if uuidStr == "" || reader == nil {
		return
	}
	r.store.Set(uuidStr, reader)
}

// Get fetches a player reader by UUID string.
func (r *PlayerRegistry) Get(uuidStr string) (uqdefines.PlayerUQReader, bool) {
	if uuidStr == "" {
		return nil, false
	}
	return r.store.Get(uuidStr)
}

// Delete removes a binding by UUID string.
func (r *PlayerRegistry) Delete(uuidStr string) {
	if uuidStr == "" {
		return
	}
	r.store.Delete(uuidStr)
}

// Iterate walks over all stored players until fn returns false.
func (r *PlayerRegistry) Iterate(fn func(uuid string, player uqdefines.PlayerUQReader) bool) {
	r.store.Iter(func(k string, v uqdefines.PlayerUQReader) (cont bool) {
		return fn(k, v)
	})
}
