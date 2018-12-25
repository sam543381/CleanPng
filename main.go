package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func displayHelp(section string) {

	switch section {
	case "commands":
		fmt.Println("Avaible commands: info, clean, help")

	case "info":
		fmt.Println("info <path>")
		fmt.Println("\tDisplays PNG image metadata related to the given image.")

	case "clean":
		fmt.Println("clean <path>")
		fmt.Println("\tCreates a copy of the given PNG image without metdata.")
		fmt.Println("WARNING: It may conserve steganographic watermarks included in mandotary data")

	default:
		fmt.Printf("No help avaible for section %s\n", section)
	}
}

func main() {

	if len(os.Args) < 2 {
		displayHelp("commands")
		os.Exit(1)
	}

	label := os.Args[1]
	args := os.Args[2:]

	processCommand(label, args)

}

func processCommand(label string, args []string) {

	switch label {

	case "help":

		if len(args) < 1 {
			displayHelp("commands")
			return
		}

		section := args[0]
		displayHelp(section)

	case "info":

		if len(args) < 1 {
			displayHelp("info")
			return
		}

		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}

		png, err := DecodePng(file)
		if err != nil {
			fmt.Println("Unable to decode PNG image")
			log.Fatal(err)
		}

		fmt.Println(png.String())

	case "clean":
		if len(args) < 1 {
			displayHelp("clean")
			return
		}

		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}

		png, err := DecodePng(file)
		file.Close()
		if err != nil {
			fmt.Println("Unable to decode PNG image")
			log.Fatal(err)
		}

		bytes, err := png.EncodeAndClean()
		if err != nil {
			log.Fatal(err)
		}

		ioutil.WriteFile("out.png", bytes, 0644)

	default:
		displayHelp("commands")
		log.Fatalf("Unknown command with name %s\n", label)
	}

}
