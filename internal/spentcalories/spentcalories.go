package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parse := strings.Split(data, ",")

	if len(parse) != 3 {
		return 0, "", 0, errors.New("длина слайса меньше 2")
	}

	steps, err := strconv.Atoi(parse[0])
	if err != nil {
		return 0, "", 0, errors.New("ошибка преобразования в целое число")
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("число шагов меньше или равно 0")
	}

	duration, err := time.ParseDuration(parse[2])
	if err != nil {
		return 0, "", 0, err
	}
	return steps, parse[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return (height * stepLengthCoefficient * float64(steps)) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration < 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, style, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	var result string
	switch style {
	case "Ходьба":
		sc, _ := WalkingSpentCalories(steps, weight, height, duration)
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		result = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость:  %.2f км/ч\nСожгли калорий: %.2f ", style, duration, distance, meanSpeed, sc)
	case "Бег":
		sc, _ := RunningSpentCalories(steps, weight, height, duration)
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		result = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость:  %.2f км/ч\nСожгли калорий: %.2f ", style, duration, distance, meanSpeed, sc)
	default:
		fmt.Println("неизвестный тип тренировки")
	}
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 {
		return 0, errors.New("число меньше или равно 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 {
		return 0, errors.New("число меньше или равно 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return ((weight * meanSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
