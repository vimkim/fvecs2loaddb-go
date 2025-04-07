package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

// readFvecs reads vectors from an fvecs file.
func readFvecs(filename string) ([][]float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 실패: %v", err)
	}
	defer file.Close()

	var vectors [][]float32
	for {
		var dim int32
		err := binary.Read(file, binary.LittleEndian, &dim)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("차원 읽기 실패: %v", err)
		}
		if dim <= 0 {
			return nil, fmt.Errorf("잘못된 차원 값: %d", dim)
		}

		vec := make([]float32, dim)
		err = binary.Read(file, binary.LittleEndian, &vec)
		if err != nil {
			return nil, fmt.Errorf("벡터 데이터 읽기 실패: %v", err)
		}
		vectors = append(vectors, vec)
	}
	return vectors, nil
}

// formatVector converts a vector of float32 to its formatted string.
func formatVector(vec []float32) string {
	strValues := make([]string, len(vec))
	for i, v := range vec {
		strValues[i] = fmt.Sprintf("%.6f", v)
	}
	vectorStr := fmt.Sprintf("[%s]", strings.Join(strValues, ", "))
	return fmt.Sprintf("'%s'\n", vectorStr)
}

func main() {
	var fvecsPath string
	var userName string
	var tableName string
	var outputPath string

	flag.StringVar(&fvecsPath, "fvecs", "", "입력 fvecs 파일 경로")
	flag.StringVar(&userName, "user", "dba", "User name")
	flag.StringVar(&tableName, "table", "", "로딩할 테이블 이름 (예: my_table)")
	flag.StringVar(&outputPath, "out", "out/loaddb_script.txt", "생성될 loaddb 문법 파일 경로")
	flag.Parse()

	if fvecsPath == "" || tableName == "" {
		fmt.Println("사용법: -fvecs 입력파일 -table 테이블이름 [-out 출력파일]")
		os.Exit(1)
	}

	// Read vectors from file.
	vectors, err := readFvecs(fvecsPath)
	if err != nil {
		fmt.Printf("fvecs 파일 읽기 오류: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("출력 파일 생성 오류: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Use buffered writer for efficient output.
	writer := bufio.NewWriter(outFile)

	// Write header.
	header := fmt.Sprintf("%%class [%s].[%s] (vec)\n", userName, tableName)
	if _, err = writer.WriteString(header); err != nil {
		fmt.Printf("헤더 작성 오류: %v\n", err)
		os.Exit(1)
	}

	// Prepare a slice to hold formatted lines.
	formattedLines := make([]string, len(vectors))

	// Use WaitGroup to ensure all goroutines complete.
	var wg sync.WaitGroup
	workerCount := 8 // Adjust based on available CPU cores.
	jobChan := make(chan struct {
		index int
		vec   []float32
	}, len(vectors))

	// Worker function.
	worker := func() {
		for job := range jobChan {
			formattedLines[job.index] = formatVector(job.vec)
			wg.Done()
		}
	}

	// Start worker goroutines.
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	// Dispatch jobs.
	for i, vec := range vectors {
		wg.Add(1)
		jobChan <- struct {
			index int
			vec   []float32
		}{index: i, vec: vec}
	}
	close(jobChan)
	wg.Wait()

	// Write formatted lines sequentially.
	for _, line := range formattedLines {
		if _, err = writer.WriteString(line); err != nil {
			fmt.Printf("데이터 라인 작성 오류: %v\n", err)
			os.Exit(1)
		}
	}

	// Flush the writer to ensure all data is written.
	if err = writer.Flush(); err != nil {
		fmt.Printf("출력 버퍼 플러시 오류: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("로딩 파일이 성공적으로 생성되었습니다: %s\n", outputPath)
}
