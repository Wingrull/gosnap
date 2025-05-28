package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Config struct {
	exclude      StringSlice // для флага -e/--exclude
	excludeNoise bool        // флаг --exclude-noise
	output       string      // путь к выходному файлу
	extensions   StringSlice // для флага -ext/--extension
}

type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	// Инициализация логгера с временными метками
	log.SetFlags(log.LstdFlags)

	// Логируем сырые аргументы командной строки для отладки
	log.Printf("Raw command-line arguments: %v", os.Args)

	// Определение флагов
	var config Config
	var exclude StringSlice
	var extensions StringSlice
	flag.Var(&exclude, "e", "Manually exclude specific files or folders")
	flag.Var(&exclude, "exclude", "Manually exclude specific files or folders")
	flag.BoolVar(&config.excludeNoise, "en", true, "Automatically exclude common development artifacts (default: true)")
	flag.BoolVar(&config.excludeNoise, "exclude-noise", true, "Automatically exclude common development artifacts (default: true)")
	flag.StringVar(&config.output, "o", "snap.txt", "Output file path")
	flag.StringVar(&config.output, "output", "snap.txt", "Output file path")
	flag.Var(&extensions, "ext", "Include only files with specified extensions (e.g., .py, .go)")
	flag.Var(&extensions, "extension", "Include only files with specified extensions (e.g., .py, .go)")
	flag.Parse()

	config.exclude = exclude
	config.extensions = extensions

	// Установка входной директории
	args := flag.Args()
	inputDir := "."
	if len(args) > 0 {
		inputDir = args[0]
		log.Printf("Input directory: %s (specified)", inputDir)
	} else {
		log.Printf("Input directory: %s (default)", inputDir)
	}
	log.Printf("Output file: %s", config.output)
	log.Printf("Exclude patterns: %v", config.exclude)
	log.Printf("Include extensions: %v", config.extensions)
	log.Printf("Exclude noise: %v", config.excludeNoise)

	// Создаем выходной файл
	log.Printf("Creating output file: %s", config.output)
	outputFile, err := os.Create(config.output)
	if err != nil {
		log.Printf("Error creating output file: %v", err)
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Собираем структуру директорий и содержимое
	log.Println("Starting snapshot generation")
	err = generateSnapshot(inputDir, outputFile, config)
	if err != nil {
		log.Printf("Error generating snapshot: %v", err)
		fmt.Printf("Error generating snapshot: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Snapshot successfully written to %s", config.output)
}

func generateSnapshot(inputDir string, outputFile *os.File, config Config) error {
	// Список стандартных артефактов для исключения
	noisePatterns := []string{".git", ".venv", "__pycache__", "node_modules", ".idea", ".DS_Store", "lib", "test", "etc", "log", "tools", ".md"}

	// Собираем структуру директорий
	var structure strings.Builder
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверяем исключения
		relPath, _ := filepath.Rel(inputDir, path)
		if shouldExclude(relPath, info, config.exclude, config.excludeNoise, noisePatterns, config.extensions) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Формируем отступы для структуры
		depth := len(strings.Split(relPath, string(os.PathSeparator)))
		indent := strings.Repeat("  ", depth-1)

		// Добавляем в структуру
		if info.IsDir() {
			structure.WriteString(fmt.Sprintf("%s%s/\n", indent, info.Name()))
		} else {
			structure.WriteString(fmt.Sprintf("%s%s\n", indent, info.Name()))
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Записываем структуру в файл
	log.Println("Writing directory structure to output file")
	_, err = outputFile.WriteString("Directory Structure:\n")
	if err != nil {
		return err
	}
	_, err = outputFile.WriteString(structure.String())
	if err != nil {
		return err
	}
	_, err = outputFile.WriteString("\n=== File Contents ===\n\n")
	if err != nil {
		return err
	}

	// Собираем содержимое файлов
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории и исключенные файлы
		relPath, _ := filepath.Rel(inputDir, path)
		if info.IsDir() || shouldExclude(relPath, info, config.exclude, config.excludeNoise, noisePatterns, config.extensions) {
			return nil
		}

		// Проверяем, является ли файл бинарным или не в UTF-8
		log.Printf("Checking file: %s", relPath)
		isBinary, err := isBinaryFile(path)
		if err != nil {
			log.Printf("Error checking file %s: %v", relPath, err)
			return nil // Пропускаем файл, но продолжаем обработку
		}
		if isBinary {
			log.Printf("Skipping binary or non-UTF-8 file: %s", relPath)
			return nil
		}

		// Читаем содержимое файла
		log.Printf("Reading file: %s", relPath)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", relPath, err)
			return err
		}

		// Записываем содержимое в выходной файл
		log.Printf("Writing file content: %s", relPath)
		_, err = outputFile.WriteString(fmt.Sprintf("File: %s\n", relPath))
		if err != nil {
			return err
		}
		_, err = outputFile.WriteString(string(content))
		if err != nil {
			return err
		}
		_, err = outputFile.WriteString("\n\n")
		return err
	})

	return err
}

func shouldExclude(path string, info os.FileInfo, exclude []string, excludeNoise bool, noisePatterns []string, extensions StringSlice) bool {
	// Проверяем явно исключенные файлы/папки
	for _, excl := range exclude {
		if strings.Contains(path, excl) || info.Name() == excl {
			return true
		}
	}

	// Проверяем стандартные артефакты, если включен excludeNoise
	if excludeNoise {
		for _, pattern := range noisePatterns {
			if strings.Contains(path, pattern) || info.Name() == pattern {
				return true
			}
		}
	}

	// Проверяем расширения файлов, если они указаны
	if len(extensions) > 0 && !info.IsDir() {
		hasValidExtension := false
		for _, ext := range extensions {
			if strings.HasSuffix(strings.ToLower(info.Name()), strings.ToLower(ext)) {
				hasValidExtension = true
				break
			}
		}
		if !hasValidExtension {
			return true
		}
	}

	return false
}

func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Читаем первые 1024 байта для анализа
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}
	buffer = buffer[:n]

	// Проверяем наличие нулевых байтов (характерно для бинарных файлов)
	for _, b := range buffer {
		if b == 0 {
			return true, nil
		}
	}

	// Проверяем, является ли содержимое валидным UTF-8
	utf8Reader := transform.NewReader(bytes.NewReader(buffer), unicode.UTF8.NewDecoder())
	_, err = io.ReadAll(utf8Reader)
	if err != nil {
		return true, nil // Считаем не-UTF-8 файлы бинарными
	}

	return false, nil
}