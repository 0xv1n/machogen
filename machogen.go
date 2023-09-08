package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Commands struct {
	Commands []string `json:"commands"`
}

func main() {
	var commandString string
	var jsonFile string

	flag.StringVar(&commandString, "commands", "", "Comma-separated list of commands")
	flag.StringVar(&jsonFile, "json", "", "JSON file containing commands")
	flag.Parse()

	var commands Commands

	if jsonFile != "" {
		file, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			return
		}

		err = json.Unmarshal(file, &commands)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
	} else if commandString != "" {
		commands.Commands = strings.Split(commandString, ",")
	} else {
		fmt.Println("Either -commands or -json must be specified.")
		return
	}

	f, err := os.Create("generated_program.cpp")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	f.WriteString("#include <cstdlib>\n")
	f.WriteString("int main() {\n")

	for _, command := range commands.Commands {
		f.WriteString(fmt.Sprintf("    system(\"%s\");\n", command))
	}

	f.WriteString("    return 0;\n")
	f.WriteString("}\n")

	exec.Command("g++", "generated_binary.cpp", "-o", "generated_binary").Run()

	// Cleanup: Remove the generated C++ file
	//if err := os.Remove("generated_program.cpp"); err != nil {
	//	fmt.Println("Error deleting file:", err)
	//}
}

