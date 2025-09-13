package service

import "changeme/background/app"

type WindowService struct{}

func (w *WindowService) HideMainAndShowFloating() {
	app.ShowFloating()
}
