package main

import (
	"arcdps/config"
	counter "arcdps/helper"
	"arcdps/logger"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
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

	logLvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(logLvl)
	}

	log.Info().Msgf("%v", log.Info().Enabled())

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
	md1 := GetMd5DigestOf(filepath)
	md2 := GetMd5FromUrl()
	log.Info().Msgf("Checking if Md5 is equal: %v", strings.Contains(md2, md1))
	if !strings.Contains(md2, md1) {
		if DoesItExist(oldFilePath, false) {
			log.Info().Msg("Restoring old d3d9.dll file...")
			os.Rename(oldFilePath, filepath)
		}
		if DoesItExist(filepath+".tmp", false) {
			log.Info().Msg("Removing temp file...")
			os.Remove(filepath + ".tmp")
		}
	} else {

		log.Info().Msg("File downloaded successfully...")
		os.Rename(filepath+".tmp", filepath)
		log.Info().Msg("Removing .tmp from filename...")
	}
	/* -- Gw2Launcher Specific -- */
	if cfg.EnableGw2Launcher {

		if DoesItExist(cfg.Gw2LauncherPath, true) {
			returnMap := replaceAllFiles(numberOfFolders(cfg.Gw2LauncherPath), filepath, cfg.Gw2LauncherPath, cfg.Filename)
			if returnMap == nil {
				log.Error().Msg("Error Occoured when replacing files")
			}
		}
	} else {
		log.Warn().Msg("GW2 Launcher Config Disabled or Not Set")
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

	counter := &counter.WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	fmt.Println()

	if err != nil {
		return err
	}
	return nil
}

func GetMd5FromUrl() string {
	log := logger.Logger()
	md5Resp, err := http.Get("https://www.deltaconnected.com/arcdps/x64/d3d9.dll.md5sum")
	if err != nil {
		log.Error().Err(err)
	}
	md5bytes, _ := io.ReadAll(md5Resp.Body)
	defer md5Resp.Body.Close()
	log.Debug().Msg(string(md5bytes))
	return string(md5bytes)
}

func GetMd5DigestOf(filepath string) string {
	log := logger.Logger()
	f, err := os.Open(filepath + ".tmp")
	if err != nil {
		log.Error().Err(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error().Err(err)
	}
	log.Debug().Msgf("%x", h.Sum(nil))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func numberOfFolders(folderPath string) int {
	log := logger.Logger()
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Warn().Msg("[GW2 Launcher]: path not set")
		log.Warn().Msg("Skipping GW2 Launcher Routines")
		return 0
	}

	numberOfFiles := 0

	for i := range files {
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
