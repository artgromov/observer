package configs

import (
	"fmt"
	"strconv"
	"strings"
)

func validateAddrString(value string) error {
	items := strings.Split(value, ":")
	if len(items) != 2 {
		return fmt.Errorf("invalid addr: %s, must be in format \"address:port\"", value)
	}
	port, err := strconv.Atoi(items[1])
	if err != nil {
		return err
	}
	if port < 0 || 65535 < port {
		return fmt.Errorf("invalid addr: %s, port must be in (0,65535) range", value)
	}
	return nil
}
