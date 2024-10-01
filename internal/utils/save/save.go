package save

import (
	"fmt"
	"os"
)

func SaveToFile(filename, data string) error {
	// Открываем файл для записи (создаём, если его нет)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// Записываем данные в файл
	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("ошибка записи данных в файл: %w", err)
	}

	return nil
}
