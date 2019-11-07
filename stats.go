package vanilla

// Stats of the player
type Stats struct {
	Str  int //Strength
	Agi  int // Agility
	Sta  int // Stamina
	Inte int // Intellect
	Spi  int // Spirit
}

// Add delta stats, and return a new stats
func (stats Stats) Add(delta Stats) (res Stats) {
	res = Stats{}

	res.Str = stats.Str + delta.Str
	res.Agi = stats.Agi + delta.Agi
	res.Sta = stats.Sta + delta.Sta
	res.Inte = stats.Inte + delta.Inte
	res.Spi = stats.Spi + delta.Spi

	return
}
