package main

func main() {
	server := &Server{}

	if err := server.Init(); err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}
}
