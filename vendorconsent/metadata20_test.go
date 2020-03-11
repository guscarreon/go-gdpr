package vendorconsent

import (
	"testing"
	"time"
)

func TestCreatedDate20(t *testing.T) {
	consent, err := Parse20(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA"))
	assertNilError(t, err)
	created := consent.Created().UTC()
	year, month, day := created.Date()
	assertIntsEqual(t, 2020, year)
	assertIntsEqual(t, int(time.February), int(month))
	assertIntsEqual(t, 27, day)
	assertIntsEqual(t, 19, created.Hour())
	assertIntsEqual(t, 51, created.Minute())
	assertIntsEqual(t, 49, created.Second())
}

func TestLastUpdated20(t *testing.T) {
	consent, err := Parse20(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA"))
	assertNilError(t, err)
	updated := consent.LastUpdated().UTC()
	year, month, day := updated.Date()
	assertIntsEqual(t, 2020, year)
	assertIntsEqual(t, int(time.February), int(month))
	assertIntsEqual(t, 27, day)
	assertIntsEqual(t, 19, updated.Hour())
	assertIntsEqual(t, 51, updated.Minute())
	assertIntsEqual(t, 49, updated.Second())
}

func TestLargeCmpID20(t *testing.T) {
	consent, err := Parse20(decode(t, "COv_46cOv_46cFZFZTENAPCAAAAAAAAAAAAAE5QBwABAAXABVAH8AgAElgJkATkAYEAgAAQACAAGAAXABUAH8AQIAwAAAA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 345, consent.CmpID())
}

func TestLargeCmpVersion20(t *testing.T) {
	consent, err := Parse20(decode(t, "COv_46cOv_46cFZFZTENAPCAAAAAAAAAAAAAE5QBwABAAXABVAH8AgAElgJkATkAYEAgAAQACAAGAAXABUAH8AQIAwAAAA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 345, consent.CmpVersion())
}

func TestLargeConsentScreen20(t *testing.T) {
	consent, err := Parse20(decode(t, "COv_46cOv_46cFZFZTENAPCAAAAAAAAAAAAAE5QBwABAAXABVAH8AgAElgJkATkAYEAgAAQACAAGAAXABUAH8AQIAwAAAA"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 19, consent.ConsentScreen())
}

func TestLanguageExtremes20(t *testing.T) {
	consent, err := Parse20(decode(t, "COv_46cOv_46cFZFZTBGAPCAAAAAAAAAAAAAE5QBwABAAXABVAH8AgAElgJkATkAYEAgAAQACAAGAAXABUAH8AQIAwAAAA"))
	assertNilError(t, err)
	assertStringsEqual(t, "BG", consent.ConsentLanguage())

	consent, err = Parse20(decode(t, "COv_46cOv_46cFZFZTSVAPCAAAAAAAAAAAAAE5QBwABAAXABVAH8AgAElgJkATkAYEAgAAQACAAGAAXABUAH8AQIAwAAAA"))
	assertNilError(t, err)
	assertStringsEqual(t, "SV", consent.ConsentLanguage())
}
