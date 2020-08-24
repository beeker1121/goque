package goque

import (
	"fmt"
	"time"
)

func getTestPath() string {
	return fmt.Sprintf("/tmp/testdb/test_db_%d", time.Now().UnixNano())
}
