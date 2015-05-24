package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type ConfigWriter struct {
	basePath string
}

func NewConfigWriter(basePath string) *ConfigWriter {
	return &ConfigWriter{
		basePath: basePath,
	}
}

func (configWriter *ConfigWriter) WriteSwitchVcl(entries []*ConfigEntry) {
	vcl := ""

	for _, entry := range entries {
		vcl = vcl + configWriter.getEntryVcl(entry)
	}

	outputDir := path.Join(configWriter.basePath, "output")

	err := os.MkdirAll(outputDir, 0775)
	if err != nil {
		log.Fatal(err)
	}

	switchVclFile := path.Join(outputDir, "switch.vcl")

	err = ioutil.WriteFile(switchVclFile, []byte(vcl), 0664)
	if err != nil {
		log.Fatal(err)
	}
}

func (ConfigWriter) getEntryVcl(entry *ConfigEntry) string {
	template := `
if (req.http.host == "{{HOST}}") {
  set req.backend_hint = {{BACKEND}};
}
	`

	vcl := template

	vcl = strings.Replace(vcl, "{{HOST}}", entry.Host, -1)
	vcl = strings.Replace(vcl, "{{BACKEND}}", entry.Backend, -1)

	return vcl
}

func (configWriter *ConfigWriter) ApplyConfiguration() error {
	scriptPath := path.Join(basePath, "varnish-restart.sh")

	cmd := exec.Command(scriptPath)

	return cmd.Start()
}
