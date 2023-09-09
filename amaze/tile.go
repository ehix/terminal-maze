package amaze

import (
	"math/rand"
)

type TileSet struct {
	player string
	wall   string
	empty  string
	start  string
	exit   string
	trail  string
}

type Spinner struct {
	frames   []string
	interval int
}

var spinners = map[string]Spinner{
	"moon":  {[]string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}, 80},
	"clock": {[]string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}, 100},
	"earth": {[]string{"🌍", "🌎", "🌏 "}, 100},
}

// <?> Maybe need two sets, one for ambiguous chars narrow and wide.
var tileSets = map[string]TileSet{
	"frog":  {"🐸", "🌱", "﹏", "🪷", "🪰", "🌀"},
	"bank":  {"🥸", "🧱", "💰", "🏦", "💎", "+1"},
	"ghost": {"👻", "🌳", "🧍", "🪦", "🏚️ ", "😱"},
	"email": {"📨", "🌐", "01", "💻", "📥", `10`},
	// "taxi":       {"🚖", "🏢", " '", "🙋", "🏡", "💸"},
	"veitnam": {"🪖", "🌴", "🏽", "🚁", "🛖", "🔥"},
	// "cablecar":   {"🚠", "🗻", "❄️ ", "🛕", "🏕️ ", "🐐"},
	"farm":      {"🚜", "🪜", "🌾", "🐮", "🥛", "🌽"},
	"lunchtime": {"👷", "🚧", "🌳", "🏗️ ", "🧰", "🏠"},
	// "briefcase2": {"👔", "👤", "📈", "🏢", "💼", "🤝"},
	// "briefcase":  {"👔", "💲", "👤", "💼", "📈", "🤝"},
	"whaling": {"🛳 ", "🪨", "﹏", "🏝️ ", "🐳", "🌊"},
	// "composer":   {"🎻", "💭", "🪶", "🎼", "🙌", "🎶"},
	// "damn":   {"🦫", "🌲", "💧", "🏔️ ", "👌", "🪵"},
	"pompei": {"🔥", "🏺", "🏛 ", "🌋", "😐", "💥"},
	// "honey":       {"🐻", "🌲", "🐝", "🏔️ ", "🍯", "🩹"},
	"toothfairy": {"🧚", "🛌", "🦷", "✨", "👄", "🪙"},
	// "nightout":    {"🕺", "💃", "🍺", "😎", "🌯‍", "🫗"},
	// "smokingarea": {"🤢", "🌬️ ", "🚬", "🚭", "🍃", "💨"},
	// "pilgrim":     {"😇", "🍀", "🐍", "🕍", "🍺", "💀"},
	// "carchase":  {"🚘", "🏢", "' ", "🏦", "✈️ ", "🚔"},
	// "kingmaker": {"🗡 ", "🌳", "🛡️ ", "🏰", "🤴", "🔥"},
}

func SetRandomTiles() (string, TileSet) {
	k := rand.Intn(len(tileSets))
	i := 0
	for key, tiles := range tileSets {
		if i == k {
			return key, tiles
		}
		i++
	}
	panic("unreachable")
}
