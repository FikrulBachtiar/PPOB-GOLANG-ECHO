package utils

import (
	"fmt"
)

func FormatPhoneNumber(msisdn string) string {
	var result string

	if msisdn[:2] == "62" {
		result = msisdn;
	} else if msisdn[:1] == "0" {
		result = fmt.Sprintf("62%s", msisdn[1:]);
	} else {
		result = fmt.Sprintf("62%s", msisdn);
	}

	return result
}