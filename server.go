package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter" //third-party Router
	"gopkg.in/yaml.v2"                    //to handle .yaml files
)

type Server []SServer

type SServer struct {
	SERVER            string `yaml:"SERVER"`
	SERVERDISPLAYNAME string `yaml:"SERVER_DISPLAY_NAME"`
	SERVERPORT        string `yaml:"SERVER_PORT"`
	SERVERPARAMS      []struct {
		PARAMETERNAME  string `yaml:"PARAMETER_NAME"`
		ISDEFAULT      bool   `yaml:"IS_DEFAULT"`
		PARAMETERTYPE  string `yaml:"PARAMETER_TYPE"`
		PARAMETERVALUE string `yaml:"PARAMETER_VALUE"`
	} `yaml:"SERVER_PARAMS"`
}

type Parameters struct {
	FirstParameter  FirstParameter  `json:"First Parameter"`
	SecondParameter SecondParameter `json:"Second Parameter"`
	ThirdParameter  ThirdParameter  `json:"Third Parameter"`
	FourthParameter FourthParameter `json:"Fourth Parameter"`
}
type FirstParameter struct {
	Level string `json:"Level"`
	Rank  int    `json:"Rank"`
}
type SecondParameter struct {
	Level string `json:"Level"`
	Rank  int    `json:"Rank"`
}
type ThirdParameter struct {
	Level string `json:"Level"`
	Rank  int    `json:"Rank"`
}
type FourthParameter struct {
	Level string `json:"Level"`
	Rank  int    `json:"Rank"`
}

var tpl *template.Template
var servers Server
var parameters Parameters

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	mux := httprouter.New()
	mux.ServeFiles("/assets/*filepath", http.Dir("static"))

	mux.GET("/", index)
	mux.GET("/servers", server)
	mux.GET("/edit/:server", edit)
	mux.POST("/save", save) //only yaml modification is implemented
	http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func server(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	servers, parameters = readAndunMarshal()
	queryParams := req.URL.Query()
	server := queryParams.Get("server")

	var err error
	if server == "first" {
		err = tpl.ExecuteTemplate(w, "index.gohtml", servers.makeMap(0, parameters))

	} else {
		err = tpl.ExecuteTemplate(w, "index.gohtml", servers.makeMap(1, parameters))
	}
	if err != nil {
		log.Fatalln(err)
	}

}

func edit(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	servers, parameters = readAndunMarshal()
	server := ps.ByName("server")

	var err error
	if server == "first" {
		err = tpl.ExecuteTemplate(w, "editor.gohtml", servers.makeMap(0, parameters))

	} else {
		err = tpl.ExecuteTemplate(w, "editor.gohtml", servers.makeMap(1, parameters))
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func save(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	req.ParseForm()
	P1, _ := strconv.ParseBool(req.Form["ISDEFAULT"][0])
	P2, _ := strconv.ParseBool(req.Form["ISDEFAULT"][1])
	Ydata := SServer{
		SERVER:            req.Form["SERVER"][0],
		SERVERDISPLAYNAME: req.Form["SERVER_DISPLAY_NAME"][0],
		SERVERPORT:        req.Form["SERVER_PORT"][0],
		SERVERPARAMS: []struct {
			PARAMETERNAME  string `yaml:"PARAMETER_NAME"`
			ISDEFAULT      bool   `yaml:"IS_DEFAULT"`
			PARAMETERTYPE  string `yaml:"PARAMETER_TYPE"`
			PARAMETERVALUE string `yaml:"PARAMETER_VALUE"`
		}{
			{req.Form["PARAMETERNAME"][0], P1, req.Form["PARAMETERTYPE"][0], req.Form["PARAMETERVALUE"][0]},
			{req.Form["PARAMETERNAME"][1], P2, req.Form["PARAMETERTYPE"][1], req.Form["PARAMETERVALUE"][1]},
		},
	}

	servers, _ = readAndunMarshal()

	if Ydata.SERVER == "First" {
		servers[0] = Ydata
	} else {
		servers[1] = Ydata
	}

	data, err := yaml.Marshal(&servers)
	if err != nil {
		log.Fatalln(err)
	}
	_ = ioutil.WriteFile("test.yaml", data, 0)

	http.Redirect(w, req, "/servers?server="+strings.ToLower(Ydata.SERVER), http.StatusMovedPermanently)

}

func readAndunMarshal() (Server, Parameters) { //reads and returns json and yaml data structured
	yfile, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	jfile, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(yfile, &servers)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(jfile, &parameters)
	if err != nil {
		log.Fatalln(err)
	}

	return servers, parameters
}

func (s Server) makeMap(i int, p Parameters) map[string]interface{} { //makes struct into map for using in template
	values := map[string]interface{}{}
	values["SERVER"] = servers[i].SERVER
	values["SERVERDISPLAYNAME"] = servers[i].SERVERDISPLAYNAME
	values["SERVERPORT"] = servers[i].SERVERPORT
	values["SERVERPARAMS"] = []map[string]interface{}{
		{
			"PARAMETERNAME":  servers[i].SERVERPARAMS[0].PARAMETERNAME,
			"ISDEFAULT":      servers[i].SERVERPARAMS[0].ISDEFAULT,
			"PARAMETERTYPE":  servers[i].SERVERPARAMS[0].PARAMETERTYPE,
			"PARAMETERVALUE": servers[i].SERVERPARAMS[0].PARAMETERVALUE,
		}, {
			"PARAMETERNAME":  servers[i].SERVERPARAMS[1].PARAMETERNAME,
			"ISDEFAULT":      servers[i].SERVERPARAMS[1].ISDEFAULT,
			"PARAMETERTYPE":  servers[i].SERVERPARAMS[1].PARAMETERTYPE,
			"PARAMETERVALUE": servers[i].SERVERPARAMS[1].PARAMETERVALUE,
		},
	}
	if i == 0 {
		values["FirstParameter"] = map[string]interface{}{"Level": p.FirstParameter.Level, "Rank": p.FirstParameter.Rank}
		values["SecondParameter"] = map[string]interface{}{"Level": p.SecondParameter.Level, "Rank": p.SecondParameter.Rank}
	} else {
		values["ThirdParameter"] = map[string]interface{}{"Level": p.ThirdParameter.Level, "Rank": p.ThirdParameter.Rank}
		values["FourthParameter"] = map[string]interface{}{"Level": p.FourthParameter.Level, "Rank": p.FourthParameter.Rank}
	}
	return values
}
