package core

type Options struct {
    InitFunc  []func()
    DeferFunc []func()
}
