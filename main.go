package main

import (
	"fmt"
	"image/color"
	"net/url"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

var lastSize = 48

func main() {
	a := app.New()
	w := a.NewWindow("App Luces")
	w.Resize(windowSize(1))

	showMainScreen(w)
	w.ShowAndRun()
}

var currentStopChan chan bool

func showBlueScreen(w fyne.Window) {
	// Detener cualquier ciclo anterior
	if currentStopChan != nil {
		select {
		case currentStopChan <- true:
		default:
		}
	}

	// Crear un rectángulo azul que cubra toda la pantalla
	blueRect := canvas.NewRectangle(color.RGBA{0, 0, 255, 255})

	// Botón para volver a la pantalla principal
	backBtn := widget.NewButton("Volver", func() {
		// Detener el ciclo de colores
		if currentStopChan != nil {
			select {
			case currentStopChan <- true:
			default:
			}
		}
		showMainScreen(w)
	})

	// Crear el contenedor con el rectángulo de fondo y el botón
	content := container.NewBorder(
		container.NewHBox(layout.NewSpacer(), backBtn, layout.NewSpacer()),
		nil, nil, nil,
		blueRect,
	)

	w.SetContent(content)

	// Iniciar el ciclo de cambio de color cada segundo
	startColorCycle(blueRect)
}

func showMainScreen(w fyne.Window) {
	btnAzul := widget.NewButton("AZUL POLICIA", func() {
		showBlueScreen(w)
	})

	btnAmarillo := widget.NewButton("AMARILLO PELIGRO", func() {
		showYellowScreen(w)
	})

	btnTexto := widget.NewButton("TEXTO PARPADEANTE", func() {
		showTextInputScreen(w)
	})

	label := widget.NewLabel("Código fuente:")
	label.Alignment = fyne.TextAlignCenter
	l := widget.NewHyperlink("https://github.com/EironnESP/app_luces", &url.URL{Scheme: "https", Host: "github.com", Path: "/EironnESP/app_luces"})
	l.Alignment = fyne.TextAlignCenter

	infoContainer := container.NewVBox(layout.NewSpacer(), label, l, layout.NewSpacer())

	// Crear el layout principal con los tres botones ocupando toda la pantalla
	mainContent := container.NewGridWithColumns(1,
		btnAzul,
		btnAmarillo,
		btnTexto,
		infoContainer,
	)

	w.SetContent(mainContent)
}

func showYellowScreen(w fyne.Window) {
	// Detener cualquier ciclo anterior
	if currentStopChan != nil {
		select {
		case currentStopChan <- true:
		default:
		}
	}

	// Crear un rectángulo amarillo que cubra toda la pantalla
	yellowRect := canvas.NewRectangle(color.RGBA{255, 255, 0, 255})

	// Botón para volver a la pantalla principal
	backBtn := widget.NewButton("Volver", func() {
		// Detener el ciclo de colores
		if currentStopChan != nil {
			select {
			case currentStopChan <- true:
			default:
			}
		}
		showMainScreen(w)
	})

	// Crear el contenedor con el rectángulo de fondo y el botón
	content := container.NewBorder(
		container.NewHBox(layout.NewSpacer(), backBtn, layout.NewSpacer()),
		nil, nil, nil,
		yellowRect,
	)

	w.SetContent(content)

	// Iniciar el ciclo de cambio de color amarillo-negro cada 0.5 segundos
	startYellowColorCycle(yellowRect)
}

func showTextInputScreen(w fyne.Window) {
	// Detener cualquier ciclo anterior
	if currentStopChan != nil {
		select {
		case currentStopChan <- true:
		default:
		}
	}

	// Crear el label con el texto "Introduce el texto:"
	label := widget.NewLabel("Introduce el texto:")
	label.Alignment = fyne.TextAlignCenter

	// Crear la entrada de texto
	textEntry := widget.NewEntry()
	textEntry.SetPlaceHolder("Escribe tu texto aquí...")

	// Crear el label para el tamaño de texto
	sizeLabel := widget.NewLabel("Tamaño del texto (1-300):")
	sizeLabel.Alignment = fyne.TextAlignCenter

	// Crear la entrada para el tamaño de texto
	sizeEntry := widget.NewEntry()
	sizeEntry.SetPlaceHolder("48")
	sizeEntry.OnChanged = func(str string) {
		if i, err := strconv.Atoi(str); err != nil || i > 300 || i <= 0 { //comprobar si es numerico, no negativo y menor a 300 sec
			sizeEntry.Text = fmt.Sprint("", lastSize)
		} else {
			lastSize = i
		}
	}

	// Crear el botón "Listo"
	listoBtn := widget.NewButton("Listo", func() {
		text := textEntry.Text
		if text != "" {
			// Obtener el tamaño del texto, por defecto 48
			sizeStr := sizeEntry.Text
			if sizeStr == "" {
				sizeStr = "48"
			}

			size, err := strconv.Atoi(sizeStr)
			if err != nil || size < 1 || size > 300 {
				// Si hay error o el valor está fuera de rango, usar 48
				size = 48
			}

			showTextDisplayScreen(w, text, float32(size))
		}
	})

	// Crear botón para volver
	backBtn := widget.NewButton("Volver", func() {
		showMainScreen(w)
	})

	label2 := widget.NewLabel("!!! Pon el móvil en horizontal !!!")
	label2.Alignment = fyne.TextAlignCenter

	// Crear el layout con los elementos centrados
	content := container.NewVBox(
		layout.NewSpacer(),
		label,
		textEntry,
		sizeLabel,
		sizeEntry,
		listoBtn,
		backBtn,
		layout.NewSpacer(),
		label2,
	)

	w.SetContent(content)
}

func showTextDisplayScreen(w fyne.Window, text string, textSize float32) {
	// Detener cualquier ciclo anterior
	if currentStopChan != nil {
		select {
		case currentStopChan <- true:
		default:
		}
	}

	// Crear el label con el texto del usuario, con fuente muy grande
	textLabel := widget.NewLabel(text)
	textLabel.Alignment = fyne.TextAlignCenter
	textLabel.Wrapping = fyne.TextWrapWord
	textLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Crear un canvas.Text para tener más control sobre el tamaño
	textObj := canvas.NewText(text, color.RGBA{255, 255, 255, 255})
	textObj.Alignment = fyne.TextAlignCenter
	textObj.TextStyle = fyne.TextStyle{Bold: true}
	textObj.TextSize = textSize // Usar el tamaño especificado por el usuario

	// Botón para volver (pequeño, en la esquina superior)
	backBtn := widget.NewButton("Volver", func() {
		if currentStopChan != nil {
			select {
			case currentStopChan <- true:
			default:
			}
		}
		showMainScreen(w)
	})

	// Layout que ocupa toda la pantalla - el texto ocupa todo el espacio
	content := container.NewBorder(
		backBtn,                      // top
		nil,                          // bottom
		nil,                          // left
		nil,                          // right
		container.NewCenter(textObj), // center - texto centrado ocupando todo el espacio
	)

	w.SetContent(content)

	// Iniciar el parpadeo del texto cada segundo
	startTextBlinkCanvas(textObj, text)
}

func windowSize(part float32) fyne.Size {
	if screenshot.NumActiveDisplays() > 0 {
		bounds := screenshot.GetDisplayBounds(0) //0 = monitor 1
		return fyne.NewSize(float32(bounds.Dx())*part, float32(bounds.Dy())*part)
	}
	return fyne.NewSize(800, 800)
}

func startColorCycle(rect *canvas.Rectangle) {
	currentStopChan = make(chan bool, 1)

	timings := []struct {
		duration time.Duration
		isBlue   bool
	}{
		{1 * time.Second, false},        // negro: 1 sec
		{1 * time.Second, true},         // azul: 1 sec
		{1 * time.Second, false},        // negro: 1 sec
		{1 * time.Second, true},         // azul: 1 sec
		{1 * time.Second, false},        // negro: 1 sec
		{100 * time.Millisecond, true},  // azul: 0.1 sec
		{100 * time.Millisecond, false}, // negro: 0.1 sec
		{100 * time.Millisecond, true},  // azul: 0.1 sec
		{500 * time.Millisecond, false}, // negro: 0.5 sec
		{100 * time.Millisecond, true},  // azul: 0.1 sec
		{100 * time.Millisecond, false}, // negro: 0.1 sec
		{100 * time.Millisecond, true},  // azul: 0.1 sec
	}

	go func() {
		step := 0
		for {
			select {
			case <-currentStopChan:
				return
			default:
				// Obtener el timing actual
				currentTiming := timings[step%len(timings)]

				// Cambiar el color
				fyne.Do(func() {
					if currentTiming.isBlue {
						rect.FillColor = color.RGBA{0, 0, 255, 255}
					} else {
						rect.FillColor = color.RGBA{0, 0, 0, 255}
					}
					rect.Refresh()
				})

				// Esperar el tiempo correspondiente
				timer := time.NewTimer(currentTiming.duration)
				select {
				case <-timer.C:
					// Continuar al siguiente paso
				case <-currentStopChan:
					timer.Stop()
					return
				}

				step++
			}
		}
	}()
}

func startYellowColorCycle(rect *canvas.Rectangle) {
	currentStopChan = make(chan bool, 1)

	go func() {
		isYellow := true
		for {
			select {
			case <-currentStopChan:
				return
			default:
				// Cambiar el color
				fyne.Do(func() {
					if isYellow {
						// Amarillo
						rect.FillColor = color.RGBA{255, 255, 0, 255}
					} else {
						// Negro
						rect.FillColor = color.RGBA{0, 0, 0, 255}
					}
					rect.Refresh()
				})

				// Esperar 0.5 segundos
				timer := time.NewTimer(500 * time.Millisecond)
				select {
				case <-timer.C:
					// Cambiar al siguiente color
					isYellow = !isYellow
				case <-currentStopChan:
					timer.Stop()
					return
				}
			}
		}
	}()
}

func startTextBlinkCanvas(textObj *canvas.Text, originalText string) {
	currentStopChan = make(chan bool, 1)

	go func() {
		isVisible := true

		for {
			select {
			case <-currentStopChan:
				// Restaurar el texto antes de salir
				fyne.Do(func() {
					textObj.Text = originalText
					textObj.Refresh()
				})
				return
			default:
				// Cambiar la visibilidad del texto
				fyne.Do(func() {
					if isVisible {
						textObj.Text = ""
					} else {
						textObj.Text = originalText
					}
					textObj.Refresh()
				})

				// Esperar 1 segundo
				timer := time.NewTimer(1 * time.Second)
				select {
				case <-timer.C:
					isVisible = !isVisible
				case <-currentStopChan:
					timer.Stop()
					// Restaurar el texto antes de salir
					fyne.Do(func() {
						textObj.Text = originalText
						textObj.Refresh()
					})
					return
				}
			}
		}
	}()
}
