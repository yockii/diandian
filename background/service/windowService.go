package service

import "changeme/background/app"

type WindowService struct{}

func (s *WindowService) HideMainAndShowFloating() {
	app.ShowFloating()
}

func (s *WindowService) ShowMainWindow() {
	app.ShowMain()
}

func (s *WindowService) IsInitializeSuccess() bool {
	return app.IsInitializeSuccess()
}

func (s *WindowService) ShowSettings() {
	app.ShowSettings()
}

func (s *WindowService) HideSettings() {
	app.HideSettings()
}

func (s *WindowService) FloatingStickySide() int {
	return app.FloatingStickySide()
}
