package vanilla

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Stats of the player
type Stats struct {
	Str int //Strength
	Agi int // Agility
	Sta int // Stamina
	Int int // Intellect
	Spi int // Spirit
}

// ReadConfig from excel
func ReadConfig(path string) (err error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return
	}

	for index, name := range f.GetSheetMap() {
		fmt.Println(index, name)
	}

	return
}
