package utils

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestReplaceEmoji(t *testing.T) {
	source := "ππββ°ππβπΒ©πΏππΆππππ§"
	result := ReplaceEmoji(source, "[e]")
	fmt.Println(result)
	assert.Equal(t, result, "[e][e][e][e][e][e][e][e][e][e][e][e][e][e][e][e]")

	source = "122223334πππ€π’ππ€£123456789"
	result = ReplaceEmoji(source, "[e]")
	fmt.Println(result)
	assert.Equal(t, result, "122223334[e][e][e][e][e][e]123456789")
}

func TestEmojiValues(t *testing.T) {
	source := "π€π€£βοΈππββ°ππβπΒ©πΏππΆππππ§"
	for _, v := range source {
		fmt.Printf("%x ", v)
		fmt.Printf("%c ", v)
	}
}
