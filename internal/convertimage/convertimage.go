package convertimage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Img2word(imgPath *string, ocrPath *string) (string, error) {
	absImgPath, err := filepath.Abs(*imgPath)
	if err != nil {
		return "", err
	}
	absOutputPath, err := filepath.Abs(*ocrPath)

	preCmd := fmt.Sprintf("shortcuts run \"extract_text_from_image\" -i %s -o %s", absImgPath, absOutputPath)
	cmd := exec.Command("zsh", "-c", preCmd)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// read tmp.txt as result
	results := ""
	file, err := os.Open(absOutputPath)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		results += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return "", err
	}

	return results, nil
}
