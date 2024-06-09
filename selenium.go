package steam_helper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tebeka/selenium"
)

func bezierCurve(t float64, p0, p1, p2, p3 float64) float64 {
	return (1-t)*(1-t)*(1-t)*p0 + 3*(1-t)*(1-t)*t*p1 + 3*(1-t)*t*t*p2 + t*t*t*p3
}

func MoveMouse(wd selenium.WebDriver, startX, startY, endX, endY int) error {
	steps := 100 // Количество шагов для перемещения
	delay := 2 * time.Millisecond

	// Генерация случайных контрольных точек для кривой Безье
	cp1XOffset := rand.Intn(601) - 150
	cp1YOffset := rand.Intn(601) - 150
	cp2XOffset := rand.Intn(601) - 150
	cp2YOffset := rand.Intn(601) - 150

	cp1X := startX + cp1XOffset
	cp1Y := startY + cp1YOffset
	cp2X := endX + cp2XOffset
	cp2Y := endY + cp2YOffset

	// Создание маркера курсора
	_, err := wd.ExecuteScript(`(function() {
        var marker = document.createElement('div');
        marker.id = 'cursorMarker';
        marker.style.position = 'absolute';
        marker.style.width = '10px';
        marker.style.height = '10px';
        marker.style.backgroundColor = 'red';
        marker.style.borderRadius = '50%';
        marker.style.zIndex = '10000';
        document.body.appendChild(marker);
    })();`, nil)
	if err != nil {
		return err
	}

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := bezierCurve(t, float64(startX), float64(cp1X), float64(cp2X), float64(endX))
		y := bezierCurve(t, float64(startY), float64(cp1Y), float64(cp2Y), float64(endY))

		script := fmt.Sprintf(`(function() {
            var marker = document.getElementById('cursorMarker');
            marker.style.left = '%dpx';
            marker.style.top = '%dpx';
        })();`, int(x), int(y))

		if _, err := wd.ExecuteScript(script, nil); err != nil {
			return err
		}

		time.Sleep(delay + time.Duration(rand.Intn(5))*time.Millisecond)
	}

	return nil
}
