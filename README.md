# Hearthstone Deckstrings

Go library for encoding and decoding [Hearthstone deckstrings](https://hearthsim.info/docs/deckstrings/). See [documentation](https://godoc.org/github.com/schmich/deckstrings) for help.

## Usage

```bash
go get github.com/schmich/deckstrings
```

```go
import "github.com/schmich/deckstrings"
```

A Hearthstone deckstring encodes a Hearthstone deck in a compact string format. The IDs used in deckstrings
and in this library are Hearthstone DBF IDs which are unique identifiers for Hearthstone entities like cards and heroes.

For additional entity metadata (e.g. hero class, card cost, card name), DBF IDs can be used in conjunction with the [official Hearthstone API](https://develop.battle.net/documentation/api-reference/hearthstone-game-data-api) or [HearthstoneJSON](https://hearthstonejson.com) database.

See the [deckstrings.Deck](https://godoc.org/github.com/schmich/deckstrings#Deck) type for details on how a
Hearthstone deck is represented.

## Decoding

[deckstrings.Decode](https://godoc.org/github.com/schmich/deckstrings#Decode) decodes a deckstring into a Hearthstone deck.

```go
deckstring := "AAECAZICCPIF+Az5DK6rAuC7ApS9AsnHApnTAgtAX/4BxAbkCLS7Asu8As+8At2+AqDNAofOAgA="
deck, err := deckstrings.Decode(deckstring)
fmt.Printf("%+v %v", deck, err)
```

```text
{Format:2 Heroes:[274] Cards:[[64 2] [95 2] [254 2] [754 1] [836 2] [1124 2] [1656 1] [1657 1] [38318 1] [40372 2] [40416 1] [40523 2] [40527 2] [40596 1] [40797 2] [41929 1] [42656 2] [42759 2] [43417 1]]} <nil>
```

## Encoding

[deckstrings.Encode](https://godoc.org/github.com/schmich/deckstrings#Encode) encodes a Hearthstone deck
into a deckstring using [base64.StdEncoding](https://golang.org/pkg/encoding/base64/#pkg-variables).

```go
cards := [][2]uint64{
    {9, 1}, {279, 1}, {436, 1}, {545, 2}, {613, 1},
    {1363, 1}, {1367, 1}, {41169, 2}, {41176, 2},
    {42046, 1}, {42597, 2}, {42598, 1}, {42804, 2},
    {42818, 1}, {42992, 2}, {43112, 2}, {46307, 2},
    {46495, 2}, {48002, 1}, {49184, 1}, {49421, 1},
}
deck := Deck{
    Format: deckstrings.FormatStandard, // Standard
    Heroes: []uint64{41887},            // Tyrande Whisperwind
    Cards:  cards,                      // Cards in deck as (DBF ID, count) pairs
}
deckstring, err := deckstrings.Encode(deck)
fmt.Println(deckstring, err)
```

```text
AAECAZ/HAgwJlwK0A+UE0wrXCr7IAubMAsLOAoL3AqCAA42CAwmhBNHBAtjBAuXMArTOAvDPAujQAuPpAp/rAgA= <nil>
```

## License

Copyright &copy; 2018 Chris Schmich  
MIT License. See [LICENSE](LICENSE) for details.
