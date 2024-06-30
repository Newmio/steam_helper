package steam_helper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tebeka/selenium"
)

var BuzierOffset int
var BuzierSteps int

func bezierCurve(t, p0, p1, p2, p3 float64) float64 {
	return (1-t)*(1-t)*(1-t)*p0 + 3*(1-t)*(1-t)*t*p1 + 3*(1-t)*t*t*p2 + t*t*t*p3
}

func MoveMouse(element selenium.WebElement, startX, startY, endX, endY int) error {

	// Генерация случайных контрольных точек для кривой Безье
	cp1XOffset := rand.Intn(BuzierOffset) - 100
	cp1YOffset := rand.Intn(BuzierOffset) - 100
	cp2XOffset := rand.Intn(BuzierOffset) - 100
	cp2YOffset := rand.Intn(BuzierOffset) - 100

	cp1X := startX + cp1XOffset
	cp1Y := startY + cp1YOffset
	cp2X := endX + cp2XOffset
	cp2Y := endY + cp2YOffset

	for i := 0; i <= BuzierSteps; i++ {
		t := float64(i) / float64(BuzierSteps)
		x := bezierCurve(t, float64(startX), float64(cp1X), float64(cp2X), float64(endX))
		y := bezierCurve(t, float64(startY), float64(cp1Y), float64(cp2Y), float64(endY))

		if err := element.MoveTo(int(x), int(y)); err != nil {
			return err
		}

		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
	}

	return nil
}

func TestMoveMouse(wd selenium.WebDriver, element selenium.WebElement, startX, startY, endX, endY int) error {
	// Генерация случайных контрольных точек для кривой Безье
	cp1XOffset := rand.Intn(BuzierOffset) - 100
	cp1YOffset := rand.Intn(BuzierOffset) - 100
	cp2XOffset := rand.Intn(BuzierOffset) - 100
	cp2YOffset := rand.Intn(BuzierOffset) - 100

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

	for i := 0; i <= BuzierSteps; i++ {
		t := float64(i) / float64(BuzierSteps)
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

		if err := element.MoveTo(int(x), int(y)); err != nil {
			return err
		}

		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
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

func TestMoveMouseAndWriteText(wd selenium.WebDriver, element selenium.WebElement, startPosition Position, text string) (Position, error) {
	end, err := GetPositionElement(element)
	if err != nil {
		return Position{}, err
	}

	if err := TestMoveMouse(wd, element, startPosition.X, startPosition.Y, end.X, end.Y); err != nil {
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

		SleepRandom(100, 250)
	}

	return end, nil
}

func MoveMouseAndWriteText(element selenium.WebElement, startPosition Position, text string) (Position, error) {
	end, err := GetPositionElement(element)
	if err != nil {
		return Position{}, err
	}

	if err := MoveMouse(element, startPosition.X, startPosition.Y, end.X, end.Y); err != nil {
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

		SleepRandom(100, 250)
	}

	return end, nil
}

func TestMoveMouseAndClick(wd selenium.WebDriver, element selenium.WebElement, startPosition Position) (Position, error) {
	end, err := GetPositionElement(element)
	if err != nil {
		return Position{}, err
	}

	if err := scrollToElement(element, end); err != nil {
		return Position{}, err
	}

	if err := TestMoveMouse(wd, element, startPosition.X, startPosition.Y, end.X, end.Y); err != nil {
		return Position{}, err
	}

	if err := element.Click(); err != nil {
		return Position{}, err
	}

	return end, nil
}

func MoveMouseAndClick(element selenium.WebElement, startPosition Position) (Position, error) {
	end, err := GetPositionElement(element)
	if err != nil {
		return Position{}, err
	}

	if err := MoveMouse(element, startPosition.X, startPosition.Y, end.X, end.Y); err != nil {
		return Position{}, err
	}

	if err := element.Click(); err != nil {
		return Position{}, err
	}

	return end, nil
}

func scrollToElement(element selenium.WebElement, position Position)error {
	window := activeScreenSize
	for{
		switch elementInWindow(position, window){

		case "right":
			if err := element.SendKeys(selenium.RightArrowKey); err != nil {
				return err
			}

			window.Width += 10

		case "down":
			if err := element.SendKeys(selenium.DownArrowKey); err != nil {
				return err
			}

			window.Height += 10

		case "stop":
			return nil

		default:
			continue
		}

		SleepRandom(1, 5)
	}
}

func elementInWindow(p Position, windowSize ScreenSize)string {

	if windowSize.Width >= p.X && p.X >= 0 && windowSize.Height >= p.Y && p.Y >= 0{
		return "stop"
	}

	if windowSize.Width < p.X && p.X >= 0 {
		return "right"
	}

	if windowSize.Height < p.Y && p.Y >= 0 {
		return "down"
	}

	return ""
}

func GetStartMousePosition(wd selenium.WebDriver) (Position, error) {
	window, err := wd.ExecuteScript("return [window.innerWidth, window.innerHeight];", nil)
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

func GetRandomWindowSize() ScreenSize {
	activeScreenSize = screenSizes[rand.Intn(len(screenSizes))]
	return activeScreenSize
}

type ScreenSize struct {
	Width  int
	Height int
}

var activeScreenSize ScreenSize

var screenSizes = []ScreenSize{
    //{Width: 640, Height: 480},     // VGA
    //{Width: 800, Height: 600},     // SVGA
    //{Width: 1024, Height: 768},    // XGA
    {Width: 1280, Height: 720},    // HD
    {Width: 1280, Height: 800},    // WXGA
    {Width: 1280, Height: 1024},   // SXGA
    {Width: 1366, Height: 768},    // WXGA
    {Width: 1440, Height: 900},    // WXGA+
    {Width: 1600, Height: 900},    // HD+
    {Width: 1680, Height: 1050},   // WSXGA+
    {Width: 1920, Height: 1080},   // Full HD
    {Width: 1920, Height: 1200},   // WUXGA
    {Width: 2048, Height: 1080},   // 2K
    {Width: 2560, Height: 1440},   // QHD
    {Width: 2560, Height: 1600},   // WQXGA
    //{Width: 2880, Height: 1800},   // Retina MacBook Pro
    {Width: 3200, Height: 1800},   // QHD+
    {Width: 3840, Height: 2160},   // 4K UHD
    {Width: 4096, Height: 2160},   // DCI 4K
    {Width: 5120, Height: 2880},   // 5K
}
