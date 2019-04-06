package helpers

func WatchFor(signal chan bool, method func() bool) {
	for {
		select {
		case <-signal:
			return
		default:
			if method() {
				return
			}
		}
	}
}
