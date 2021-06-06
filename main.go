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

	filepath := fmt.Sprintf("%v%v", cfg.Destination, cfg.Filename)
	yy, mm, dd := time.Now().Date()

	log.Info().Msg(filepath)

	oldFilePath := fmt.Sprintf("%v.%v-%v-%v.old", filepath, yy, mm, dd)

	renameFlag := DoesItExist(filepath, false)
	log.Info().Msgf("%v", renameFlag)
	if renameFlag {
		log.Info().Msgf("Old d3d9.dll found, renaming it to %v", oldFilePath)
		err := os.Rename(filepath, oldFilePath)
		if err != nil {
			log.Info().Msg("Error renaming the file...not really sure how this can even happen...")
			log.Fatal().Msg(err.Error())
		}
	}

	err = DownloadFile(cfg.URL, filepath)
	if err != nil {
		log.Warn().Msg("There was an error when downloading the file...")
		if DoesItExist(filepath+".old", false) {
			log.Info().Msg("Restoring old d3d9.dll file...")
			os.Rename(oldFilePath, filepath)
		}
		if DoesItExist(filepath+".tmp", false) {
			log.Info().Msg("Removing temp file...")
			os.Remove(filepath + ".tmp")
		}

		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("File downloaded successfully...")
	os.Rename(filepath+".tmp", filepath)
	log.Info().Msg("Removing .tmp from filename...")
	/* -- Gw2Launcher Specific -- */

	if DoesItExist(cfg.Gw2LauncherPath, true) {
		returnMap := replaceAllFiles(numberOfFolders(cfg.Gw2LauncherPath), filepath, cfg.Gw2LauncherPath, cfg.Filename)
		if returnMap == nil {
			log.Error().Msg("Error Occoured when replacing files")
		}
	}

}

func DoesItExist(filepath string, isDir bool) bool {

	if stat, err := os.Stat(filepath); err != nil {
		if isDir {
			if stat.IsDir() && os.IsNotExist(err) {
				return false
			}

		} else if os.IsNotExist(err) {
			return false
		}

	}
	return true
}

func DownloadFile(url string, filepath string) error {

	log := logger.Logger()
	log.Info().Msg(filepath + ".tmp")

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
[Y] Loggerfile
*/

func numberOfFolders(folderPath string) int {
	log := logger.Logger()
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Warn().Msg("[GW2 Launcher]: path not set")
		log.Warn().Msg("Skipping GW2 Launcher Routines")
		return 0
	}

	numberOfFiles := 0

	for i, _ := range files {
		numberOfFiles += i
	}
	log.Info().Msgf("Number of folders: %v", numberOfFiles)
	return numberOfFiles
}

func replaceAllFiles(n int, filepath string, gw2LauncherPath string, filename string) *map[int]bool {
	log := logger.Logger()
	if n < 1 {
		return nil
	}

	returnMap := make(map[int]bool)
	reader, err := os.Open(filepath)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	var replaceLogMsg []string

	for i := 1; i <= n; i++ {
		returnMap[i] = CopyFile(i, filepath, gw2LauncherPath, filename)
		// log.Info().Msgf("%v", returnMap[i])
		replaceLogMsg = append(replaceLogMsg, fmt.Sprintf("%v: %v", i, returnMap[i]))
	}
	err = reader.Close()

	log.Info().Msgf("%v", replaceLogMsg)
	if err != nil {
		fmt.Println(err)
	}

	return &returnMap
}

func CopyFile(i int, srcpath string, gw2LauncherPath string, filename string) bool {
	log := logger.Logger()

	src, err := os.Open(srcpath)
	if err != nil {
		fmt.Println(err)
	}

	destpath := fmt.Sprintf("%v/%v/bin64/%v", gw2LauncherPath, i, filename)
	log.Info().Msgf("Copying file from %v to %v", srcpath, destpath)
	dest, err := os.Create(destpath)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return false
	}
	defer dest.Close()
	defer src.Close()
	_, err = io.Copy(dest, src)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return false
	}
	return true
}
