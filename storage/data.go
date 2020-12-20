package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/Pauloo27/tuner/search"
	"github.com/Pauloo27/tuner/utils"
)

var dataFile string

type Playlist struct {
	Name  string
	Songs []*search.SearchResult
}

type TunerData struct {
	Version                                        string
	Playlists                                      []*Playlist
	Cache, FetchAlbum, LoadMPRIS, SearchSoundCloud bool
}

func CreateDataFolder(dataFolder string) {
	err := os.Mkdir(dataFolder, 0744)
	utils.HandleError(err, "Cannot create data folder at "+dataFolder)
}

func Load() *TunerData {
	dataFolder := utils.LoadDataFolder()
	dataFile = path.Join(dataFolder, "data.json")

	_, err := os.Stat(dataFile)

	if os.IsNotExist(err) {
		_, err = os.Create(dataFile)
		utils.HandleError(err, "Cannot create data file at "+dataFile)

		data := &TunerData{}
		Save(data)
		return data
	}

	file, err := os.OpenFile(dataFile, os.O_CREATE, 0644)
	utils.HandleError(err, "Cannot open data file at "+dataFile)

	defer file.Close()

	var data *TunerData

	buffer, err := ioutil.ReadAll(file)
	utils.HandleError(err, "Cannot read data file at "+dataFile)

	err = json.Unmarshal(buffer, &data)
	utils.HandleError(err, "Cannot unmarshal data file at "+dataFile)

	return data
}

func Save(data *TunerData) {
	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_TRUNC, 0644)
	utils.HandleError(err, "Cannot open data file at "+dataFile)

	file.Truncate(0)

	defer file.Close()

	buffer, err := json.Marshal(*data)
	utils.HandleError(err, "Cannot marshal data")

	file.Write(buffer)
}
