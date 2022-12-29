package csv

import (
	"fmt"
	"mud/utils/ui/logger"
	"os"
	"sort"
	"strings"
	"sync"
)

// Represents a csv file
type CSVFile struct {
	Columns       []string    // The column headers of the csv file
	LineLocations []int64     // Indicates the byte offsets of each line in the file
	LineCount     int64       // The number of lines in the file
	Filepath      string      // The path of the file
	lock          *sync.Mutex // For concurrency
}

func checkError(e interface{}) {
	if e != nil {
		logger.ErrorCustomCaller(1, e)
		panic(e)
	}
}

// Fetches the given number of bytes from the given file,
// Returns the bytes and the number of actual bytes read.
func fetchBytes(f *os.File, n int) ([]byte, int) {
	var buff []byte = make([]byte, n)
	nOut, err := f.Read(buff)
	if nOut > 0 {
		checkError(err) // Because EOF is an error
		return buff[:nOut], nOut
	}
	return []byte{}, 0
}

// Parses a buffer of bytes and returns an array of string column data.
func parseCSVLine(line []byte) []string {
	var inQuote bool = false
	var escaped bool = false
	var lastEntryEnd int = 0
	var entries []string

	for i, b := range line {
		if escaped {
			escaped = false
			continue
		}

		r := rune(b)

		if r == '\\' {
			escaped = true
		} else if r == '"' {
			inQuote = !inQuote
		} else if r == ',' && !inQuote {
			entries = append(entries, string(line[lastEntryEnd:i]))
			lastEntryEnd = i + 1
		}
	}

	entries = append(entries, string(line[lastEntryEnd:]))

	return entries
}

// Reads a single line from the csv file into the given buffer
// Returns the string column data, the left over bytes in the buffer,
// and whether the end of the file was reached.
func readLine(f *os.File, buffer []byte) ([]string, []byte, bool) {
	rawData, buffOut, eof := readLineRaw(f, buffer)
	return parseCSVLine(rawData), buffOut, eof
}

// Reads raw data from the given line,
// Returns the raw string as bytes from the buffer,
// the left over bytes that were read in, but not used,
// and whether the end of the file was reached or not.
func readLineRaw(f *os.File, buffer []byte) ([]byte, []byte, bool) {
	var eof bool = false

	data := buffer
	for true {
		var foundLine bool = false
		var lineIndexEnd int = -1

		for i, b := range data {
			if rune(b) == '\n' || (eof && i == len(data)-1) {
				foundLine = true
				if eof && i == len(data)-1 {
					lineIndexEnd = i + 1
				} else {
					lineIndexEnd = i
				}
				break
			}
		}

		if foundLine {
			line := data[:lineIndexEnd]
			if lineIndexEnd < len(data)-1 {
				data = data[lineIndexEnd+1:]
			} else {
				data = []byte{}
			}
			return line, data, eof
		}

		newData, nOut := fetchBytes(f, 1024)
		eof = nOut < 1024
		data = append(data, newData...)
	}

	return []byte{}, buffer, true
}

