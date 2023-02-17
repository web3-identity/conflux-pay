package utils

import "time"

func Retry(count int, interval time.Duration, fn func() error) error {
	for i := 0; i < count; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		if i == count-1 {
			return err
		}
		time.Sleep(interval)
	}
	return nil
}
