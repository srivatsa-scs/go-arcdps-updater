package main

import (
	"arcdps/config"
	counter "arcdps/helper"
	"arcdps/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {

	log := logger.Logger()

	log.Info().Msg("Reading configuration file...")
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		log.Info().Msg("Configuration file not found...")
		log.Error().Stack().Msg(err.Error())
		log.Fatal()
	}

	filepath := cfg.Destination
	yy, mm, dd := time.Now().Date()

	oldFilePath := fmt.Sprintf("%v.%v-%v-%v.old", filepath, yy, mm, dd)

	renameFlag := DoesFileExist(filepath)

	if renameFlag {
		log.Info().Msgf("Old d3d9.dll found, renaming it to %v", oldFilePath)
		err := os.Rename(filepath, oldFilePath)
		if err != nil {
			log.Info().Msg("Error renaming the file...not really sure how this can even happen...")
			log.Fatal().Msg(err.Error())
		}
	}

	err = DownloadFile(cfg.URL, cfg.Destination)
	if err != nil {
		log.Info().Msg("There was an error when downloading the file...")
		if DoesFileExist(filepath + ".old") {
			log.Info().Msg("Restoring old d3d9.dll file...")
			os.Rename(oldFilePath, filepath)
		}
		if DoesFileExist(filepath + ".tmp") {
			log.Info().Msg("Removing temp file...")
			os.Remove(filepath + ".tmp")
		}

		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("File downloaded successfully...")
	os.Rename(filepath+".tmp", filepath)
	log.Info().Msg("Removing .tmp from filename...")

}

func DoesFileExist(filepath string) bool {

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false
		}

	}
	return true
}

func DownloadFile(url string, filepath string) error {

	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()
	log.Info().Msg("File created with .tmp to file extension...")

	log.Info().Msg("Fetching the latest d3d9.dll....")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/* WIP MD5 check
	md5Resp, err := http.Get("https://www.deltaconnected.com/arcdps/x64/d3d9.dll.md5sum")
	h := md5.New()
	if _, err := io.Copy(h, out); err != nil {
		return err
	}
	*/

	counter := &counter.WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	fmt.Println()
	if err != nil {
		return err
	}
	return nil
}

/*
[Y] Read paths from config file (JSON)
[Y] Download File
[Y] Rename old file -> temp file
[Y] Move downloaded file -> bin64
[Y] If error delete .tmp, rename .old -> basename
[X] MD5 Checksum
[X] Loggerfile
*/
