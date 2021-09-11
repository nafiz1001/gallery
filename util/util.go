package util

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func CreateId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", sha1.Sum([]byte(strconv.FormatInt(int64(rand.Int()), 36))))[:7]
}
