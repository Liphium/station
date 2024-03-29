package pipeshutil

import "log"

var Log = log.New(log.Writer(), "pipes-handler ", log.Flags())

func RemoveString(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}
