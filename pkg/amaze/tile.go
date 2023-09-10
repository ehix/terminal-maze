package amaze

import (
	"math/rand"
)

type Spinner struct {
	frames   []string
	interval int
}

var spinners = map[string]Spinner{
	"moon":  {[]string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}, 80},
	"clock": {[]string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}, 100},
	"earth": {[]string{"🌍", "🌎", "🌏 "}, 100},
}

// Constitutes each tile type used to create mazes.
type TileSet struct {
	player string
	wall   string
	empty  string
	start  string
	exit   string
	trail  string
}

// <!> Maybe need two sets, one for ambiguous chars narrow and wide.
// Alternatively, some way to control how these emojis are being used in each terminal.
var tileSets = map[string]TileSet{
	"frog":    {"🐸", "🌱", "﹏", "🪷", "🪰", "🌀"},
	"bankjob": {"🥸", "🧱", "💰", "🏦", "💎", "+1"},
	"spooky":  {"👻", "🌳", "🧍", "🪦", "🏚️ ", "😱"},
	"email":   {"📨", "🌐", "01", "💻", "📥", `10`},
	// "taxi":       {"🚖", "🏢", " '", "🙋", "🏡", "💸"},
	"veitnam": {"🪖", "🌴", "🏽", "🚁", "🛖", "🔥"},
	// "cablecar":   {"🚠", "🗻", "❄️ ", "🛕", "🏕️ ", "🐐"},
	"farmboy":   {"🚜", "🪜", "🌾", "🐮", "🥛", "🌽"},
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

// Returns a random TileSet along with its name.
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
