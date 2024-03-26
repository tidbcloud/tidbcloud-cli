package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConnectInfo struct {
	Endpoint   []Endpoint   `json:"endpoint"`
	Os         []Os         `json:"os"`
	Ca         []Ca         `json:"ca"`
	Client     []Client     `json:"client"`
	Variable   []Variable   `json:"variable"`
	Connection []Connection `json:"connection"`
}
type Endpoint struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Os struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Ca struct {
	Os   string `json:"os"`
	Type string `json:"type"`
	Path string `json:"path"`
}
type Options struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Client struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Options []Options `json:"options,omitempty"`
}
type Variable struct {
	ID          string `json:"id"`
	Placeholder string `json:"placeholder"`
}
type Connection struct {
	Endpoint   string   `json:"endpoint"`
	Client     string   `json:"client"`
	Type       string   `json:"type"`
	Path       string   `json:"path"`
	DownloadCa []string `json:"download_ca,omitempty"`
	Doc        string   `json:"doc,omitempty"`
	Content    string   `json:"content,omitempty"`
}

const ConnectInfoStruct = `
type ConnectInfo struct {
	Endpoint   []Endpoint 
	Os         []Os         
	Ca         []Ca        
	Client     []Client     
	Variable   []Variable   
	Connection []Connection 
}
type Endpoint struct {
	ID   string 
	Name string 
}
func (e Endpoint) String() string {
	return e.Name
}
type Os struct {
	ID   string 
	Name string 
}
func (o Os) String() string {
	return o.Name
}
type Ca struct {
	Os   string
	Type string 
	Path string 
}
type Options struct {
	ID   string 
	Name string 
}
func (o Options) String() string {
	return o.Name
}
type Client struct {
	ID      string   
	Name    string    
	Options []Options 
}
func (c Client) String() string {
	return c.Name
}
type Variable struct {
	ID          string 
	Placeholder string 
}
type Connection struct {
	Endpoint   string  
	Client     string  
	Type       string  
	Path       string  
	DownloadCa []string 
	Doc        string  
	Content   string
}

`

const license = `// Copyright 2024 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

`

func main() {
	// read manifest file and unmarshal it to ConnectInfo struct
	connectInfo, err := ReadManifest()
	if err != nil {
		panic(err)
	}
	// generate code
	generateGoCode := GenerateCode(connectInfo)

	// write code to file
	err = os.WriteFile("../../../internal/util/connect-info.go", []byte(generateGoCode), 0644)
	if err != nil {
		panic(err)
	}
}

func ReadManifest() (ConnectInfo, error) {
	var connectInfo ConnectInfo
	jsonString, err := os.ReadFile("meta/manifest.json")
	if err != nil {
		return connectInfo, err
	}
	err = json.Unmarshal(jsonString, &connectInfo)
	if err != nil {
		return connectInfo, err
	}
	return connectInfo, nil
}

func GenerateCode(connectInfo ConnectInfo) string {
	GenerateGoCode := license
	GenerateGoCode += "// Code is generated. DO NOT EDIT.\n\n"
	GenerateGoCode += "package util\n"
	GenerateGoCode += ConnectInfoStruct

	endpoint := connectInfo.Endpoint
	connectInfoOs := connectInfo.Os
	ca := connectInfo.Ca
	client := connectInfo.Client
	variable := connectInfo.Variable
	connection := connectInfo.Connection

	// Endpoint
	GenerateGoCode += "var ConnectInfoEndpoint = []Endpoint{\n"
	for _, v := range endpoint {
		GenerateGoCode += "\t{ID: \"" + v.ID + "\", Name: \"" + v.Name + "\"},\n"
	}
	GenerateGoCode += "}\n"

	// Os
	GenerateGoCode += "var ConnectInfoOs = []Os{\n"
	for _, v := range connectInfoOs {
		GenerateGoCode += "\t{ID: \"" + v.ID + "\", Name: \"" + v.Name + "\"},\n"
	}
	GenerateGoCode += "}\n"

	// Ca
	GenerateGoCode += "var ConnectInfoCa = []Ca{\n"
	for _, v := range ca {
		GenerateGoCode += "\t{Os: \"" + v.Os + "\", Type: \"" + v.Type + "\", Path: \"" + v.Path + "\"},\n"
	}
	GenerateGoCode += "}\n"

	// Client
	GenerateGoCode += "var ConnectInfoClient = []Client{\n"
	for _, v := range client {
		GenerateGoCode += "\t{ID: \"" + v.ID + "\", Name: \"" + v.Name + "\", Options: []Options{\n"
		for _, o := range v.Options {
			GenerateGoCode += "\t\t{ID: \"" + o.ID + "\", Name: \"" + o.Name + "\"},\n"
		}
		GenerateGoCode += "\t}},\n"
	}
	GenerateGoCode += "}\n"
	// Variable
	GenerateGoCode += "var ConnectInfoVariable = []Variable{\n"
	for _, v := range variable {
		GenerateGoCode += "\t{ID: \"" + v.ID + "\", Placeholder: \"" + v.Placeholder + "\"},\n"
	}
	GenerateGoCode += "}\n"

	// Connection
	GenerateGoCode += "var ConnectInfoConnection = []Connection{\n"
	for i, v := range connection {
		path := v.Path[1:]
		connectionString, err := os.ReadFile("meta" + path)
		if err != nil {
			panic(err)
		}
		connection[i].Content = string(connectionString)
		GenerateGoCode += "\t{Endpoint: \"" + v.Endpoint + "\", Client: \"" + v.Client + "\", Type: \"" + v.Type + "\", Path: \"" + v.Path + "\", DownloadCa: []string{\n"
		for _, d := range v.DownloadCa {
			GenerateGoCode += "\t\t\"" + d + "\",\n"
		}
		GenerateGoCode += "\t}, Doc: \"" + v.Doc + "\", Content: " + fmt.Sprintf("`%s`},\n", v.Content)
	}
	GenerateGoCode += "}\n"
	return GenerateGoCode
}
