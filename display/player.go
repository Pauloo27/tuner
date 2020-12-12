package display

import (
	"fmt"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/tuner/keybind"
	"github.com/Pauloo27/tuner/player"
	"github.com/Pauloo27/tuner/search"
	"github.com/Pauloo27/tuner/state"
	"github.com/Pauloo27/tuner/utils"
)

const (
	pausedIcon  = ""
	playingIcon = ""
)

func ShowPlaying(result *search.YouTubeResult, mpv *player.MPV) {
	if !state.Playing {
		return
	}

	utils.ClearScreen()

	icon := playingIcon

	playback, _ := mpv.Player.GetPlaybackStatus()
	if playback != mpris.PlaybackPlaying {
		icon = pausedIcon
	}

	extra := ""
	if status, err := mpv.Player.GetLoopStatus(); err == nil {
		if status == mpris.LoopTrack {
			extra += utils.ColorWhite + "  "
		} else if status == mpris.LoopPlaylist {
			extra += utils.ColorBlue + "  "
		}
	}

	if mpv.IsPlaylist() {
		fmt.Printf("Playing: %s (%d/%d)\n",
			mpv.Playlist.Name,
			mpv.PlaylistIndex+1,
			len(mpv.Playlist.Songs),
		)
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

	if mpv.ShowURL {
		fmt.Printf("%s%s%s\n", utils.ColorBlue, result.URL(), utils.ColorReset)
	}

	if mpv.ShowHelp {
		fmt.Println("\n" + utils.ColorBlue + "Keybinds:")
		for _, bind := range keybind.ListBinds() {
			fmt.Printf("  %s: %s\n", bind.KeyName, bind.Description)
		}
	}

	if mpv.ShowLyric {
		fmt.Println(utils.ColorBlue)
		lines := len(mpv.LyricLines)
		if lines == 0 {
			fmt.Println("Fetching lyric...")
		}
		for i := mpv.LyricIndex; i < mpv.LyricIndex+15; i++ {
			if i == lines {
				break
			}
			fmt.Println(mpv.LyricLines[i])
		}
	}

	fmt.Print(utils.ColorReset)

}