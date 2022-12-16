package utils

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestReplaceEmoji(t *testing.T) {
	source := "😁🙏✂➰🚀🛀Ⓜ🉑©🗿😀😶🚁🛅🌍🕧"
	result := ReplaceEmoji(source, "[e]")
	fmt.Println(result)
	assert.Equal(t, result, "[e][e][e][e][e][e][e][e][e][e][e][e][e][e][e][e]")

	source = "122223334👌👌🤝😢😂🤣123456789"
	result = ReplaceEmoji(source, "[e]")
	fmt.Println(result)
	assert.Equal(t, result, "122223334[e][e][e][e][e][e]123456789")
}

func TestEmojiValues(t *testing.T) {
	source := "🤝🤣⌛️😁🙏✂➰🚀🛀Ⓜ🉑©🗿😀😶🚁🛅🌍🕧"
	for _, v := range source {
		fmt.Printf("%x ", v)
		fmt.Printf("%c ", v)
	}
}
