package utils

import "errors"

func FormatPhoneNumber(msisdn string) (string, error) {
	if msisdn[:2] == "62" {
		return msisdn, nil;
	}

	if msisdn[:1] == "0" {
		return "62" + msisdn[1:], nil;
	}

	if msisdn[:1] == "+" && msisdn[:3] == "+62" {
		return msisdn[1:], nil;
	}

	return "", errors.New("Nomor HP tidak dikenali");
}