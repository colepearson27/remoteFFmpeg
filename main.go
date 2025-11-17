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

	// cat test.mp4 | ssh robot@robot ffmpeg -i - -map_metadata -1 -c:v libsvtav1 -crf 30 -preset 1 -b:v 0 -g 60 -movflags +faststart -c:a copy -f matroska - > c.mkv

	// Steam to remote, transcode, and stream back to local
	everythingArgs := []string{fmt.Sprintf("%s@%s", config.Hostname, config.Host), "ffmpeg", "-i", "-"}
	everythingArgs = append(everythingArgs, ffmpegOutputArgs...)
	everythingArgs = append(everythingArgs, []string{"-f", "matroska", "-"}...)
	fmt.Println("Doing Everything...")
	fmt.Println(everythingArgs)
	everythingCmd := exec.Command("ssh", everythingArgs...)

	c1 := exec.Command("cat", inputFile)
    everythingCmd.Stdin, _ = c1.StdoutPipe()
	c2 := exec.Command("tee", outputFile)
	c2.Stdin, _ = everythingCmd.StdoutPipe()
    _ = everythingCmd.Start()
	_ = c2.Start()
    _ = c1.Start()
    _ = everythingCmd.Wait()
	_ = c1.Wait()
	_ = c2.Wait()

	// Remove the temp folder when done with it
	fmt.Println("Removing Temp Dir...")
	deleteTempArgs := []string{fmt.Sprintf("%s@%s", config.Hostname, config.Host), "-p", config.SshPort, "rm", "-r", tmpFolder}
	deleteTempCmd := exec.Command("ssh", deleteTempArgs...)
	cmdOutput, runErr := deleteTempCmd.Output()
	if runErr != nil {
		log.Fatal(runErr)
	}
}
