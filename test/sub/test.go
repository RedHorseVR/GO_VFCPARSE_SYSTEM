package test

import (
	"github.com/MaxxSportsTV/maxx-video-backend-go/pkg/logic"
	"github.com/q191201771/naza/pkg/nazalog"
)
	go func() {
	 	config.StartgRPCServer(listner)
	}()

func Assert(expected interface{}, actual interface{}, extInfo ...string) {
	if !nazareflect.Equal(expected, actual) {
		var v string
		if len(extInfo) == 0 {
			v = fmt.Sprintf("assert failed. excepted=%+v, but actual=%+v", expected, actual)
		} else {
			v = fmt.Sprintf("assert failed. excepted=%+v, but actual=%+v, extInfo=%s", expected, actual, extInfo)
		}
		switch global.GetOption().AssertBehavior {
		case AssertError:
			global.Out(LevelError, 2, v)
		case AssertFatal:
			global.Out(LevelFatal, 2, v)
			fake.Os_Exit(1)
		case AssertPanic:
			global.Out(LevelPanic, 2, v)
			panic(v)
		}
	}
}

func Out(level Level, calldepth int, s string) {
	global.Out(level, calldepth, s)
}

func Sync() {
	global.Sync()
}

func WithPrefix(s string) Logger {
	return global.WithPrefix(s)
}

func GetOption() Option {
	return global.GetOption()
}

func main() {

	defer nazalog.Sync()

	sm := logic.NewMvbServer()
	err := sm.RunLoop()
	nazalog.Infof("server manager done. err=%+v", err)
	// router := gin.Default()

	listner, err := net.Listen("tcp", ":8485")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// config.LoadEnv()
	// run database
	// config.ConnectDB()

	// routes
	// routes.LiveStreamRoute(router)
	// router.Run(":8484")
	
	go func() {
	 	config.StartgRPCServer(listner)
	}()
	
	http.Serve(listner, router)

}
