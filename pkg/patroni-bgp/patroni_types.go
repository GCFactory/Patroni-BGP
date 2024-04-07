package patroniBgp

const (
	PatroniStateUndefined = iota
	PatroniStateError
	PatroniStateMaster
	PatroniStateReplica
	PatroniStateSyncReplica
	PatroniStateAsyncReplica
)
