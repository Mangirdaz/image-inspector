package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mangirdaz/ocp-demo/config"

	iicmd "github.com/mangirdaz/image-inspector/pkg/cmd"
	ii "github.com/mangirdaz/image-inspector/pkg/inspector"
	"github.com/mangirdaz/image-inspector/pkg/storage"
)

func scanImage(image string) {

	inspectorOptions := iicmd.NewDefaultImageInspectorOptions()

	inspectorOptions.Image = image
	inspectorOptions.DstPath = config.Get("EnvImagePath")
	inspectorOptions.ScanType = "openscap"
	inspectorOptions.OpenScapHTML = true

	if err := inspectorOptions.Validate(); err != nil {
		log.Fatal(err)
	}

	storage := storage.InitKVStorage()
	inspector := ii.NewDefaultImageInspector(*inspectorOptions, storage)
	inspector.Inspect()

}
