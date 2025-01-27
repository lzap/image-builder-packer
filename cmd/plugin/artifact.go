package main

type artifact struct {
}

func (a artifact) BuilderId() string {
	return "osbuild.image-builder"
}

func (a artifact) Files() []string {
	return nil
}

func (a artifact) Id() string {
	return ""
}

func (a artifact) String() string {
	return ""
}

func (a artifact) State(name string) interface{} {
	return nil
}

func (a artifact) Destroy() error {
	return nil
}