// Creates a new csv file with the given columns and the initial data.
func CreateCSV(path string, columns []string, lines [][]string) CSVFile {
	logger.Info("Creating %s CSV File...", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
	checkError(err)

	firstLinePos, err := f.WriteString(fmt.Sprintf("%d\n%s\n", 0, strings.Join(columns, ",")))
	checkError(err)

	f.Close()

	file := CSVFile{
		columns,
		[]int64{int64(firstLinePos)},
		0,
		path,
		&sync.Mutex{},
	}

	for _, line := range lines {
		file.AppendLine(line)
	}

	return file
}

// Parses an existing csv file and returns the data in the csv struct
func ParseCSV(path string) CSVFile {
	logger.Info("Reading %s CSV File...", path)

	f, err := os.Open(path)
	checkError(err)

	var buffer []byte
	var eof bool
	var headers []string
	var currentCursor int = 0
	var i int64

	// Read Line Count
	var lineCount int64
	var rawLine []byte
	rawLine, buffer, eof = readLineRaw(f, buffer)
	_, err = fmt.Sscanf(string(rawLine), "%d", &lineCount)
	checkError(err)

	// Read Column Headers
	headers, buffer, eof = readLine(f, buffer)

	// Adjust cursor
	for _, h := range headers {
		currentCursor += len(h)
	}
	currentCursor += len(headers) + len(rawLine) + 1

	// read remaining lines
	var lineLocations []int64
	var line []string
	for i = 0; len(buffer) > 0 && !eof && i < lineCount; i++ {
		lineLocations = append(lineLocations, int64(currentCursor))

		line, buffer, eof = readLine(f, buffer)

		// Adjust the cursor
		for _, l := range line {
			currentCursor += len(l)
		}
		currentCursor += len(line)
	}

	// So that we can keep track of how long the last line is
	lineLocations = append(lineLocations, int64(currentCursor))

	f.Close()

	return CSVFile{
		headers,
		lineLocations,
		lineCount,
		path,
		&sync.Mutex{},
	}
}

// Fetches a specific line from the csv file
func (cf *CSVFile) ReadSpecificLine(line int64) []string {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	f, err := os.Open(cf.Filepath)
	checkError(err)

	var result []string
	lineStart := cf.LineLocations[line]
	_, err = f.Seek(lineStart, 0)
	checkError(err)

	lineStop := cf.LineLocations[line+1]
	lengthOfLine := lineStop - lineStart - 1

	buff := make([]byte, lengthOfLine)
	nOut, err := f.Read(buff)
	checkError(err)
	if int64(nOut) != lengthOfLine {
		logger.Warn("Read different byte amount than expected for line %d of %s, read %d, expected %d",
			line, cf.Filepath, nOut, lengthOfLine)
	}

	result = parseCSVLine(buff)

	f.Close()

	return result
}

// Shifts the contents of the file by a certain number of bytes
// starting at the given starting byte offset
// the shift distance can be both positive and negative.
func shiftFileContents(f *os.File, offset int64, start int64) {
	var nextBuffer []byte
	var newNOut int
	var nOut int
	var seekTarget int64
	var buff []byte
	var nOff int64
	var totalSize int64 = start + offset

	nOff, err := f.Seek(start, 0)
	checkError(err)

	buff, nOut = fetchBytes(f, 1024)
	for nOut > 0 {
		totalSize += int64(nOut)
		seekTarget = nOff + offset
		if nOut == 1024 {
			nOff, err = f.Seek(int64(nOut), 1)
			checkError(err)

			nextBuffer, newNOut = fetchBytes(f, 1024)

			seekTarget -= int64(nOut)
			nOut = newNOut
		} else {
			nOut = 0
		}

		nOff, err = f.Seek(seekTarget, 0)
		checkError(err)

		nWrit, err := f.Write(buff)
		checkError(err)
		if nWrit < len(buff) {
			logger.Warn("Ran out of space in the file while writing data, expected to write %d, but only wrote %d",
				len(buff), nWrit)
		}

		buff = nextBuffer
	}

	f.Truncate(totalSize)
}

// Replaces an existing line of data in the csv with a new one
func (cf *CSVFile) ModifyLine(line int, data []string) {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	f, err := os.OpenFile(cf.Filepath, os.O_RDWR, 0777)
	checkError(err)

	var newLength int64 = int64(len(data))
	for _, col := range data {
		newLength += int64(len(col))
	}

	_, err = f.Seek(cf.LineLocations[line], 0)
	checkError(err)

	oldLength := cf.LineLocations[line+1] - cf.LineLocations[line]

	diff := newLength - oldLength
	if diff != 0 && line < len(cf.LineLocations)-1 {
		start := cf.LineLocations[line+1]
		shiftFileContents(f, diff, start)
	}

	_, err = f.Seek(cf.LineLocations[line], 0)
	checkError(err)

	_, err = f.WriteString(strings.Join(data, ",") + "\n")
	checkError(err)

	// Update the line locations
	for li := int64(line + 1); li <= cf.LineCount; li++ {
		cf.LineLocations[li] += diff
	}

	f.Close()

	cf.syncLineCount()
}

// Deletes a single line of data from the csv
func (cf *CSVFile) DeleteLine(line int64) {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	f, err := os.OpenFile(cf.Filepath, os.O_RDWR, 0777)
	checkError(err)

	diff := cf.LineLocations[line+1] - cf.LineLocations[line]
	if line < cf.LineCount-1 {
		// It's not the last line
		shiftFileContents(f, -diff, cf.LineLocations[line+1])

		// Update the line locations
		for li := line + 1; li <= cf.LineCount; li++ {
			cf.LineLocations[li-1] = cf.LineLocations[li] - diff
		}
	}

	cf.LineLocations = cf.LineLocations[:len(cf.LineLocations)-1]
	cf.LineCount--

	f.Truncate(cf.LineLocations[len(cf.LineLocations)-1])

	f.Close()
	cf.syncLineCount()
}

// Deletes multiple lines of data from the csv
// sorts the lines first so that they're in descending order
func (cf *CSVFile) DeleteLines(lines []int64) {
	// Sort in Descending order
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] > lines[j]
	})

	for _, line := range lines {
		cf.DeleteLine(line)
	}
}

// Adds a new line of data to the csv file and returns it's line number
func (cf *CSVFile) AppendLine(data []string) int {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	f, err := os.OpenFile(cf.Filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	checkError(err)

	newIndex := len(cf.LineLocations) - 1

	lineLength, err := f.WriteString(strings.Join(data, ",") + "\n")
	checkError(err)
	cf.LineLocations = append(cf.LineLocations,
		cf.LineLocations[newIndex]+int64(lineLength))

	cf.LineCount++

	f.Close()
	cf.syncLineCount()
	return newIndex
}

// Writes the csv file's line count to the first line of the csv file
func (cf CSVFile) syncLineCount() {
	lineString := fmt.Sprintf("%d", cf.LineCount)
	f, err := os.OpenFile(cf.Filepath, os.O_RDWR, 0777)
	checkError(err)
	defer f.Close()

	oldData, _, _ := readLineRaw(f, []byte{})

	oldLength := len(oldData)
	newLength := len(lineString)

	if oldLength != newLength {
		diff := newLength - oldLength
		shiftFileContents(f, int64(diff), int64(newLength))
	}

	_, err = f.Seek(0, 0)
	checkError(err)

	_, err = f.WriteString(lineString)
	checkError(err)
}
