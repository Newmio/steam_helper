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
	cp1XOffset := rand.Intn(401) - 150
	cp1YOffset := rand.Intn(401) - 150
	cp2XOffset := rand.Intn(401) - 150
	cp2YOffset := rand.Intn(401) - 150

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

func SleepRandom(min int, max int) {
	if min > max {
		min, max = max, min
	}

	time.Sleep(time.Duration(rand.Intn(max-min+1)+min) * time.Millisecond)
}

type Position struct {
	X int
	Y int
}

func MoveMouseAndWriteText(wb selenium.WebDriver, nameCssSelector, text string, startPosition Position) (Position, error) {
	element, err := wb.FindElement(selenium.ByCSSSelector, nameCssSelector)
	if err != nil {
		return Position{}, err
	}

	end, err := GetPositionElement(element)
	if err != nil {
		return Position{}, err
	}

	if err := MoveMouse(wb, startPosition.X, startPosition.Y, end.X, end.Y); err != nil {
		return Position{}, err
	}

	SleepRandom(100, 500)

	if err := element.Click(); err != nil {
		return Position{}, err
	}

	for _, value := range text {
		if err := element.SendKeys(string(value)); err != nil {
			return Position{}, err
		}

		SleepRandom(50, 150)
	}

	return end, nil
}

func MoveMouseAndClick(wb selenium.WebDriver, nameCssSelector string) (Position, error) {
	element, err := wb.FindElement(selenium.ByCSSSelector, nameCssSelector)
	if err != nil {
		return Position{}, err
	}

	start, err := GetRandomStartMousePosition(wb)
	if err != nil {
		return Position{}, err
	}

	end, err := GetPositionElement(element)
	if err != nil {
		return Position{}, err
	}

	return end, MoveMouse(wb, start.X, start.Y, end.X, end.Y)
}

func GetRandomStartMousePosition(wb selenium.WebDriver) (Position, error) {
	window, err := wb.ExecuteScript("return [window.innerWidth, window.innerHeight];", nil)
	if err != nil {
		return Position{}, err
	}

	windowSize := window.([]interface{})

	return Position{
		X: int(windowSize[0].(float64)),
		Y: int(windowSize[1].(float64)),
	}, nil
}

func GetPositionElement(element selenium.WebElement) (Position, error) {
	location, err := element.Location()
	if err != nil {
		return Position{}, err
	}

	size, err := element.Size()
	if err != nil {
		return Position{}, err
	}

	return Position{
		X: location.X + rand.Intn(int(size.Width)),
		Y: location.Y + rand.Intn(int(size.Height)),
	}, nil
}

// func convertToHTTPCookies(seleniumCookies []*selenium.Cookie) []*http.Cookie {
// 	var httpCookies []*http.Cookie
// 	for _, sc := range seleniumCookies {

// 		var expires time.Time
//         if sc.Expiry != 0 {
//             expires = time.Unix(int64(sc.Expiry), 0)
//         }

// 		hc := &http.Cookie{
// 			Name:     sc.Name,
// 			Value:    sc.Value,
// 			Domain:   sc.Domain,
// 			Path:     sc.Path,
// 			Expires:  expires,
// 			Secure:   sc.Secure,
// 			HttpOnly: sc.HTTPOnly,
// 		}
// 		httpCookies = append(httpCookies, hc)
// 	}
// 	return httpCookies
// }

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
