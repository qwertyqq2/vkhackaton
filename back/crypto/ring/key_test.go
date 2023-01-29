package ring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	pk := GeneratePrivate()
	pub := pk.Public()
	fmt.Println(pk.String())
	fmt.Println(pub.String())
}

func TestParse(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		pk := GeneratePrivate()
		pub := pk.Public()
		pubs := pub.String()
		pks := pk.String()
		pkCopy := ParsePrivate(pks)
		pubCopy := ParsePublic(pubs)
		assert.Equal(t, pks, pkCopy.String())
		assert.Equal(t, pubs, pubCopy.String())
	})

}
