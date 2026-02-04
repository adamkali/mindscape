package schemas

import (
	"embed"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"github.com/google/uuid"
)

type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WidgetLayout struct {
	DefaultSize WidgetSize `json:"defaultSize"`
	MinSize     WidgetSize `json:"minSize"`
	MaxSize     WidgetSize `json:"maxSize"`
	Resizable   bool       `json:"resizable"`
}

type WidgetProperty struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Value       any    `json:"value"` // Changed from string to any to support different value types (string, bool, number)
	Format      string `json:"format"`
	Description string `json:"description"`
	Enum        []any  `json:"enum"`
}

type WidgetSchema struct {
	Description string                       `json:"description"`
	ID          uuid.UUID                    `json:"id"`
	Type        string                       `json:"type"`
	Title       string                       `json:"title"`
	Layout      WidgetLayout                 `json:"layout"`
	Properties  map[string]WidgetProperty    `json:"properties"`
	Required    []string                     `json:"required"`
}

type WidgetSchemaStorage struct {
	Storage map[uuid.UUID]*WidgetSchema
}

var (
	WidgetAdditionError = "Could not add widget to internal storage"
	WidgetNotFoundError = "Could not find widget in internal storage"
)

func BuiltinError(source string, err string) error {
	return fmt.Errorf("%s: %s", source, err)
}

func (wss *WidgetSchemaStorage) Add(schema WidgetSchema) error {
	for _, v := range wss.Storage {
		if v.ID == schema.ID {
			// get the type of the function called
			methodName := runtime.FuncForPC(
				reflect.ValueOf(wss.Add).Pointer(),
			).Name()
			return BuiltinError(methodName, WidgetAdditionError)
		}
	}
	wss.Storage[schema.ID] = &schema
	return nil
}

func (wss *WidgetSchemaStorage) Get(widgetID uuid.UUID) (*WidgetSchema, error) {
	widget := wss.Storage[widgetID]
	if widget == nil {
		methodName := runtime.FuncForPC(
			reflect.ValueOf(wss.Get).Pointer(),
		).Name()
		return nil, BuiltinError(methodName, WidgetNotFoundError)
	}
	return widget, nil
}

func (wss *WidgetSchemaStorage) GetAll() []WidgetSchema {
	storage := make([]WidgetSchema, len(wss.Storage))
	index := 0
	for _, v := range wss.Storage {
		storage[index] = *v
		index++
	}
	return storage
}

//go:embed searchbar-schema.json githubprofile-wide.json githubprofile-lg.json githubprofile-sm.json coolify-schema.json
var fs embed.FS

const mapping = `[
	"searchbar-schema.json",
	"githubprofile-wide.json",
	"githubprofile-lg.json",
	"githubprofile-sm.json",
	"coolify-schema.json"
]`

func EmbeddedScemas() (*WidgetSchemaStorage, error) {
	embeddedScemaFilenames := new ([]string)
	err := json.Unmarshal([]byte(mapping), &embeddedScemaFilenames) 
	if err != nil {
		return nil, err
	}
	schemas := WidgetSchemaStorage{
		Storage: make(
			map[uuid.UUID]*WidgetSchema,
			len(*embeddedScemaFilenames),
		),
	}
	for _, filename:= range *embeddedScemaFilenames {
		data, err := fs.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		var schema WidgetSchema
		err = json.Unmarshal(data, &schema)
		if err != nil {
			return nil, err
		}
		err = schemas.Add(schema)
		if err != nil {
			return nil, err
		}
	}

	return &schemas, nil
}
