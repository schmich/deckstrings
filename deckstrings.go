// This package encodes and decodes Hearthstone deckstrings.
//
// A Hearthstone deckstring encodes a Hearthstone deck in a compact string format.
// The IDs used in deckstrings and in this library are Hearthstone DBF IDs
// which are unique identifiers for Hearthstone entities like cards and heroes.
//
// For additional entity metadata (e.g. hero class, card cost, card name), DBF IDs
// can be used in conjunction with the HearthstoneJSON database. See
// https://hearthstonejson.com/ for details.
package deckstrings

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// The deckstring version supported by this package. Decoding a deckstring
// with a newer version is not supported. All deckstrings encoded by this
// package include this version.
const Version = 1

// The game format for which the deck was built. Wild and Standard are the current
// Hearthstone game formats.
type Format uint64

const (
	FormatWild     Format = 1
	FormatStandard Format = 2
)

// Deck represents a Hearthstone deck with its associated game format, hero,
// and card inventory.
//
// The Format field will typically be FormatWild or FormatStandard. Since Format
// is just a type alias for uint64, however, any uint64 value can be encoded to
// or decoded from a deckstring.
//
// The Heroes field is an array of hero DBF IDs for whom this deck was built.
// While multiple heroes (or no heroes) can be associated with a deck, Hearthstone
// does not currently support such a concept, so this will typically have just a
// single value.
//
// The Heroes field refers to specific characters (e.g. Malfurion or Lunara), not the
// general class (e.g. Druid). You can use metadata from HearthstoneJSON to map
// from the individual hero to the deck's class.
//
// The Cards field is an inventory of the cards present in the deck. It's an array
// of uint64 pairs with the first element being the card's unique DBF ID and the
// second element being the count of that card in the deck (typically 1 or 2). A
// count of 0 is invalid. Counts greater than 2 are valid but are typically not
// seen in Hearthstone decks. The count of cards will typically sum to 30, but a
// deckstring can encode an arbitrary number of cards.
//
// See HearthstoneJSON for hero and card metadata using DBF IDs:
// https://hearthstonejson.com/
type Deck struct {
	Format Format
	Heroes []uint64
	Cards  [][2]uint64
}

// Decode a deckstring into a Hearthstone deck.
//
// Decodings are canonical: the resulting deck's Heroes and Cards fields are
// ordered by DBF ID ascending.
//
// Returns an error if the string is not base64 encoded, if the deckstring version
// is not supported, or if the general format is invalid. See the Deck type for
// details about possible values and ranges for format, heroes, and cards.
func Decode(deckstring string) (deck Deck, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "deckstring decode")
		}
	}()

	reader := bufio.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(deckstring)))
	varint := &varintReader{reader}

	header := [4]uint64{}
	if err = varint.ReadMany(header[:]); err != nil {
		return Deck{}, err
	}

	if reserved := header[0]; reserved != 0 {
		return Deck{}, fmt.Errorf("unexpected reserved byte: %d", reserved)
	}

	if version := header[1]; version != Version {
		return Deck{}, fmt.Errorf("unsupported version: %d", version)
	}

	format, length := header[2], header[3]

	heroes := make([]uint64, length)
	for i := uint64(0); i < length; i++ {
		hero, err := varint.Read()
		if err != nil {
			return Deck{}, err
		}

		heroes[i] = hero
	}

	// Sort heroes.
	sort.Slice(heroes, func(i, j int) bool { return heroes[i] < heroes[j] })

	cards := make([][2]uint64, 0, 30)

	for group := 1; group <= 3; group++ {
		var err error
		var length uint64
		if length, err = varint.Read(); err != nil {
			return Deck{}, err
		}

		for i := uint64(0); i < length; i++ {
			dbfID, err := varint.Read()
			if err != nil {
				return Deck{}, err
			}

			count := uint64(group)
			if group >= 3 {
				if count, err = varint.Read(); err != nil {
					return Deck{}, err
				}
			}

			card := [2]uint64{dbfID, count}
			cards = append(cards, card)
		}
	}

	// Sort cards by DBF ID.
	sort.Slice(cards, func(i, j int) bool { return cards[i][0] < cards[j][0] })

	return Deck{
		Format: Format(format),
		Heroes: heroes,
		Cards:  cards,
	}, nil
}

// Encode a Hearthstone deck into a deckstring using base64.StdEncoding.
//
// Encodings are canonical: the deck's Heroes and Cards fields are encoded
// in ascending DBF ID order.
//
// Returns an error if any card count is 0. See the Deck type for details
// about possible values and ranges for format, heroes, and cards.
func Encode(deck Deck) (deckstring string, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "deckstring encode")
		}
	}()

	var buf bytes.Buffer
	writer := base64.NewEncoder(base64.StdEncoding, &buf)
	varint := &varintWriter{writer}

	values := []uint64{
		0,       // Reserved. Must be zero.
		Version, // Deckstring encoding version.
		uint64(deck.Format),
		uint64(len(deck.Heroes)),
	}

	if err = varint.WriteMany(values); err != nil {
		return "", err
	}

	// Sort heroes.
	heroes := make([]uint64, len(deck.Heroes))
	copy(heroes, deck.Heroes)
	sort.Slice(heroes, func(i, j int) bool { return heroes[i] < heroes[j] })

	if err = varint.WriteMany(heroes); err != nil {
		return "", err
	}

	// Gather cards into groups based on their count in the deck.
	// There are only three groups: 1x cards, 2x cards, and any other multiple.
	groups := make(map[int][][2]uint64)
	for _, card := range deck.Cards {
		dbfID, count := card[0], card[1]
		if count < 1 {
			return "", fmt.Errorf("invalid card count for DBF ID %d", dbfID)
		}

		groupID := 3
		if count < 3 {
			groupID = int(count)
		}

		if _, ok := groups[groupID]; !ok {
			groups[groupID] = [][2]uint64{card}
		} else {
			groups[groupID] = append(groups[groupID], card)
		}
	}

	for groupID := 1; groupID <= 3; groupID++ {
		var ok bool
		var group [][2]uint64
		if group, ok = groups[groupID]; !ok {
			group = [][2]uint64{}
		}

		// Sort group by card DBF ID.
		sort.Slice(group, func(i, j int) bool { return group[i][0] < group[j][0] })

		if err = varint.Write(uint64(len(group))); err != nil {
			return "", err
		}

		for _, card := range group {
			dbfID, count := card[0], card[1]
			if err = varint.Write(dbfID); err != nil {
				return "", err
			}

			// For cards with unusual counts (e.g. not 1x or 2x),
			// we write an explicit count as well.
			if groupID == 3 {
				if err = varint.Write(count); err != nil {
					return "", err
				}
			}
		}
	}

	if err = writer.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
