package deckstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeSimplest(t *testing.T) {
	deckstring := "AAEAAAAAAA=="
	deck := Deck{
		Format: Format(0),
		Heroes: []uint64{},
		Cards:  [][2]uint64{},
	}

	encoded, err := Encode(deck)
	assert.Nil(t, err)
	assert.Equal(t, deckstring, encoded, "deckstrings should be equal")

	decoded, err := Decode(deckstring)
	assert.Nil(t, err)
	assert.Equal(t, deck, decoded, "decks should be equal")
}

func TestEncodeDecode(t *testing.T) {
	deckstring := "AAECAR8GxwPJBLsFmQfZB/gIDI0B2AGoArUDhwSSBe0G6wfbCe0JgQr+DAA="
	deck := Deck{
		Format: FormatStandard,
		Heroes: []uint64{31},
		Cards:  [][2]uint64{{141, 2}, {216, 2}, {296, 2}, {437, 2}, {455, 1}, {519, 2}, {585, 1}, {658, 2}, {699, 1}, {877, 2}, {921, 1}, {985, 1}, {1003, 2}, {1144, 1}, {1243, 2}, {1261, 2}, {1281, 2}, {1662, 2}},
	}

	encoded, err := Encode(deck)
	assert.Nil(t, err)
	assert.Equal(t, deckstring, encoded, "deckstrings should be equal")

	decoded, err := Decode(deckstring)
	assert.Nil(t, err)
	assert.Equal(t, deck, decoded, "decks should be equal")
}

func TestEncodeHeroSort(t *testing.T) {
	p, err := Encode(Deck{Heroes: []uint64{0, 1}, Cards: [][2]uint64{}})
	assert.Nil(t, err)

	q, err := Encode(Deck{Heroes: []uint64{1, 0}, Cards: [][2]uint64{}})
	assert.Nil(t, err)

	assert.Equal(t, p, q, "deckstrings should be equal")
}

func TestEncodeCardSort(t *testing.T) {
	p, err := Encode(Deck{Heroes: []uint64{}, Cards: [][2]uint64{{0, 1}, {1, 1}, {2, 2}, {3, 2}, {4, 3}}})
	assert.Nil(t, err)

	q, err := Encode(Deck{Heroes: []uint64{}, Cards: [][2]uint64{{3, 2}, {4, 3}, {1, 1}, {0, 1}, {2, 2}}})
	assert.Nil(t, err)

	assert.Equal(t, p, q, "deckstrings should be equal")
}

func TestEncodeHighCount(t *testing.T) {
	deck := Deck{
		Heroes: []uint64{},
		Cards:  [][2]uint64{{1, 3}, {2, 3}, {3, 3}, {4, 4}, {5, 4}, {6, 10}, {7, 100}, {8, 1000}},
	}

	encoded, err := Encode(deck)
	assert.Nil(t, err)
	assert.Equal(t, "AAEAAAAACAEDAgMDAwQEBQQGCgdkCOgH", encoded, "deckstrings should be equal")

	decoded, err := Decode(encoded)
	assert.Nil(t, err)
	assert.Equal(t, deck, decoded, "decks should be equal")
}

func TestEncodeInvalidCount(t *testing.T) {
	_, err := Encode(Deck{Heroes: []uint64{}, Cards: [][2]uint64{{10, 1}, {20, 2}, {30, 0}}})
	assert.NotNil(t, err)
}

func TestDecodeUnsortedHeroes(t *testing.T) {
	deckstring := "AAEAAgIBAAAA"
	deck := Deck{Heroes: []uint64{1, 2}, Cards: [][2]uint64{}}

	decoded, err := Decode(deckstring)
	assert.Nil(t, err)
	assert.Equal(t, deck, decoded, "decks should be equal")
}

func TestDecodeUnsortedCards(t *testing.T) {
	deckstring := "AAEAAAMDAgEAAA=="
	deck := Deck{
		Heroes: []uint64{},
		Cards:  [][2]uint64{{1, 1}, {2, 1}, {3, 1}},
	}

	decoded, err := Decode(deckstring)
	assert.Nil(t, err)
	assert.Equal(t, deck, decoded, "decks should be equal")
}

func TestDecodeEmptyDeckstring(t *testing.T) {
	_, err := Decode("")
	assert.NotNil(t, err)
}

func TestDecodeInvalidBase64(t *testing.T) {
	_, err := Decode("{}''\n\t @$%^&*()")
	assert.NotNil(t, err)
}

func TestDecodeInvalidReserved(t *testing.T) {
	_, err := Decode("BB")
	assert.NotNil(t, err)
}

func TestDecodeInvalidVersion(t *testing.T) {
	_, err := Decode("AABB")
	assert.NotNil(t, err)
}

func TestDecodeUnexpectedEOF(t *testing.T) {
	_, err := Decode("AAEB0")
	assert.NotNil(t, err)
}

