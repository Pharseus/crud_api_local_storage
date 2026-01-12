package main

import "crud-api-local-storage/route"

// @title File Management API
// @version 1.0
// @description API untuk upload, download, list, dan delete file menggunakan local storage
// @termsOfService http://swagger.io/terms/

func main() {
	route.StartRoute()
}
