package main
import (
	"strings"
)

func get_suffix_rune(url *string) (suffix string, runes *[]rune) {
	if strings.HasPrefix(*url, "comics_adventure/th/") {
		suffix = ".webp"
		runes = &letterRunes
	} else if strings.HasPrefix(*url, "comics_events/th/") || strings.HasPrefix(*url, "comics_mythic/th/") {
		suffix = ".jpg"
		runes = &letterRunes
	} else if strings.HasPrefix(*url, "comics/th/") {
		suffix = "@2x.webp"
		runes = &letternRunes
	}
	return
}