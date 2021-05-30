package main

import (
	"arcdps/config"
	counter "arcdps/helper"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	log.Default()

	log.Println("Reading configuration file...")
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		log.Println("Configuration file not found...")
		log.Fatal(err)
	}

	filepath := cfg.Destination

	renameFlag := DoesFileExist(filepath)

	if renameFlag {
		log.Println("Appending .old to filename...")
		err := os.Rename(filepath, filepath+".old")
		if err != nil {
			log.Println("Error appending .old to filename... not really sure how this can even happen...")
			log.Fatal(err)
		}
	}

	err = DownloadFile(cfg.URL, cfg.Destination)
	if err != nil {
		log.Println("There was an error when downloading the file...")
		if DoesFileExist(filepath + ".old") {
			log.Println("Restoring old d3d9.dll file...")
			os.Rename(filepath+".old", filepath)
		}
		if DoesFileExist(filepath + ".tmp") {
			log.Println("Removing temp file...")
			os.Remove(filepath + ".tmp")
		}

		log.Fatal(err)
	}
	log.Println("File downloaded successfully...")
	os.Rename(filepath+".tmp", filepath)
	log.Println("Removing .tmp from filename...")

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
	log.Println("File created with .tmp to file extension...")

	log.Println("Fetching the latest d3d9.dll....")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	counter := &counter.WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
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
*/
