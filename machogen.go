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
	var commandString, jsonFile, simulation string

	flag.StringVar(&commandString, "commands", "", "Comma-separated list of commands")
	flag.StringVar(&jsonFile, "json", "", "JSON file containing commands")
	flag.StringVar(&simulation, "s", "", "Simulation mode. Format: <type>:<param>, e.g., N:1.1.1.1 will open a TCP socket to 1.1.1.1")

	flag.Parse()

	// EXAMPLE Behavior
	// ./machogen -commands "echo Hello World" -s N:1.1.1.1
	// Will run (in sequence):
	//	1. sh -c echo Hello World
	//	2. open TCP Socket to 1.1.1.1:80 sourcing directly from the generated bin
	//		rather than a separate sh -c command. This is to differentiate behavior
	//		directly from the binary from shell commands launched as childprocs.

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
	f.WriteString("#include <iostream>\n")
	f.WriteString("#include <sys/socket.h>\n")
	f.WriteString("#include <netinet/in.h>\n")
	f.WriteString("#include <arpa/inet.h>\n")
	f.WriteString("#include <unistd.h>\n")
	f.WriteString("int main() {\n")

	for _, command := range commands.Commands {
		f.WriteString(fmt.Sprintf("    system(\"%s\");\n", command))
	}

	if simulation != "" {
		simType := strings.Split(simulation, ":")[0]
		simParam := strings.Split(simulation, ":")[1]
		// Network Connection Simulation
		// atm im using TCP sockets rather than just a PING since ICMP may not always be avail
		// port can be modified, with a default of 80 - may add cli parsing to read it in from user
		if simType == "N" {
			f.WriteString("    int sockfd = socket(AF_INET, SOCK_STREAM, 0);\n")
			f.WriteString("    struct sockaddr_in servaddr;\n")
			f.WriteString(fmt.Sprintf("    servaddr.sin_family = AF_INET;\n"))
			f.WriteString(fmt.Sprintf("    servaddr.sin_port = htons(80);\n"))
			f.WriteString(fmt.Sprintf("    inet_pton(AF_INET, \"%s\", &servaddr.sin_addr);\n", simParam))
			f.WriteString("    connect(sockfd, (struct sockaddr *)&servaddr, sizeof(servaddr));\n")
			f.WriteString("    close(sockfd);\n")
		}
	}

	f.WriteString("    return 0;\n")
	f.WriteString("}\n")

	exec.Command("clang", "generated_binary.cpp", "-o", "generated_binary").Run()

	// Cleanup: Remove the generated C++ file - can comment below out if you're curious what the generated code looks like
	if err := os.Remove("generated_binary.cpp"); err != nil {
		fmt.Println("Error deleting file:", err)
	}
}