func TestDeckstrings(t *testing.T) {
	deckstrings := []string{
		"AAEBAf0GAA/yAaIC3ALgBPcE+wWKBs4H2QexCMII2Q31DfoN9g4A",
		"AAECAZICCPIF+Az5DK6rAuC7ApS9AsnHApnTAgtAX/4BxAbkCLS7Asu8As+8At2+AqDNAofOAgA=",
		"AAECAaIHCLIC7QLdCJG8Asm/ApTQApziAp7iAgu0AagF1AXcrwKStgKBwgKbwgLrwgLKywKmzgKnzgIA",
		"AAECAR8E8gXtCZG8AobTAg2oArUD5QfrB5cIxQj+DLm0Auq7AuTCAo7DAtPNAtfNAgA=",
		"AAECAQcES+0FoM4Cn9MCDZAD/ASRBvgH/weyCPsMxsMC38QCzM0Cjs4Cns4C8dMCAA==",
		"AAECAf0GHjCKAZMB9wTtBfIF2waSB7YH4Qf7B40IxAjMCPMM2LsC2bwC3bwCysMC3sQC38QC08UC58sCos0C980Cn84CoM4Cws4Cl9MCl+gCAAA=",
		"AAECAZ8FDPIF9QX6Bo8JvL0C/70CucEC78ICps4Cws4CnOIC0OICCdmuArO7ApW8ApvCAsrDAuPLAqfOAvfQApboAgA=",
		"AAECAZICCEDyBfkMrqsC4LsClL0Cz8cCmdMCC1+KAf4B3gXEBuQIvq4CtLsCy7wCoM0Ch84CAA==",
		"AAEBAaIHCLIC9gTUBe0FpAeQEJG8AoHCAgu0AcsDzQObBbkGiAfdCIYJrxDEFpK2AgA=",
		"AAEBAZ8FCqcF4AX6BusPnhCEF9muArq9AuO+ArnBAgrbA6cI6g/TqgLTvAKzwQKdwgKxwgKIxwLjywIA",
		"AAEBAf0EArgI1hEOigHAAZwCyQOrBMsE5gTtBJYF+Af3DZjEAtrFArnRAgA=",
		"AAEBAa0GBgm0A5IPtxeoqwKFuAIMlwKhBNMK1wr6EaGsAui/AtHBAuXMAubMArTOAvDPAgA=",
		"AAEBAf0GCLYH+g7CD/UP8BHdvAL3zQKX0wILigGTAdMB4QeNCNwKjg6tEN4Wqa0C58sCAA==",
		"AAEBAZICCrQDxQTtBbkGig7WEegV7BWuqwLguwIKQF/+AdMDxAbkCJdovq4CoM0Ch84CAA==",
		"AAEBAQcG+QzVEbAVxsMCoM4C9s8CDEuRA9QEkQb4B/8H+wzkD4KtAszNAo7OAvHTAgA=",
		"AAEBAR8C/gyG0wIO0wG1A4cEgAfhB5cIxQjcCvcNuRHUEcsU3hbTzQIA",
		"AAECAR8CuwXFCA6oArUD6weXCNsJ7QmBCv4Mzq4C6rsC5MICjsMC080Cps4CAA==",
		"AAECAaoIBNAHiq0C9r0Cm8ICDVrvAYECgQT+BfAHkwn3qgL6qgL1rALDtAKuvAL5vwIA",
		"AAECAf0GBPcEoQaxCMUJDTDcAvUF+wXZB8IIxAi0rAL2rgLnwQKrwgLrwgKVzgIA",
		"AAECAZICCNUB/gHTAosE+wTjBdoK+QoLKUBa2AGBAqECtALgBOYFngnZCgA=",
		"AAECAQcAAZEDAA==",
		"AAECAf0EBk20AvEFigfsB5YNDCla2AG7AoUDiwOrBLQElgWABrwI2QoA",
		"AAECAZ8FAkaeCQ6EAfoBgQKhAoUDvQPcA+4EiAXjBc8GrwfQB/UMAA==",
		"AAEBAa0GAA/lBJ0GyQalCdIK0wrXCvIM8wyFEJYUiq0C7K4C0sECm8ICAA==",
		"AAEBAQcI+AeyCPkM6A+wFYawAvHTAqTnAgtLnQKQA6IE1ASRBv8H+wyCrQLMzQKOzgIA",
		"AAECAf0EBskDxQTcCum6AtDBAvbqAgzAAZUDqwSBsgKCtAKwvALBwQKYxALHxwLezQK50QLN6wIA",
		"AAECAaoICO0Fsgb7DJPBAqvnAvPnAuDqAu/3AgvuAYEE9QT+BcfBAsnHApvLArbNAp7wAqbwAu/xAgA=",
		"AAECAf0EBE1x7/EC74ADDbsClQOrBLQE5gSWBewFwcECj9MC++wC6vYClf8Cuf8CAA==",
		"AAECAZICCPIF+Az5DK6rAuC7ApS9AsnHApnTAgtAX/4BxAbkCLS7Asu8As+8At2+AqDNAofOAgA=",
	}

	for _, deckstring := range deckstrings {
		decoded, err := Decode(deckstring)
		assert.Nil(t, err)

		encoded, err := Encode(decoded)
		assert.Nil(t, err)

		assert.Equal(t, deckstring, encoded, "deckstrings should be equal")
	}
}
