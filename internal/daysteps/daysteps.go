package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/es-x/4f/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parse := strings.Split(data, ",")

	if len(parse) != 2 {
		return 0, 0, errors.New("длина слайса меньше 2")
	}
	// конвертируем шаги в число, если ошибка то вызываем ошибку
	steps, err := strconv.Atoi(parse[0])
	if err != nil {
		return 0, 0, errors.New("ошибка преобразования в целое число")
	}

	if steps <= 0 {
		return 0, 0, errors.New("число шагов меньше или равно 0")
	}

	duration, err := time.ParseDuration(parse[1])
	if err != nil {
		return 0, 0, err
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	// проверяем, вернулась ли ошибка
	if err != nil {
		// если да, то выводим ее в консоль и завершаем программу
		fmt.Println(err.Error())
		return ""
	}

	if steps < 0 {
		return ""
	}

	dist := (float64(steps) * stepLength) / mInKm

	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, duration, calories)

	return result
}
