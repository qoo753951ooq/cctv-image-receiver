package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

var CurrentConfig Configuration

type Configuration struct {
	Param Param   `json:"param"`
	Os    Os      `json:"os"`
	Cctvs []*Cctv `json:"cctv"`
}

type Param struct {
	ProjectName  string `json:"project_name"`
	CronSpec     string `json:"cron_spec"`
	ImageOutPath string `json:"image_out_path"`
}

type Os struct {
	IsWindows    bool   `json:"is_windows"`
	TimeFilePath string `json:"time_file_path"`
}

type Cctv struct {
	User           string `json:"user"`
	Pwd            string `json:"pwd"`
	No             string `json:"no"`
	Name           string `json:"name"`
	Ip             string `json:"ip"`
	ImageUrlSuffix string `json:"image_url_suffix"`
	ImageWidth     uint   `json:"image_width"`
	ImageHeight    uint   `json:"image_height"`
}

func init() {
	file, err := os.Open("conf/config.json")

	if err != nil {
		file, _ = os.Open("./conf/config.json")
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&CurrentConfig)

	if err != nil {
		fmt.Println("Decode config error : ", err)
	}
}
