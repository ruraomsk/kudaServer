package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath string `toml:"logpath"`
	Port    int    `toml:"port"`
}

func init() {
}
