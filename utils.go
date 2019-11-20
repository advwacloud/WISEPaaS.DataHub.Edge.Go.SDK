package agent

import (
	"encoding/json"
	"time"
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func setInterval(someFunc func(), seconds int, async bool) chan bool {
	interval := time.Duration(seconds) * time.Second
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				close(clear)
				return
			}
		}
	}()
	return clear
}
