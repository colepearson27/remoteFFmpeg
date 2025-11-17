package main

import (
	"encoding/json"
	"io"
	"strings"
	"os"
	"os/exec"
	"log"
	"fmt"
)

// This program is not secure!
// It is only meant to be used in local networks where configuration is known to not be malicious

// Maybe use a different way to ssh connect at some point

// Find a way to do like a file stream at some point to not wait for copying the file each way
type Config struct {
	Hostname string `json:"hostname"`
	Host string `json:"host"`
	SshPort string `json:"sshPort"`
}

func main() {
	// Read Config file
	var config *Config
	configFile, openErr := os.Open("config.json")
	if os.IsNotExist(openErr) {
	} else if openErr != nil {
		log.Fatalf("While trying to open config.json, got %s", openErr)
	}

	configData, readErr := io.ReadAll(configFile)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(configData, &config)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	log.Printf("Config: %v", config)

	// Get the FFmpegArgs to pass through
	var ffmpegArgs []string
	ffmpegArgs = os.Args
	var ffmpegOutputArgs []string

	var inputFile string
	var outputFile string

	// Get the files that are used in simple FFmpeg commands
	for i := 1; i < len(ffmpegArgs) - 1; i++ {
		if ffmpegArgs[i] == "-i" {
			inputFile = ffmpegArgs[i + 1]
			i++ // Skip the input file as well
			continue
		}
		ffmpegOutputArgs = append(ffmpegOutputArgs, ffmpegArgs[i])
	}
	outputFile = ffmpegArgs[len(ffmpegArgs) - 1]

	if strings.HasPrefix(outputFile, "file:") {
		outputFile = outputFile[5:]
		log.Fatal(outputFile)
	}

	log.Printf("FFmpeg args %v", ffmpegArgs)
	log.Printf("Input file: %s", inputFile)
	log.Printf("Output file: %s", outputFile)

	// Get a temp folder from host
	fmt.Println("Getting Temp Dir...")
	mkTempArgs := []string{fmt.Sprintf("%s@%s", config.Hostname, config.Host), "-p", config.SshPort, "mktemp", "-d"}
	mkTempCmd := exec.Command("ssh", mkTempArgs...)
	cmdOutput, mkTempRunErr := mkTempCmd.Output()
	if mkTempRunErr != nil {
		log.Fatal(mkTempRunErr)
	}

	var tmpFolder string
	// Remove the newline character at the end
	tmpFolder = string(cmdOutput[:len(cmdOutput) - 1])
	log.Printf("Temp Folder: %s", tmpFolder)

	// Copy input file to temp directory to host
	fmt.Println("Copying File to host...")
	scpInputArgs := []string{inputFile, fmt.Sprintf("%s@%s:%s/", config.Hostname, config.Host, tmpFolder)}
	scpInputCmd := exec.Command("scp", scpInputArgs...)
	scpOutput, scpInputErr := scpInputCmd.CombinedOutput()
	if scpInputErr != nil {
		log.Println(string(scpOutput))
		log.Fatal(scpInputErr)
	}

	// Run the FFmpeg command
	// This assumes that the FFmpeg command can run on the target computer
	fmt.Println("Running FFmpeg Command...")
	ffmpegCommandArgs := []string{fmt.Sprintf("%s@%s", config.Hostname, config.Host), "-p", config.SshPort, "ffmpeg", "-i", fmt.Sprintf("%s/%s", tmpFolder, inputFile)}
	ffmpegCommandArgs = append(ffmpegCommandArgs, ffmpegOutputArgs...)
	ffmpegCommandArgs = append(ffmpegCommandArgs, fmt.Sprintf("%s/%s", tmpFolder, outputFile))
	fmt.Println(ffmpegCommandArgs)
	ffmpegCmd := exec.Command("ssh", ffmpegCommandArgs...)
	ffmpegOutput, ffmpegErr := ffmpegCmd.CombinedOutput()
	if ffmpegErr != nil {
		log.Println(string(ffmpegOutput))
		log.Fatal(ffmpegErr)
	}
	log.Println(string(ffmpegOutput))

	// Copy the resulting file back to our computer
	fmt.Println("Copying File from Host...")
	scpOutputArgs := []string{fmt.Sprintf("%s@%s:%s/%s", config.Hostname, config.Host, tmpFolder, outputFile), outputFile}
	scpOutputCmd := exec.Command("scp", scpOutputArgs...)
	scpOutput, scpOutputErr := scpOutputCmd.CombinedOutput()
	if scpOutputErr != nil {
		log.Println(string(scpOutput))
		log.Fatal(scpOutputErr)
	}

	// Remove the temp folder when done with it
	fmt.Println("Removing Temp Dir...")
	deleteTempArgs := []string{fmt.Sprintf("%s@%s", config.Hostname, config.Host), "-p", config.SshPort, "rm", "-r", tmpFolder}
	deleteTempCmd := exec.Command("ssh", deleteTempArgs...)
	cmdOutput, runErr := deleteTempCmd.Output()
	if runErr != nil {
		log.Fatal(runErr)
	}
}
