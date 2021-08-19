package object

// environment 是一个 hash map<string,Object>的结构
// 为了解决let identifier 的生命周期问题觉得使用层级的方式构建 env

type Environment interface {
	Get(name string) (Object, bool)
	Set(name string, val Object) Object
}

// globalEnv
type globalEnv struct {
	store map[string]Object
}

func (e *globalEnv) Get(name string) (val Object, ok bool) {
	val, ok = e.store[name]
	return
}

func (e *globalEnv) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

var (
	gEnv = &globalEnv{make(map[string]Object)}
)

func GlobalEnv() Environment {
	return gEnv
}

// localEnv
type localEnv struct {
	Environment // parent
	localStore  map[string]Object
}

func (e *localEnv) Get(name string) (Object, bool) {
	val, ok := e.localStore[name]
	if !ok && e.Environment != nil {
		return e.Environment.Get(name) // 没找到去上一层找
	}
	return val, ok

}

func (e *localEnv) Set(name string, val Object) Object {
	e.localStore[name] = val // 只在本地设置值
	return val
}

// 用于 带有{} 的语句
func WithLocalEnv(parent Environment) Environment {
	return &localEnv{parent, map[string]Object{}}
}
