package converts

import (
	"regexp"
	"strconv"
	"strings"
)

func StringTouint(input string) (uint, error) {
	if input == "" {
		return 0, nil
	}

	output, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(output), nil
}

func StringToArrUint(data string) ([]uint, error) {

	dataWithoutSpace := regexp.
		MustCompile(" ").
		ReplaceAllLiteralString(data, "")

	info := strings.Split(dataWithoutSpace, ",")

	uints := []uint{}

	for _, data := range info {
		if data == "" {
			continue
		}
		num, err := StringTouint(data)
		if err != nil {
			return nil, err
		}
		uints = append(uints, num)
	}
	return uints, nil
}

func StringToFloat32(input string) (float32, error) {
	if input == "" {
		return 0, nil
	}

	output, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, err
	}

	return float32(output), nil
}

func StringToUint(input string) (*uint, error) {
	if input == "" {
		return nil, nil
	}

	output, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return nil, err
	}

	value2 := uint(output)
	var output2 *uint = &value2

	return output2, nil
}

func StringToUint64(input string) (uint64, error) {
	if input == "" {
		return 0, nil
	}

	output, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, err
	}

	return output, nil
}
