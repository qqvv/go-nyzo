package message

type MsgType uint16

const (
	Invalid MsgType = iota + 1
	_               // BootstrapRequest, deprecated
	_               // BootstrapResponse, deprecated
	NodeJoin
	NodeJoinResponse
	Transaction
	TransactionResponse
	PrevHashRequest
	PrevHashResponse
	NewBlock
	NewBlockResponse
	BlockRequest
	BlockResponse
	TxPoolRequest
	TxPoolResponse
	MeshRequest
	MeshResponse
	BlockVote
	BlockVoteResponse
	NewVerivierVote
	NewVerivierVoteResponse
	MissingBlockVoteRequest
	MissingBlockVoteResponse
	MissingBlockRequest
	MissingBlockResponse
	TimestampRequest
	TimestampResponse
	HashVoteOverrideRequest
	HashVoteOverrideResponse
	ConsensusThresholdOverrideRequest
	ConsensusThresholdOverrideResponse
	NewVerivierVoteOverrideRequest
	NewVerivierVoteOverrideResponse
	BootstrapRequestV2
	BootstrapResponseV2
	BlockWithVotesRequest
	BlockWithVotesResponse
)

const (
	Ping MsgType = iota + 200
	PingResponse
)

const (
	UpdateRequest MsgType = iota + 300
	UpdateResponse
)

const (
	BlockRejectionRequest MsgType = iota + 400
	BlockRejectionResponse
	DetachmentRequest
	DetachmentResponse
	UnfrozenBlockPoolPurgeRequest
	UnfrozenBlockPoolPurgeResponse
	UnfrozenBlockPoolStatusRequest
	UnfrozenBlockPoolStatusResponse
	MeshStatusRequest
	MeshStatusResponse
	TogglePauseRequest
	TogglePauseResponse
	ConsensusTallyStatusRequest
	ConsensusTallyStatusResponse
	NewVerifierTallyStatusRequest
	NewVerifierTallyStatusResponse
	BlacklistStatusRequest
	BlacklistStatusResponse
)

const (
	ResetRequest MsgType = iota + 500
	ResetResponse
)

const (
	IncomingRequest MsgType = iota + 65533
	Error
	Unknown
)
