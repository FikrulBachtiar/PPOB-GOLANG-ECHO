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

func RandomNumber(length int) int {
	rand.Seed(time.Now().UnixNano());
	var numbers int
	var listNumbers []int
	var result int
	zero := 1;

	for x := 1; x < length; x++ {
		randomNumber := rand.Intn(9);
		zero *= 10;
		numbers = zero * randomNumber;
		listNumbers = append(listNumbers, numbers);
	}

	for _, num := range listNumbers {
		result += num;
	}

	return result;
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