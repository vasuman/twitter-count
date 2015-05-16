package res

//go:generate go-res-pack ./data/ res_gen.go

func Setup() {
	genInit()
}
