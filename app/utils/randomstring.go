package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RandomString(length int) (string, error) {
	randomBytes := make([]byte, length);

	_, err := rand.Read(randomBytes);
	if err != nil {
		return "", err;
	}

	regex := regexp.MustCompile("[^a-zA-Z0-9]+");

	randomString := base64.URLEncoding.EncodeToString(randomBytes);
	randomString = strings.ToUpper(randomString);
	randomString = randomString[:length];
	randomString = regex.ReplaceAllString(randomString, "");

	return randomString, nil;
}

func RandomNumber(number int) int {
	rand.Seed(time.Now().UnixNano());
	randomNumber := rand.Intn(number);
	return randomNumber;
}

func SerialNumberString(number int) (string, error) {
	numb := "";
	
	numbString := strconv.Itoa(number);
	length := len(numbString);

	zeroLength := 10 - length;
	var zero []string;

	for s := 0; s < zeroLength; s++ {
		zero = append(zero, "0");
	}

	numb = fmt.Sprintf("%s%s", strings.Join(zero, ""), numbString);

	return numb, nil;
}