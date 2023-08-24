package pack

import (
	"sort"
	"sync"

	"github.com/meta-node-blockchain/meta-node/types"
)

type PackPool struct {
	mu    sync.Mutex
	packs []types.Pack
}

func NewPackPool() *PackPool {
	return &PackPool{
		packs: make([]types.Pack, 0),
	}
}

func (p *PackPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.packs)
}

func (p *PackPool) AddPack(pack types.Pack) {
	p.mu.Lock()
	p.packs = append(p.packs, pack)
	p.mu.Unlock()
}

func (p *PackPool) AddPacks(packs []types.Pack) {
	p.mu.Lock()
	p.packs = append(p.packs, packs...)
	p.mu.Unlock()
}

func (p *PackPool) TakePack(numberOfPack uint64) []types.Pack {
	p.mu.Lock()
	sort.Slice(p.packs, func(i int, u int) bool {
		return p.packs[i].Timestamp() < p.packs[u].Timestamp()
	})
	if int(numberOfPack) > len(p.packs) {
		numberOfPack = uint64(len(p.packs))
	}
	rs := p.packs[:numberOfPack]
	p.packs = p.packs[numberOfPack:]
	p.mu.Unlock()
	return rs
}

func (p *PackPool) Copy() types.PackPool {
	rs := NewPackPool()
	rs.packs = append(rs.packs, p.packs...)
	return rs
}
