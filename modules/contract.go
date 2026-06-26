package modules

// ModuleContract 定义了每个模块向系统暴露的生命周期契约。
type ModuleContract interface {
	Start() error
	Stop() error
}
