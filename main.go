package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/tuner/command"
	"github.com/Pauloo27/tuner/commands"
	"github.com/Pauloo27/tuner/controller"
	"github.com/Pauloo27/tuner/options"
	"github.com/Pauloo27/tuner/search"
	"github.com/Pauloo27/tuner/utils"
	"github.com/eiannone/keyboard"
)

var close = make(chan os.Signal)
var playing = false
var mpvInstance controller.MPV

const (
	pausedIcon  = ""
	playingIcon = ""
)

func doSearch(searchTerm string, limit int) (results []search.YouTubeResult, err error) {
	c := make(chan bool)

	go utils.PrintWithLoadIcon(fmt.Sprintf("Searching for %s", searchTerm), c, 100*time.Millisecond, true)
	results = search.SearchYouTube(searchTerm, limit)

	c <- true
	<-c
	return
}

func setupCloseHandler() {
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-close
			if !playing {
				utils.MoveCursorTo(1, 1)
				utils.ClearScreen()
				fmt.Println("Bye!")
				os.Exit(0)
			}
		}
	}()
}

func listResults(results []search.YouTubeResult) {
	utils.ClearScreen()
	for i, result := range results {
		bold := ""
		if i%2 == 0 {
			bold = utils.ColorBold
		}

		defaultColor := bold + utils.ColorWhite
		altColor := bold + utils.ColorGreen

		duration := result.Duration

		if duration == "" {
			duration = utils.ColorRed + "LIVE"
		}

		fmt.Printf("  %s%d: %s %sfrom %s - %s%s\n",
			defaultColor, i+1,
			altColor+result.Title,
			defaultColor,
			altColor+result.Uploader,
			defaultColor+duration,
			utils.ColorReset,
		)
	}
}

func listenToKeyboard(cmd *exec.Cmd, playerCtl *controller.MPV) {
	err := keyboard.Open()
	utils.HandleError(err, "Cannot open keyboard")
	for {
		c, key, err := keyboard.GetKey()
		if err != nil {
			break
		}

		if bind, ok := byKey[key]; ok {
			bind.Handler(cmd, playerCtl)
		} else if bind, ok := byChar[c]; ok {
			bind.Handler(cmd, playerCtl)
		}
	}
}

func displayPlayingScreen(result search.YouTubeResult, mpv *controller.MPV) {
	for {
		if !playing {
			break
		}
		utils.ClearScreen()

		icon := playingIcon

		playback, _ := mpv.Player.GetPlaybackStatus()
		if playback != mpris.PlaybackPlaying {
			icon = pausedIcon
		}

		extra := utils.ColorWhite
		if status, err := mpv.Player.GetLoopStatus(); err == nil && status == mpris.LoopTrack {
			extra += "  "
		}

		fmt.Printf(" %s  %s %sfrom %s%s%s\n",
			icon,
			utils.ColorGreen+result.Title,
			utils.ColorWhite,
			utils.ColorGreen+result.Uploader,
			extra,
			utils.ColorReset,
		)

		if status, _ := mpv.Player.GetPlaybackStatus(); status != "" {
			volume, _ := mpv.Player.GetVolume()
			fmt.Printf("Volume: %s%.0f%%%s\n", utils.ColorGreen, volume*100, utils.ColorReset)
		}

		if mpv.ShowHelp {
			fmt.Println("\n" + utils.ColorBlue + "Keybinds:")
			for _, bind := range listBinds() {
				fmt.Printf("  %s: %s\n", bind.KeyName, bind.Description)
			}
			fmt.Print(utils.ColorReset)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func play(result search.YouTubeResult) {
	url := fmt.Sprintf("https://youtube.com/watch?v=%s", result.ID)

	parameters := []string{url}
	if !options.Options.ShowVideo {
		parameters = append(parameters, "--no-video", "--ytdl-format=worst")
	}
	if !options.Options.Cache {
		parameters = append(parameters, "--cache=no")
	}

	playing = true
	cmd := exec.Command("mpv", parameters...)

	go func() {
		mpvInstance = controller.ConnectToMPV(cmd)
		go displayPlayingScreen(result, &mpvInstance)
		go listenToKeyboard(cmd, &mpvInstance)
	}()

	err := cmd.Run()

	if err != nil && err.Error() != "exit status 4" && err.Error() != "signal: killed" {
		utils.HandleError(err, "Cannot run MPV")
	}

	keyboard.Close()
	playing = false
}

func tuneIn(warning *string) {
	utils.ClearScreen()

	fmt.Printf("%sUse /help to see the commands%s\n\n", utils.ColorBlue, utils.ColorReset)
	if *warning != "" {
		fmt.Printf("%s%s%s\n", utils.ColorYellow, *warning, utils.ColorReset)
		*warning = ""
	}

	rawSearchTerm, err := utils.AskFor("Search term (add ! prefix to play the first, Ctrl+C to exit)")
	utils.HandleError(err, "Cannot read user input")

	if strings.HasPrefix(rawSearchTerm, "/") {
		found, out := command.InvokeCommand(strings.TrimPrefix(rawSearchTerm, "/"))
		if !found {
			*warning = "Invalid command"
		} else {
			*warning = out
		}
		return
	}

	limit := 10
	if strings.HasPrefix(rawSearchTerm, "!") {
		limit = 1
	}
	searchTerm := strings.TrimPrefix(rawSearchTerm, "!")

	results, err := doSearch(searchTerm, limit)
	utils.HandleError(err, "Cannot search")

	if len(results) == 0 {
		*warning = "No results found"
		return
	}

	var entry search.YouTubeResult

	if len(results) == 1 {
		entry = results[0]
	} else {
		listResults(results)
		enteredIndex, err := utils.AskFor("Insert index of the video")
		utils.HandleError(err, "Cannot read user input")

		index, err := strconv.ParseInt(enteredIndex, 10, 64)

		if err != nil || index > int64(len(results)) || index <= 0 {
			*warning = "Invalid index"
			return
		}
		index--

		entry = results[index]
	}

	play(entry)
}

func main() {
	commands.SetupDefaultCommands()
	setupCloseHandler()
	registerDefaultKeybinds()
	warning := ""
	for {
		tuneIn(&warning)
	}
}
