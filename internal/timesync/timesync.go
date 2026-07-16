package timesync

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
)

const httpDateLayout = "Mon, 02 Jan 2006 15:04:05 GMT"
const dialAndReadTimeout = 5 * time.Second

func Sync(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(dialAndReadTimeout))

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	request := fmt.Sprintf("HEAD / HTTP/1.0\r\nHost: %s\r\n\r\n", host)
	if _, err := conn.Write([]byte(request)); err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			return fmt.Errorf("no Date header in response")
		}
		rest, ok := strings.CutPrefix(line, "Date: ")
		if !ok {
			continue
		}

		serverTime, err := time.Parse(httpDateLayout, rest)
		if err != nil {
			return fmt.Errorf("parsing Date header %q: %w", rest, err)
		}
		runtime.AdjustTimeOffset(-1 * int64(time.Since(serverTime)))
		return nil
	}
}
