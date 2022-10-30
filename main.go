package main

import (
	"encoding/json"
	"fmt"

	"github.com/ebfe/scard"
	"github.com/gogetth/sscard"
)

var err error
var context *scard.Context

var reader string
var readers []string

var card *scard.Card

type IDCard struct {
	CiticenID string `json:"citizen_id"`
	Fullname  string `json:"fullname"`
	Picture   []byte `json:"picture"`
}

func onError() {
	if r := recover(); r != nil {
		b, _ := json.Marshal(map[string]interface{}{
			"success": false,
			"message": r,
			"reader":  reader,
			"readers": readers,
		})
		fmt.Println(string(b))
		// fmt.Println(base64.RawURLEncoding.EncodeToString(b))
	}
}

func onSuccess(data *IDCard) {
	b, _ := json.Marshal(map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    data,
		"reader":  reader,
		"readers": readers,
	})
	fmt.Println(string(b))
	// fmt.Println(base64.RawURLEncoding.EncodeToString(b))
}

func main() {
	defer onError()

	context, err = scard.EstablishContext()

	if err != nil {
		panic(err.Error())
	}

	// Release the PC/SC context (when needed)
	defer context.Release()

	// List available readers
	readers, err = context.ListReaders()

	if err != nil {
		panic(err.Error())
	}

	// Use the first reader
	reader = readers[0]
	// fmt.Println("Using reader:", reader)

	// Connect to the card
	card, err = context.Connect(reader, scard.ShareShared, scard.ProtocolAny)

	if err != nil {
		panic(err.Error())
	}

	// Disconnect (when needed)
	defer card.Disconnect(scard.LeaveCard)

	_, err = sscard.APDUGetRsp(card, sscard.APDUThaiIDCardSelect)
	if err != nil {
		fmt.Println("Error Transmit:", err)
		return
	}

	cardPhotoJpg, err := sscard.APDUGetBlockRsp(card, sscard.APDUThaiIDCardPhoto, sscard.APDUThaiIDCardPhotoRsp)

	if err != nil {
		panic(err.Error())
	}

	onSuccess(&IDCard{Picture: cardPhotoJpg})
}
