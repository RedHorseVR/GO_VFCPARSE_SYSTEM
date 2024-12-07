package logic

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerManager struct {
	// option          Option
	serverStartTime int64
	config          *Config
	dbClient        *DbConnClient

	// servicesManager
	// httpServerManager *base.HttpServerManager
	// httpServerHandler *HttpServerHandler
	// hlsServerHandler  *hls.ServerHandler

	// rtmpServer     *rtmp.Server
	// rtspServer     *rtsp.Server
	httpApiServer *HttpApiServer
	grpcApiServer *GrpcApiServer
	// grpcClientConn *GrpcClientConn
	// pprofServer    *http.Server
	exitChan chan struct{}

	// mutex sync.Mutex
	// groupManager IGroupManager

	// simpleAuthCtx *SimpleAuthCtx
}

func NewServerManager() *ServerManager {
	sm := &ServerManager{
		serverStartTime: time.Now().UnixNano(),
		exitChan:        make(chan struct{}, 1),}

	sm.config = LoadConfAndInitLog()

	// Enable Gin HTTP Server
	sm.httpApiServer = NewHttpApiServer(sm.config.HttpApiConfig.Addr, sm)

	// Enable gRPC Server
	sm.grpcApiServer = NewGrpcApiServer(sm.config.GrpcApiConfig.Addr, sm)

	// Connect to database
	sm.dbClient = NewDbConn(&sm.config.DbClientConfig, sm)

	// Connect to database
	// sm.dbClient = NewDbConn()

	return sm
}

func (sm *ServerManager) RunLoop() error {

	go sm.RunKillSignalHandler()

	if sm.httpApiServer != nil {
		if err := sm.httpApiServer.Listen(); err != nil {
			return err
		}
		go func() {
			if err := sm.httpApiServer.RunLoop(); err != nil {
				Log.Error(err)
			}
		}()	}

	if sm.grpcApiServer != nil {
		if err := sm.grpcApiServer.Listen(); err != nil {
			return err
		}
		go func() {
			if err := sm.grpcApiServer.RunLoop(); err != nil {
				Log.Error(err)
			}
		}()	}

	if sm.dbClient != nil {
		go func() {
			if err := sm.dbClient.RunLoop(); err != nil {
				Log.Error(err)
			}
		}()	}

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	var tickCount uint32
	for {
		select {
		case <-sm.exitChan:
			return nil
		case <-t.C:
			tickCount++

			// sm.mutex.Lock()

			// // 关闭空闲的group
			// sm.groupManager.Iterate(func(group *Group) bool {
			// 	if group.IsInactive() {
			// 		Log.Infof("erase inactive group. [%s]", group.UniqueKey)
			// 		group.Dispose()
			// 		return false
			// 	}

			// 	group.Tick(tickCount)
			// 	return true
			// })

			// // 定时打印一些group相关的debug日志
			// if sm.config.DebugConfig.LogGroupIntervalSec > 0 &&
			// 	tickCount%uint32(sm.config.DebugConfig.LogGroupIntervalSec) == 0 {
			// 	groupNum := sm.groupManager.Len()
			// 	Log.Debugf("DEBUG_GROUP_LOG: group size=%d", groupNum)
			// 	if sm.config.DebugConfig.LogGroupMaxGroupNum > 0 {
			// 		var loggedGroupCount int
			// 		sm.groupManager.Iterate(func(group *Group) bool {
			// 			loggedGroupCount++
			// 			if loggedGroupCount <= sm.config.DebugConfig.LogGroupMaxGroupNum {
			// 				Log.Debugf("DEBUG_GROUP_LOG: %d %s", loggedGroupCount, group.StringifyDebugStats(sm.config.DebugConfig.LogGroupMaxSubNumPerGroup))
			// 			}
			// 			return true
			// 		})
			// 	}
			// }

			// sm.mutex.Unlock()

			// // 定时通过http notify发送group相关的信息
			// if uis != 0 && (tickCount%uis) == 0 {
			// 	updateInfo.ServerId = sm.config.ServerId
			// 	updateInfo.Groups = sm.StatAllGroup()
			// 	sm.option.NotifyHandler.OnUpdate(updateInfo)
			// }
		}
	}

	// never reach here
}

func (sm *ServerManager) Dispose() {
	Log.Debug("dispose server manager.")

	// Dospose Db Connection at the end
	if sm.httpApiServer != nil {
		sm.httpApiServer.Dispose()
	}

	// sm.mutex.Lock()
	// sm.groupManager.Iterate(func(group *Group) bool {
	// 	group.Dispose()
	// 	return true
	// })
	// sm.mutex.Unlock()

	sm.exitChan <- struct{}{}
}

func (sm *ServerManager) RunKillSignalHandler() {

	Log.Info("Starting RunKillSignalHandler")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c

	Log.Infof("kill signal recv. s=%+v", s)

	sm.Dispose()
}
