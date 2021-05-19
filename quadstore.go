package quadgraph

import (
	"sync/atomic"

	"github.com/cayleygraph/quad"
)

type QuadStore struct {
	// values is a reverse hash map from the individual values to their identifiers
	values map[quad.Value]int64

	// quads is a reverse hash map from the quads direction vector to the quad identifier
	quads map[vector4]int64

	// entries is a hashmap of entries indexed by their ID
	entries map[int64]*Entry

	// the last ID we added
	lastID int64
}

const (
	EntryKindValue = iota
	EntryKindQuad
)

// Entry is an entry in the quadstore. It either consist of just a value (EntryKindValue) or
// an actual quad (EntryKindQuad), which directions are references (IDs) to value entries.
type Entry struct {
	// the ID of this entry
	ID int64

	Kind byte

	// the value of this entry
	Value quad.Value

	// the directions of this quad, i.e. references to others
	Directions vector4
}

type vector4 [4]int64

func NewStore() *QuadStore {
	return &QuadStore{
		values:  make(map[quad.Value]int64),
		entries: make(map[int64]*Entry),
	}
}

func (store *QuadStore) AddQuad(q quad.Quad) (quadID int64) {
	var directions vector4

	// loop through the individual values (called directions) of the quad
	for i := quad.Subject; i <= quad.Label; i++ {
		value := q.Get(i)

		if value == nil {
			continue
		}

		// try to map the individual value to an ID. this will either
		// return an existing one or create a new one
		valueID := store.mapValue(value)

		if valueID != 0 {
			directions[i-1] = valueID
		}
	}

	if id, ok := store.quads[directions]; ok {
		return id
	}

	return store.addEntry(&Entry{
		Directions: directions,
		Kind:       EntryKindQuad,
	})
}

func (store *QuadStore) mapValue(value quad.Value) (valueID int64) {
	var ok bool

	if valueID, ok = store.values[value]; ok {
		return valueID
	}

	// not found, need to add it
	valueID = store.addEntry(&Entry{
		Value: value,
		Kind:  EntryKindValue,
	})

	// store it reverse map
	store.values[value] = valueID

	return valueID
}

func (store *QuadStore) addEntry(entry *Entry) int64 {
	id := atomic.AddInt64(&store.lastID, 1)

	entry.ID = id

	store.entries[id] = entry

	return id
}
