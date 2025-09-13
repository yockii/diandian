package service

import "changeme/background/app"

type WindowService struct{}

func (s *WindowService) HideMainAndShowFloating() {
	app.ShowFloating()
}

func (s *WindowService) ShowMainWindow() {
	app.ShowMain()
}
