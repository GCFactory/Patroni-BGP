package patroni_bgp

const (
	PatroniStateUndefined = iota
	PatroniStateError
	PatroniStateMaster
	PatroniStateReplica
	PatroniStateSyncReplica
	PatroniStateAsyncReplica
)
