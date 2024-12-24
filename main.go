package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Определение структуры для хранения параметров
type Params struct {
	Fields    []int
	Delimiter string
	Separated bool
}

// Основная функция
func main() {
	// Обрабатываем введенную строку и разбиваем ее на параметры
	params := setParams()
	processInput(params)
}

// Функция для обработки введеной строки
func setParams() Params {
	// Определение параметров
	fields := flag.String("f", "", "Выбрать поля (колонки), например, 1,2,3")
	delimiter := flag.String("d", "\t", "Использовать другой разделитель")
	separated := flag.Bool("s", false, "Только строки с разделителем")

	flag.Parse()

	// Проверка на то что пользователь не ввел поля
	if *fields == "" {
		fmt.Fprintln(os.Stderr, "Необходимо указать поля через ключ -f")
		os.Exit(1)
	}

	return Params{
		Fields:    parseFields(*fields),
		Delimiter: *delimiter,
		Separated: *separated,
	}
}

func parseFields(fields string) []int {
	parts := strings.Split(fields, ",")
	var indices []int
	for _, part := range parts {
		var index int
		_, err := fmt.Sscanf(part, "%d", &index)
		if err != nil || index <= 0 {
			fmt.Fprintln(os.Stderr, "Некорректное значение поля:", part)
			os.Exit(1)
		}
		indices = append(indices, index-1) // Индексы начинаются с 0
	}
	return indices
}

func processInput(config Params) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if shouldSkipLine(line, config) {
			continue
		}
		columns := strings.Split(line, config.Delimiter)
		output := selectFields(columns, config.Fields)
		if output != "" {
			fmt.Println(output)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения ввода:", err)
		os.Exit(1)
	}
}

func shouldSkipLine(line string, config Params) bool {
	return config.Separated && !strings.Contains(line, config.Delimiter)
}

func selectFields(columns []string, indices []int) string {
	var selected []string
	for _, index := range indices {
		if index < len(columns) {
			selected = append(selected, columns[index])
		}
	}
	return strings.Join(selected, "\t")
}
