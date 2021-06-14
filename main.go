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
		log.Error().Stack().Err(err).Msg("")
		log.Fatal()
	}

	logLvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Warn().Err(err).Msg("")
	} else {
		zerolog.SetGlobalLevel(logLvl)
		log.Info().Msgf("LogLevel: %v", logLvl)
	}

	log.Info().Msgf("Is Debug Mode Enabled? %v", log.Debug().Enabled())

	log.Debug().Msgf("Configuration: %v", cfg)

	filepath := fmt.Sprintf("%v%v", cfg.Destination, cfg.Filename)
	yy, mm, dd := time.Now().Date()

	log.Debug().Msg(filepath)

	oldFilePath := fmt.Sprintf("%v.%v-%v-%v.old", filepath, yy, mm, dd)

	log.Debug().Msg(oldFilePath)

	renameFlag := DoesItExist(filepath, false)
	log.Debug().Msgf("Rename Flag: %v", renameFlag)
	if renameFlag {
		log.Info().Msgf("Old d3d9.dll found, renaming it to %v", oldFilePath)
		err := os.Rename(filepath, oldFilePath)
		if err != nil {
			log.Info().Msg("Error renaming the file...not really sure how this can even happen...")
			log.Fatal().Stack().Err(err).Msg("")
		}
	}

	err = DownloadFile(cfg.URL, filepath)
	if err != nil {
		log.Error().Msg("There was an error when downloading the file...")
		RestoreOldVersionAndRemoveTemp(oldFilePath, filepath)
		log.Fatal().Stack().Err(err).Msg("")
	}
	fileDigest := GetMd5DigestOfFile(filepath)
	urlDigest := GetMd5FromUrl()
	log.Info().Msgf("Checking if MD5 Digest Matches: %v", strings.Contains(urlDigest, fileDigest))
	if !strings.Contains(urlDigest, fileDigest) {
		RestoreOldVersionAndRemoveTemp(oldFilePath, filepath)
	} else {
		log.Info().Msg("File downloaded successfully...")
		os.Rename(filepath+".tmp", filepath)
		log.Debug().Msg("Removing .tmp from filename...")
	}
	/* -- Gw2Launcher Specific -- */
	if cfg.EnableGw2Launcher {
		log.Debug().Msg("GW2 Launcher Config Enabled...")
		if DoesItExist(cfg.Gw2LauncherPath, true) {
			log.Debug().Msg("GW2 Launcher Folder exists...")
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

func RestoreOldVersionAndRemoveTemp(oldFilePath string, filepath string) {
	log := logger.Logger()
	if DoesItExist(oldFilePath, false) {
		log.Info().Msg("Restoring old d3d9.dll file...")
		os.Rename(oldFilePath, filepath)
	}
	if DoesItExist(filepath+".tmp", false) {
		log.Info().Msg("Removing temp file...")
		os.Remove(filepath + ".tmp")
	}

}

func DownloadFile(url string, filepath string) error {

	// url = "https://www.deltaconnected.com/arcdps/x64/d3d9.dll"

	log := logger.Logger()
	log.Debug().Msg(filepath + ".tmp")

	out, err := os.Create(filepath + ".tmp")

	if err != nil {
		return err
	}
	defer out.Close()
	log.Debug().Msg("File created with .tmp to file extension...")

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
		log.Error().Stack().Err(err).Msg("")
	}
	md5bytes, _ := io.ReadAll(md5Resp.Body)
	defer md5Resp.Body.Close()
	log.Debug().Msg(string(md5bytes))
	return string(md5bytes)
}

func GetMd5DigestOfFile(filepath string) string {
	log := logger.Logger()
	f, err := os.Open(filepath + ".tmp")
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error().Stack().Err(err).Msg("")
	}
	log.Debug().Msgf("%x", h.Sum(nil))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func numberOfFolders(folderPath string) int {
	log := logger.Logger()
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Warn().Msg("[GW2 Launcher]: path not set")
		log.Warn().Msg("[GW2 Launcher]: Skipping GW2 Launcher Routines")
		return 0
	}
	numberOfFolders := len(files)
	log.Info().Msgf("Number of folders: %v", numberOfFolders)
	return numberOfFolders
}

func replaceAllFiles(n int, filepath string, gw2LauncherPath string, filename string) *map[int]bool {
	log := logger.Logger()
	log.Debug().Msgf("[GW2 Launcher] replaceAllFiles(%v,%v,%v,%v)", n, filepath, gw2LauncherPath, filename)
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

	log.Debug().Msgf("[GW2 Launcher] The following replacements occoured: %v", replaceLogMsg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	return &returnMap
}

func CopyFile(i int, srcpath string, gw2LauncherPath string, filename string) bool {
	log := logger.Logger()

	src, err := os.Open(srcpath)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	destpath := fmt.Sprintf("%v%v"+string(os.PathSeparator)+"bin64"+string(os.PathSeparator)+"%v", gw2LauncherPath, i, filename)
	log.Debug().Msgf("Copying file from %v to %v", srcpath, destpath)
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
