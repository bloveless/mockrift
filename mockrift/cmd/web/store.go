package main

import "mockrift/pkg/models"

func getAppFromApps(appName string, apps []*models.App) *models.App {
	for _, app := range apps {
		if app.Name == appName {
			return app
		}
	}

	return nil
}
