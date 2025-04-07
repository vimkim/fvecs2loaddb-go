package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"
)

// readFvecs는 fvecs 파일을 열어 각 벡터를 []float32 슬라이스로 읽어들인다.
func readFvecs(filename string) ([][]float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 실패: %v", err)
	}
	defer file.Close()

	var vectors [][]float32
	// 파일 끝까지 반복 읽기
	for {
		var dim int32
		err := binary.Read(file, binary.LittleEndian, &dim)
		if err != nil {
			// 파일 끝에 도달하면 break
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("차원 읽기 실패: %v", err)
		}

		// 차원이 음수이면 잘못된 파일
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

func main() {
	// 명령행 인자: fvecs 파일, 테이블 이름, 출력 파일 (선택)
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

	// fvecs 파일 읽기
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

	// loaddb 파일 작성
	// %class 명령어를 이용해 테이블(클래스) 및 vec 칼럼을 지정한다.
	header := fmt.Sprintf("%%class [%s].[%s] (vec)\n", userName, tableName)
	_, err = outFile.WriteString(header)
	if err != nil {
		fmt.Printf("헤더 작성 오류: %v\n", err)
		os.Exit(1)
	}

	// 각 벡터를 데이터 라인으로 작성: 인스턴스 번호 없이 벡터만 출력
	for _, vec := range vectors {
		strValues := make([]string, len(vec))
		for i, v := range vec {
			// 소수점 이하 6자리까지 문자열로 변환
			strValues[i] = fmt.Sprintf("%.6f", v)
		}
		// 각 벡터를 콤마와 공백으로 구분하여 대괄호([])로 감싼다.
		vectorStr := fmt.Sprintf("[%s]", strings.Join(strValues, ", "))
		// 데이터 라인: 벡터 문자열을 작은따옴표(')로 감싼다.
		line := fmt.Sprintf("'%s'\n", vectorStr)
		_, err = outFile.WriteString(line)
		if err != nil {
			fmt.Printf("데이터 라인 작성 오류: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("로딩 파일이 성공적으로 생성되었습니다: %s\n", outputPath)
}
