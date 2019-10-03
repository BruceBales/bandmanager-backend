package access

type WorkerChannels struct {
	CreateBandChan chan CreateBandJob
	EditBandChan   chan EditBandJob
	MemberChan     chan MemberJob
}
