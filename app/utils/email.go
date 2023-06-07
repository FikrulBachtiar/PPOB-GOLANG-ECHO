package utils

func EmailParsingString(value bool) string {
	var emailStatus string
	if value {
		emailStatus = "VERIFIED"
	} else {
		emailStatus = "UNVERIFIED"
	}

	return emailStatus
}