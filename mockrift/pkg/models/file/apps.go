package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mockrift/pkg/models"
	"os"
	"path/filepath"
)

type AppModel struct{}

func (m *AppModel) GetAll() []*models.App {
	appFiles, gErr := filepath.Glob("./requests/*.json")
	if gErr != nil {
		log.Fatal("Error loading request files: " + gErr.Error())
	}

	var apps []*models.App

	for _, appFile := range appFiles {
		appBytes, rErr := ioutil.ReadFile(appFile)
		if rErr != nil {
			log.Fatal("Unable to read app file: " + rErr.Error())
		}

		var app models.App
		uErr := json.Unmarshal(appBytes, &app)
		if uErr != nil {
			log.Fatal("Unable to unmarshal app file file: " + uErr.Error())
		}

		apps = append(apps, &app)
	}

	return apps
}

func (m *AppModel) Get(slug string) *models.App {
	var a models.App

	fmt.Println("Loading app from /home/appuser/app/requests/" + slug + ".json")
	jsonFile, err := os.Open("/home/appuser/app/requests/" + slug + ".json")
	if err != nil {
		// If the file doesn't exist then create a new app and return it.
		// It will be saved as soon as it has a request to save.
		return &models.App{
			Slug: slug,
		}
	}
	defer jsonFile.Close()

	jsonBytes, jsonBytesErr := ioutil.ReadAll(jsonFile)
	if jsonBytesErr != nil {
		log.Fatal(fmt.Printf("Unable to read JSON file (%s): %s\n", slug, jsonBytesErr.Error()))
	}

	unmarshalErr := json.Unmarshal(jsonBytes, &a)
	if unmarshalErr != nil {
		log.Fatal("Unable to unmarshal json file: " + unmarshalErr.Error())
	}

	return &a
}

func (m *AppModel) Save(app *models.App) {
	appJson, mErr := json.MarshalIndent(app, "", "  ")
	if mErr != nil {
		log.Fatal(mErr)
	}

	f, oErr := os.OpenFile("./requests/"+app.Slug+".json", os.O_WRONLY|os.O_CREATE, 0644)
	if oErr != nil {
		log.Fatal("Unable to open file for writing: " + oErr.Error())
	}

	_, wErr := f.Write(appJson)
	if wErr != nil {
		log.Fatal("Unable to write json to file: " + wErr.Error())
	}
}
