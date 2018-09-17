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

const Version = 1

type Format uint64

const (
	FormatWild     Format = 1
	FormatStandard Format = 2
)

type Deck struct {
	Format Format
	Heroes []uint64
	Cards  [][2]uint64
}

func Decode(deckstring string) (_ Deck, err error) {
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

func Encode(deck Deck) (_ string, err error) {
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
