package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/trace"
	"strings"
	"time"
)

func quit() {
	fmt.Println("Quitting program...")
	os.Exit(3)
}

// CREATE
func createFile(scanner *bufio.Scanner) {
	fmt.Printf("\nCOMMAND: 'create'\n\n")
	infoLogger.Println("User selected CREATE command...")

	fmt.Println("What would you like to name the file? (No need to add file extension)")
	for scanner.Scan() {
		// check if user wants to quit instead
		input := strings.TrimSuffix(scanner.Text(), ".txt")
		if input == "/q" {
			infoLogger.Println("User quitting...")
			quit()
		}
		// Check if empty filename
		if input == "" {
			fmt.Printf("Please enter a valid filename.")
			fmt.Printf("\n\nReturning to beginning...\n\n")
			time.Sleep(3 * time.Second)
			displayCommands()
			return
		}

		// Attempt to create file provided from the user
		_, err := os.Create(filepath.Join("./text-files", filepath.Base(input+".txt")))
		if err != nil {
			ErrorReturnToMenu(err)
			return
		}
		infoLogger.Printf("User created file: %v\n", input)
		fmt.Printf("%v.txt has been created.\n", input)
		fmt.Printf("\n\nReturning to beginning...\n\n")
		time.Sleep(3 * time.Second)
		displayCommands()
		return
	}
}

// READ
func readFile(scanner *bufio.Scanner) {
	var f *os.File
	fmt.Printf("\nCOMMAND: 'READ'\n\n")
	infoLogger.Println("User selected OPEN command...")

	// show list of files to choose from
	showFiles(scanner)
	fmt.Println("Which one of the files above would you like to open? (No need to add file extension)")
	for scanner.Scan() {
		input := scanner.Text()
		pathWithFile := fmt.Sprintf("./text-files/%v.txt", input)

		// Open the file that the user specified
		file, err := os.Open(pathWithFile)
		if err != nil {
			ErrorReturnToMenu(err)
			return
		}
		f = file
		fmt.Println()
		break
	}

	// Create new scanner to read the file
	rd := bufio.NewReader(f)
	fmt.Println("---BOF---")
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			ErrorReturnToMenu(err)
			return
		}
		fmt.Print(line)
	}
	fmt.Println("---EOF---")
	displayCommands()
}

// UPDATE
func writeToFile(scanner *bufio.Scanner) {
	fmt.Printf("\nCOMMAND: 'WRITE'\n\n")
	infoLogger.Println("User selected WRITE command...")

	var fileToWrite *os.File
	showFiles(scanner)
	fmt.Println("Which file above would you like to write to? (No need to add file extension)")
	for scanner.Scan() {
		input := scanner.Text()
		path := fmt.Sprintf("./text-files/%v.txt", input)
		if input == "/q" {
			infoLogger.Println("User quitting...")
			quit()
		}
		// Open file
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		fileToWrite = file
		if err != nil {
			ErrorReturnToMenu(err)
			return
		}
		break
	}
	fmt.Println("What would you like to write into the file?")

	// Get text to actuqlly write into the file
	textToWrite := bufio.NewScanner(os.Stdin)
	for textToWrite.Scan() {
		input := textToWrite.Text()
		_, err := fileToWrite.Write([]byte(input))
		if err != nil {
			ErrorReturnToMenu(err)
			return
		}
		break
	}
	defer fileToWrite.Close()

}

// DELETE
func removeFile(scanner *bufio.Scanner) {
	fmt.Printf("\nCOMMAND: 'REMOVE'\n\n")
	infoLogger.Println("User selected REMOVE command...")
	// show current list of files
	showFiles(scanner)
	fmt.Println("What file would you like to remove from the list above? (No need to add file extension)")
	for scanner.Scan() {
		input := scanner.Text()
		s := fmt.Sprintf("./text-files/%v.txt", input)
		err := os.Remove(s)
		if err != nil {
			ErrorReturnToMenu(err)
			return
		}
		infoLogger.Printf("%v.txt has been successfully removed.\n", input)
		fmt.Printf("%v.txt has been successfully removed.\n", input)
		fmt.Printf("\n\nReturning to beginning...\n\n")
		time.Sleep(3 * time.Second)
		displayCommands()
		return
	}

}

// READ -- List files
func showFiles(scanner *bufio.Scanner) {
	fmt.Printf("\nCOMMAND: 'SHOW'\n\n")
	infoLogger.Println("User selected SHOW command...")
	files, err := ioutil.ReadDir("./text-files")
	if err != nil {
		ErrorReturnToMenu(err)
		return
	}
	if len(files) == 0 {
		fmt.Println("No files yet, would you like to create a new file? (y/n)")
		for scanner.Scan() {
			input := scanner.Text()
			strings.ToLower(input)
			switch input {
			case "y":
				createFile(scanner)
				return
			case "yes":
				createFile(scanner)
				return
			case "n":
				return
			case "no":
				return
			}
		}
	}
	for _, f := range files {
		// format files in table like format
		fmt.Printf("|%-7s|\n", f.Name())
	}
	return
}

// Error occurred -- return to menu
func ErrorReturnToMenu(err error) {
	fmt.Println("An error occurred, check logs...")
	errorLogger.Println(err)
	fmt.Printf("\n\nReturning to beginning...\n\n")
	time.Sleep(1 * time.Second)
	fmt.Println(".")
	time.Sleep(1 * time.Second)
	fmt.Println(".")
	time.Sleep(1 * time.Second)
	fmt.Println(".")
	fmt.Println(" ")
	displayCommands()
	return
}
func displayCommands() {
	// TODO: turn into an array
	fmt.Println("Commands: \n1. '/h': show commands \n2. '/q': quit \n3. 'show': show list of files in dir \n4. 'write': show file contents \n5. 'read': Read File \n6. 'create': Create new file \n7. 'remove': Remove file")
}

// Log Levels
var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
)

// This method does not need to be explicitly called, it is called when main() is ran
func init() {

	// Create file if it does not exist, this is the log file
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fatalLogger.Fatal("Could not create required log file.")
	}

	// intialize the different types of logs -- formats
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(file, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Tracing -- To "READ" the trace.out file run the command: 'go tool trace trace.out'
	f, err := os.Create("trace.out")
	if err != nil {
		fatalLogger.Fatalf("File creation error: %v\n", err)
	}

	// close trace on program exit
	defer func() {
		if err := f.Close(); err != nil {
			fatalLogger.Fatalf("Could not close trace.out file: %v\n", err)
		}
	}()

	// Start trace
	if err := trace.Start(f); err != nil {
		fatalLogger.Fatalf("Could not start trace: %v\n", err)
	}

	// stop trace
	defer trace.Stop()
}

func main() {

	// Create text-files directory if none found
	if _, err := os.Stat("./text-files"); os.IsNotExist(err) {
		// susceptible to race condition but this app is just an example
		err := os.Mkdir("./text-files", 0777)
		if err != nil {
			fatalLogger.Fatal(err)
		}
	}

	// Show commands
	displayCommands()
	// Create scanner to read input
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "/h":
			displayCommands()
		case "/q":
			quit()
		case "show":
			showFiles(scanner)
		case "create":
			createFile(scanner)
		case "read":
			readFile(scanner)
		case "remove":
			removeFile(scanner)
		case "write":
			writeToFile(scanner)
		}
	}

}
