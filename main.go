/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "everlasting/cmd"

// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	cmd.Execute()
}
