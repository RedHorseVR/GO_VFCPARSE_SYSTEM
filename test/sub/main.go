package main

import (
	"github.com/MaxxSportsTV/maxx-video-backend-go/pkg/logic"
	"github.com/q191201771/naza/pkg/nazalog"
)

func main() {

	defer nazalog.Sync()

	sm := logic.NewMvbServer()
	err := sm.RunLoop()
	nazalog.Infof("server manager done. err=%+v", err)
	// router := gin.Default()

	// listner, err := net.Listen("tcp", ":8485")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// config.LoadEnv()
	// //run database
	// config.ConnectDB()

	// //routes
	// routes.LiveStreamRoute(router)
	// router.Run(":8484")
	// go func() {
	// 	config.StartgRPCServer(listner)
	// }()
	// http.Serve(listner, router)

}
