package grpc

//
//import (
//	"eventbot/Logger"
//	"google.golang.org/grpc_proto"
//)
//
//type App struct {
//	gRPCServer *grpc_proto.Server
//	port       int
//}
//
//func New(port int) *App {
//	gRPCServer := grpc_proto.NewServer()
//
//	return &App{
//		gRPCServer: gRPCServer,
//		port:       port,
//	}
//}
//
//func (a *App) MustRun() {
//	if err := a.Run(); err != nil {
//		panic(err)
//	}
//}
//
//func (a *App) Run() error {
//	const op = "grpcapp.Run"
//
//	Logger.Sugar.Infoln("op", op)
//	Logger.Sugar.Infoln("port", a.port)
//	Logger.Sugar.Infoln("gRPC server is running")
//
//	//if err := a.gRPCServer.Serve(l); err != nil {
//	//	return Logger.Sugar.Errorf("%s: %w", op, err)
//	//}
//
//	return nil
//}
//
//func (a *App) Stop() {
//	const op = "graceful.Stop"
//
//	Logger.Sugar.Infof("stopping gRPC server, port: %v", a.port)
//
//	a.gRPCServer.GracefulStop()
//}
