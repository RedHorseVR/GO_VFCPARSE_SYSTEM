package logic

type IMvbServer interface {
	RunLoop() error
	Dispose()

	// // AddCustomizePubSession 定制化增强功能。业务方可以将自己的流输入到 ILalServer 中
	// //
	// // @example 示例见 lal/app/demo/customize_lalserver
	// //
	// // @doc 文档见 <lalserver二次开发 - pub接入自定义流> https://pengrl.com/lal/#/customize_pub
	// //
	// AddCustomizePubSession(streamName string) (ICustomizePubSessionContext, error)

	// // DelCustomizePubSession 将 ICustomizePubSessionContext 从 ILalServer 中删除
	// //
	// DelCustomizePubSession(ICustomizePubSessionContext)

	// // StatLalInfo StatAllGroup StatGroup CtrlStartPull CtrlKickOutSession
	// //
	// // 一些获取状态、发送控制命令的API。
	// // 目的是方便业务方在不修改logic包内代码的前提下，在外层实现一些特定逻辑的定制化开发。
	// //
	// StatLalInfo() base.LalInfo
	// StatAllGroup() (sgs []base.StatGroup)
	// StatGroup(streamName string) *base.StatGroup
	// CtrlStartRelayPull(info base.ApiCtrlStartRelayPullReq) base.ApiCtrlStartRelayPull
	// CtrlStopRelayPull(streamName string) base.ApiCtrlStopRelayPull
	// CtrlKickSession(info base.ApiCtrlKickSession) base.HttpResponseBasic
}

func NewMvbServer() IMvbServer {
	return NewServerManager()
}
