package main

import "golang-cohort-5/tasks/sesi-5/routers"

func main() {
	var PORT = ":8080"

	routers.StartServer().Run(PORT)
}
