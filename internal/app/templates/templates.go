package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

// Storage for templates
type Storage struct {
	templates map[string]*template.Template
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		templates: make(map[string]*template.Template),
	}
}

// LoadDir ...
func (s *Storage) LoadDir(path string, ext string) error {
	dirInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("could not load templates dir: %v", err)
	}

	validFilename, err := regexp.Compile(`[a-z_\.]+`)
	if err != nil {
		return fmt.Errorf("could not load templates dir: %v", err)
	}

	for _, fileInfo := range dirInfo {
		fileExt := filepath.Ext(fileInfo.Name())
		fileName := fileInfo.Name()

		if !fileInfo.IsDir() && validFilename.MatchString(fileName) && fileExt == ext {
			f, err := os.Open(path + "/" + fileName)
			if err != nil {
				return fmt.Errorf("could not load template '%s': %v", fileName, err)
			}
			defer f.Close()
			fileContent, err := ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("could not load template '%s': %v", fileName, err)
			}

			tmpl, err := template.New("message_template").Parse(string(fileContent))
			if err != nil {
				return fmt.Errorf("could not parse template '%s': %v", fileName, err)
			}

			s.templates[fileName[0:len(fileName)-len(fileExt)]] = tmpl
		}
	}

	return nil
}

// Get return template object
func (s *Storage) Get(title string) (*template.Template, error) {
	tmpl, ok := s.templates[title]
	if !ok {
		return nil, fmt.Errorf("no such template '%s'", title)
	}
	return tmpl, nil
}
