package main

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	notionRedownload()
}
