package types

// BITCOINBridge defines the interface of custom bitcoin bridge.
type BITCOINBridge interface {
	// Deposit transfers amout to address
	Deposit(string, string, int64) (string, error)
}
