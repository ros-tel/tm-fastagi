package main

import (
	"net/url"
	"regexp"
)

// Parse AGI reguest return path and query params
func parseAgiReq(request string) (string, url.Values, error) {
	req, err := url.Parse(request)
	if err != nil {
		return "", nil, err
	}
	query, err := url.ParseQuery(req.RawQuery)
	if err != nil {
		return "", nil, err
	}
	return req.Path, query, nil
}

func carNumberToSay(car_num string) string {
	result := ""

	var re = regexp.MustCompile(`(?i)^(?:[^\d]*)?(\d)(\d)(\d)(?:[^\d]*)?`)
	if len(re.FindStringIndex(car_num)) > 0 {
		sub := re.FindStringSubmatch(car_num)
		if sub[1] == "0" {
			if sub[2] == "0" {
				result = "digits/0&digits/0&digits/" + sub[3]
			} else if sub[2] == "1" {
				result = "digits/0&digits/" + sub[2] + sub[3]
			} else {
				if sub[3] == "0" {
					result = "digits/0&digits/" + sub[2] + "0"
				} else {
					result = "digits/0&digits/" + sub[2] + "0&digits/" + sub[3]
				}
			}
		} else {
			if sub[2] == "0" {
				if sub[3] == "0" {
					result = "digits/" + sub[1] + "00"
				} else {
					result = "digits/" + sub[1] + "00&digits/" + sub[3]
				}
			} else if sub[2] == "1" {
				result = "digits/" + sub[1] + "00&digits/" + sub[2] + sub[3]
			} else if sub[3] == "0" {
				result = "digits/" + sub[1] + "00&digits/" + sub[2] + sub[3]
			} else {
				result = "digits/" + sub[1] + "00&digits/" + sub[2] + "0&digits/" + sub[3]
			}
		}
	}

	return result
}
