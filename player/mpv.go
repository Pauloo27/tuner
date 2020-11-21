package player

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/tuner/lyric"
	"github.com/Pauloo27/tuner/search"
	"github.com/Pauloo27/tuner/utils"
	"github.com/godbus/dbus/v5"
)

type UpdateHandler func(result *search.YouTubeResult, mpv *MPV)
type SaveFunction func(result *search.YouTubeResult, mpv *MPV)

type MPV struct {
	Pid                                  int
	Cmd                                  *exec.Cmd
	Player                               *mpris.Player
	ShowHelp, ShowLyric, ShowURL, Saving bool
	LyricIndex                           int
	LyricLines                           []string
	Result                               *search.YouTubeResult
	onUpdate                             UpdateHandler
	save                                 SaveFunction
	Exitted                              bool
}

func ConnectToMPV(cmd *exec.Cmd, result *search.YouTubeResult, onUpdate UpdateHandler, save SaveFunction) *MPV {
	for {
		if cmd.Process != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	pid := cmd.Process.Pid

	nameWithPID := fmt.Sprintf("org.mpris.MediaPlayer2.mpv.instance%d", pid)

	conn, err := dbus.SessionBus()
	utils.HandleError(err, "Cannot connect to dbus")

	names, err := mpris.List(conn)
	utils.HandleError(err, "Cannot list players")

	playerName := ""

	for _, name := range names {
		if name == nameWithPID {
			playerName = name
			break
		}
	}

	if playerName == "" {
		playerName = "org.mpris.MediaPlayer2.mpv"
	}

	player := mpris.New(conn, playerName)
	utils.HandleError(err, "Cannot connect to mpv")

	mpv := MPV{pid, cmd, player, false, false, false, false, 0, []string{}, result, onUpdate, save, false}

	mpv.Update()

	go func() {
		ch := make(chan *dbus.Signal)
		err := mpv.Player.OnSignal(ch)
		utils.HandleError(err, "Cannot add signal handler")

		for range ch {
			if mpv.Exitted {
				break
			}
			mpv.Update()
		}
	}()

	return &mpv
}

func (i *MPV) Update() {
	i.onUpdate(i.Result, i)
}

func (i *MPV) Save() {
	i.save(i.Result, i)
}

func (i *MPV) PlayPause() {
	_ = i.Player.PlayPause()
	i.Update()
}

func (i *MPV) Exit() {
	i.Exitted = true
}

func (i *MPV) FetchLyric() {
	path, err := lyric.SearchFor(fmt.Sprintf("%s %s", i.Result.Title, i.Result.Uploader))
	if err != nil {
		i.LyricLines = []string{"Cannot get lyric"}
		i.Update()
		return
	}

	l, err := lyric.Fetch(path)
	if err != nil {
		i.LyricLines = []string{"Cannot get lyric"}
		i.Update()
		return
	}

	i.LyricLines = strings.Split(l, "\n")
	i.Update()
}
