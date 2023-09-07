package amaze

import "math/rand"

type TileSet struct {
	player string
	wall   string
	empty  string
	start  string
	exit   string
	trail  string
}

var tileSets = map[string]TileSet{
	"human":    {"🧑", "🧱", "  ", "🚪", "🏁", "👣"},
	"frog":     {"🐸", "🌱", "﹏", "🪷 ", "🪰 ", "🌀"},
	"ghost":    {"👻", "🌳", "🧍", "🪦 ", "🏚️ ", "😱"},
	"email":    {"📨", "🌐", "01", "💻", "📥", `10`},
	"taxi":     {"🚖", "🏢", " '", "📱", "🏡", "💸"},
	"veitnam":  {"🪖 ", "🌴", "🏽", "🚁", "🛖 ", "🔥"},
	"cablecar": {"🚠", "⛰️ ", "☁️ ", "🛕", "🏕️ ", "🐐"},
	"farm":     {"🚜", "🪜 ", "🌾", "🐮", "🥛", "🌽"},
}

// 🚢

func SetRandomTiles() TileSet {
	k := rand.Intn(len(tileSets))
	i := 0
	for _, tiles := range tileSets {
		if i == k {
			return tiles
		}
		i++
	}
	panic("unreachable")
}
