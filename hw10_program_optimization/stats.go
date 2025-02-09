package hw10programoptimization

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/bytedance/sonic" //nolint:all
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	var user User

	i := 0
	dec := sonic.ConfigFastest.NewDecoder(r)

	for {
		if err = dec.Decode(&user); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return
		}
		result[i] = user
		i++
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	reg, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := reg.MatchString(user.Email)
		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
