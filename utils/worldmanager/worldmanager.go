package worldmanager

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/df-mc/goleveldb/leveldb/opt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

// WorldManager manages multiple worlds, dragonfly does not have multi-world management itself,
// so we must implement it ourselves.
type WorldManager struct {
	s *server.Server

	folderPath string

	log *logrus.Logger

	worldsMu sync.RWMutex
	worlds   map[string]*world.World
}

// New ...
func New(server *server.Server, folderPath string, log *logrus.Logger) *WorldManager {
	_ = os.Mkdir(folderPath, 0644)
	defaultWorld := server.World()
	return &WorldManager{
		s:          server,
		folderPath: folderPath,
		log:        log,
		worlds: map[string]*world.World{
			defaultWorld.Name(): defaultWorld,
		},
	}
}

// DefaultWorld ...
func (m *WorldManager) DefaultWorld() *world.World {
	return m.s.World()
}

// Worlds ...
func (m *WorldManager) Worlds() []*world.World {
	m.worldsMu.RLock()
	worlds := make([]*world.World, 0, len(m.worlds))
	for _, w := range m.worlds {
		worlds = append(worlds, w)
	}
	m.worldsMu.RUnlock()
	return worlds
}

// World ...
func (m *WorldManager) World(name string) (*world.World, bool) {
	m.worldsMu.RLock()
	w, ok := m.worlds[strings.ToLower(name)]
	m.worldsMu.RUnlock()
	return w, ok
}

// LoadWorld ...
func (m *WorldManager) LoadWorld(worldName string, dimension world.Dimension, generator world.Generator) error {
	if _, ok := m.World(worldName); ok {
		return fmt.Errorf("world is already loaded")
	}

	p, err := mcdb.New(m.log, m.folderPath+"/"+worldName, opt.DefaultCompressionType)
	if err != nil {
		return fmt.Errorf("error loading world: %v", err)
	}
	p.SaveSettings(&world.Settings{
		Name:            worldName,
		TimeCycle:       false,
		WeatherCycle:    true,
		DefaultGameMode: world.GameModeSurvival,
		Difficulty:      world.DifficultyNormal,
	})

	readonly := true
	if strings.EqualFold(worldName, "god") {
		readonly = false
	}

	w := world.Config{
		Log:             m.log,
		Dim:             dimension,
		Provider:        p,
		Generator:       generator,
		ReadOnly:        readonly,
		RandomTickSpeed: -1,
	}.New()

	if _, ok := generator.(world.NopGenerator); ok {
		w.SetBlock(cube.Pos{0, 0, 0}, block.Grass{}, nil)
	}

	m.worldsMu.Lock()
	m.worlds[worldName] = w
	m.worldsMu.Unlock()
	return nil
}

// UnloadWorld ...
func (m *WorldManager) UnloadWorld(w *world.World) error {
	if w == m.DefaultWorld() {
		return fmt.Errorf("the default world cannot be unloaded")
	}

	if _, ok := m.World(w.Name()); !ok {
		return fmt.Errorf("world isn't loaded")
	}

	m.log.Debugf("Unloading world '%v'\n", w.Name())
	for _, p := range m.s.Players() {
		if p.World() == w {
			m.DefaultWorld().AddEntity(p)
			p.Teleport(m.DefaultWorld().Spawn().Vec3Middle())
		}
	}

	m.worldsMu.Lock()
	delete(m.worlds, w.Name())
	m.worldsMu.Unlock()

	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing world: %v", err)
	}
	m.log.Debugf("Unloaded world '%v'\n", w.Name())
	return nil
}

// Close ...
func (m *WorldManager) Close() error {
	m.worldsMu.Lock()
	for _, w := range m.worlds {
		// Let dragonfly close this.
		if w == m.DefaultWorld() {
			continue
		}

		m.log.Debugf("Closing world '%v'\n", w.Name())
		if err := w.Close(); err != nil {
			return err
		}
	}
	m.worlds = nil
	m.worldsMu.Unlock()
	return nil
}
