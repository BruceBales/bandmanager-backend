package access

type WorkerChannels struct {
	CreateBandChan chan CreateBandJob
	EditBandChan   chan EditBandJob
	MemberChan     chan MemberJob
}

type Response struct {
	JobID    int64
	Err      error
	Message  string
	HTTPCode int
}
