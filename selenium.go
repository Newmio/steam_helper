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

func GetRandomWindowSize(userAgent string) ScreenSize {
	for _, value := range pcUserAgents {
		if value == userAgent {
			return computerScreenSizes[rand.Intn(len(computerScreenSizes))]
		}
	}

	return mobileAndTabletScreenSizes[rand.Intn(len(mobileAndTabletScreenSizes))]
}

type ScreenSize struct {
	Width  int
	Height int
}

var computerScreenSizes = []ScreenSize{
	{1920, 1080}, {1366, 768}, {1440, 900}, {1536, 864},
	{1600, 900}, {1680, 1050}, {2560, 1440}, {3840, 2160},
	{1280, 1024}, {1600, 1200}, {1920, 1200}, {2560, 1600},
	{3440, 1440}, {5120, 2880}, {2048, 1080}, {2880, 1800},
	{3200, 1800}, {2736, 1824}, {3000, 2000}, {3840, 1600},
	{4096, 2160}, {1080, 1920}, {1440, 2560}, {2160, 3840},
	{2400, 1350}, {3072, 1920}, {1600, 2560}, {1200, 1920},
}

var mobileAndTabletScreenSizes = []ScreenSize{
	{320, 480}, {360, 640}, {375, 667}, {375, 812},
	{414, 896}, {768, 1024}, {800, 1280}, {1024, 1366},
	{600, 1024}, {800, 1280}, {720, 1280}, {1080, 1920},
	{1440, 2560}, {2160, 3840}, {480, 800}, {540, 960},
}
