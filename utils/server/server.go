package server_utils

import (
	"encoding/json"
	"fmt"
	"io"
	"mud/utils/crud"
	"mud/utils/ui/logger"
	"net/http"
)

func ReadFullResponse(r io.ReadCloser) []byte {
	var buffer []byte = make([]byte, 1024)
	var result []byte
	var nRead int
	var err interface{}

	for nRead, err = io.ReadFull(r, buffer); nRead > 0 && err == nil; nRead, err = io.ReadFull(r, buffer) {
		result = append(result, buffer...)
	}

	if err == io.ErrUnexpectedEOF {
		result = append(result, buffer[:nRead]...)
	} else if err != nil && err != io.EOF {
		panic(err)
	}

	return result
}

func CreateRouteHandlers(dict map[string]func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			http.Error(w, "You didn't even send in a request", http.StatusBadRequest)
			return
		}

		logger.Info("%s %s", r.Method, r.URL.Path)

		handler, ok := dict[r.Method]
		if !ok {
			http.Error(w, fmt.Sprintf("%s method is not implemented", r.Method), http.StatusMethodNotAllowed)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				logger.Error(2, r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		handler(w, r)
	}
}

func CreateHandlers(dict map[string]map[string]func(http.ResponseWriter, *http.Request)) {
	for route := range dict {
		logger.Info("Creating HTTP handlers for %s", route)
		http.HandleFunc(route, CreateRouteHandlers(dict[route]))
	}
}

func DefaultUpdateParser(d []byte) (interface{}, interface{}) {
	var payload crud.CrudUpdate
	json.Unmarshal(d, &payload)
	return payload.Key, payload.NewData
}

type CrudParsers struct {
	Create   func([]byte) []interface{}
	Retrieve func([]byte) interface{}
	Update   func([]byte) (interface{}, interface{})
	Delete   func([]byte) interface{}
}

func DefaultErrorHandler(_ http.ResponseWriter, _ []byte) bool {
	return false
}

func DefaultUpdateErrorHandler(w http.ResponseWriter, d []byte) bool {
	var payload crud.CrudUpdate
	err := json.Unmarshal(d, &payload)
	if err != nil {
		http.Error(w, "Updates must follow update format", http.StatusBadRequest)
		logger.Error(err)
	}
	return err != nil
}

type errorHandler func(http.ResponseWriter, []byte) bool

type CrudErrorHandlers struct {
	Create   errorHandler
	Retrieve errorHandler
	Update   errorHandler
	Delete   errorHandler
}

func DefaultTranslator(v interface{}) interface{} {
	return v
}

type CrudTranslators struct {
	Create   func(interface{}) interface{}
	Retrieve func(interface{}) interface{}
	Update   func(interface{}) interface{}
}

func FetchRequestBody(r *http.Request) []byte {
	return ReadFullResponse(r.Body)
}

func HandleCrudData(data interface{}, responseTranslator func(interface{}) interface{}) string {
	response, err := json.Marshal(responseTranslator(data))
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return string(response)
}

func CreateCrudRoutes(crudObj crud.Crud, parsers CrudParsers, errors CrudErrorHandlers, translators CrudTranslators) map[string]func(http.ResponseWriter, *http.Request) {
	return map[string]func(http.ResponseWriter, *http.Request){
		"POST": func(w http.ResponseWriter, r *http.Request) {
			data := FetchRequestBody(r)
			if errors.Create(w, data) {
				return
			}
			crudData := crudObj.Create(parsers.Create(data)...)
			fmt.Fprintf(w, HandleCrudData(crudData, translators.Create))
		},
		"GET": func(w http.ResponseWriter, r *http.Request) {
			data := FetchRequestBody(r)
			if errors.Retrieve(w, data) {
				return
			}
			crudData := crudObj.Retrieve(parsers.Retrieve(data))
			fmt.Fprintf(w, HandleCrudData(crudData, translators.Retrieve))
		},
		"PATCH": func(w http.ResponseWriter, r *http.Request) {
			data := FetchRequestBody(r)
			if errors.Update(w, data) {
				return
			}
			key, newData := parsers.Update(data)
			crudData := crudObj.Update(key, newData)
			fmt.Fprintf(w, HandleCrudData(crudData, translators.Update))
		},
		"DELETE": func(w http.ResponseWriter, r *http.Request) {
			data := FetchRequestBody(r)
			if errors.Delete(w, data) {
				return
			}
			crudObj.Delete(parsers.Delete(data))
			fmt.Fprintf(w, string(data))
		},
	}
}
