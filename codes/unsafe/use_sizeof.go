

// ClientConn in other package 
//type ClientConn struct {
//	ctx    context.Context
//	cancel context.CancelFunc
//
//	target       string
//	parsedTarget resolver.Target
//	authority    string
//	dopts        dialOptions
//	csMgr        *connectivityStateManager
//
//	balancerBuildOpts balancer.BuildOptions
//	blockingpicker    *pickerWrapper
//
//	mu              sync.RWMutex
//	resolverWrapper *ccResolverWrapper
//	sc              *ServiceConfig
//	conns           map[*addrConn]struct{}
//	// Keepalive parameter can be updated if a GoAway is received.
//	mkp             keepalive.ClientParameters
//	curBalancerName string
//	balancerWrapper *ccBalancerWrapper
//	retryThrottler  atomic.Value
//
//	firstResolveEvent *grpcsync.Event
//
//	channelzID int64 // channelz unique identification number
//	czData     *channelzData
//}


var (
	Target *string
)

func (cm *ClientConnManager) Call(method string, in interface{}, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	cc, err := cm.getClientConn()
	if err != nil {
		return err
	}
  // use unsafe.Pointer and uintptr get unexported field parsedTarget
	Target = (* string)(unsafe.Pointer(uintptr(unsafe.Pointer(cc)) + unsafe.Sizeof(context.TODO()) +
		unsafe.Sizeof(* new(context.CancelFunc)) + unsafe.Sizeof("1")))

	return cc.Invoke(ctx, method, in, out)
}
