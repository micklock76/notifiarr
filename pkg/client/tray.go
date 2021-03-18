// +build darwin windows

package client

import (
	"os"

	"github.com/Go-Lift-TV/discordnotifier-client/pkg/bindata"
	"github.com/Go-Lift-TV/discordnotifier-client/pkg/notifiarr"
	"github.com/Go-Lift-TV/discordnotifier-client/pkg/ui"
	"github.com/getlantern/systray"
	"golift.io/version"
)

/* This file handles the OS GUI elements. */

// startTray Run()s readyTray to bring up the web server and the GUI app.
func (c *Client) startTray() {
	systray.Run(c.readyTray, c.exitTray)
}

func (c *Client) exitTray() {
	c.sigkil = nil

	if err := c.Exit(); err != nil {
		c.Errorf("Shutting down web server: %v", err)
		os.Exit(1) // web server problem
	}
	// because systray wants to control the exit code? no..
	os.Exit(0)
}

// readyTray creates the system tray/menu bar app items, and starts the web server.
func (c *Client) readyTray() {
	b, err := bindata.Asset(ui.SystrayIcon)
	if err == nil {
		systray.SetTemplateIcon(b, b)
	} else {
		c.Errorf("Reading Icon: %v", err)
		systray.SetTitle("DNC")
	}

	systray.SetTooltip(c.Flags.Name() + " v" + version.Version)

	c.makeChannels() // make these before starting the web server.
	c.menu["info"].Disable()
	c.menu["dninfo"].Hide()
	c.menu["alert"].Hide() // currently unused.

	go c.watchKillerChannels()
	c.StartWebServer()
	c.watchGuiChannels()
}

//nolint:lll
func (c *Client) makeChannels() {
	c.menu["stat"] = ui.WrapMenu(systray.AddMenuItem("Running", "web server state unknown"))

	conf := systray.AddMenuItem("Config", "show configuration")
	c.menu["conf"] = ui.WrapMenu(conf)
	c.menu["view"] = ui.WrapMenu(conf.AddSubMenuItem("View", "show configuration"))
	c.menu["edit"] = ui.WrapMenu(conf.AddSubMenuItem("Edit", "edit configuration"))
	c.menu["key"] = ui.WrapMenu(conf.AddSubMenuItem("API Key", "set API Key"))
	c.menu["load"] = ui.WrapMenu(conf.AddSubMenuItem("Reload", "reload configuration"))

	link := systray.AddMenuItem("Links", "external resources")
	c.menu["link"] = ui.WrapMenu(link)
	c.menu["info"] = ui.WrapMenu(link.AddSubMenuItem(c.Flags.Name(), version.Print(c.Flags.Name())))
	c.menu["hp"] = ui.WrapMenu(link.AddSubMenuItem("DiscordNotifier.com", "open DiscordNotifier.com"))
	c.menu["wiki"] = ui.WrapMenu(link.AddSubMenuItem("DiscordNotifier Wiki", "open DiscordNotifier wiki"))
	c.menu["disc1"] = ui.WrapMenu(link.AddSubMenuItem("DiscordNotifier Discord", "open DiscordNotifier discord server"))
	c.menu["disc2"] = ui.WrapMenu(link.AddSubMenuItem("Go Lift Discord", "open Go Lift discord server"))
	c.menu["gh"] = ui.WrapMenu(link.AddSubMenuItem("GitHub Project", c.Flags.Name()+" on GitHub"))

	logs := systray.AddMenuItem("Logs", "log file info")
	c.menu["logs"] = ui.WrapMenu(logs)
	c.menu["logs_view"] = ui.WrapMenu(logs.AddSubMenuItem("View", "view the application log"))
	c.menu["logs_http"] = ui.WrapMenu(logs.AddSubMenuItem("HTTP", "view the HTTP log"))
	c.menu["logs_rotate"] = ui.WrapMenu(logs.AddSubMenuItem("Rotate", "rotate both log files"))

	data := systray.AddMenuItem("Notifiarr", "plex sessions, system snapshots, network monitors")
	c.menu["data"] = ui.WrapMenu(data)
	c.menu["snap_log"] = ui.WrapMenu(data.AddSubMenuItem("Log Full Snapshot", "write snapshot data to log file"))
	c.menu["plex_test"] = ui.WrapMenu(data.AddSubMenuItem("Test Plex Sessions", "send plex sessions to notifiarr test endpoint"))
	c.menu["snap_test"] = ui.WrapMenu(data.AddSubMenuItem("Test System Snapshot", "send system snapshot to notifiarr test endpoint"))
	c.menu["netw_test"] = ui.WrapMenu(data.AddSubMenuItem("Test Network Snapshot", "send network snapshot to notifiarr test endpoint"))
	c.menu["plex_dev"] = ui.WrapMenu(data.AddSubMenuItem("Dev Plex Sessions", "send plex sessions to notifiarr dev endpoint"))
	c.menu["snap_dev"] = ui.WrapMenu(data.AddSubMenuItem("Dev System Snapshot", "send system snapshot to notifiarr dev endpoint"))
	c.menu["netw_dev"] = ui.WrapMenu(data.AddSubMenuItem("Dev Network Snapshot", "send network snapshot to notifiarr dev endpoint"))
	c.menu["plex_prod"] = ui.WrapMenu(data.AddSubMenuItem("Prod Plex Sessions", "send plex sessions to notifiarr"))
	c.menu["snap_prod"] = ui.WrapMenu(data.AddSubMenuItem("Prod System Snapshot", "send system snapshot to notifiarr"))
	c.menu["netw_prod"] = ui.WrapMenu(data.AddSubMenuItem("Prod Network Snapshot", "send network snapshot to notifiarr"))
	c.menu["netw_dev"].Disable()  // these are not ready yet.
	c.menu["netw_test"].Disable() // these are not ready yet.
	c.menu["netw_prod"].Disable() // these are not ready yet.

	// These start hidden.
	c.menu["update"] = ui.WrapMenu(systray.AddMenuItem("Update", "Check GitHub for Update"))
	c.menu["dninfo"] = ui.WrapMenu(systray.AddMenuItem("Info!", "info from DiscordNotifier.com"))
	c.menu["alert"] = ui.WrapMenu(systray.AddMenuItem("Alert!", "alert from DiscordNotifier.com"))

	c.menu["exit"] = ui.WrapMenu(systray.AddMenuItem("Quit", "Exit "+c.Flags.Name()))
}

func (c *Client) watchKillerChannels() {
	defer systray.Quit() // this kills the app

	for {
		select {
		case sigc := <-c.sighup:
			c.Printf("Caught Signal: %v (reloading configuration)", sigc)
			c.reloadConfiguration("caught signal " + sigc.String())
		case sigc := <-c.sigkil:
			c.Errorf("Need help? %s\n=====> Exiting! Caught Signal: %v", helpLink, sigc)
			return
		case <-c.menu["exit"].Clicked():
			c.Errorf("Need help? %s\n=====> Exiting! User Requested", helpLink)
			return
		}
	}
}

// nolint:errcheck,cyclop
func (c *Client) watchGuiChannels() {
	for {
		select {
		case <-c.menu["stat"].Clicked():
			c.toggleServer()
		case <-c.menu["gh"].Clicked():
			ui.OpenURL("https://github.com/Go-Lift-TV/discordnotifier-client/")
		case <-c.menu["hp"].Clicked():
			ui.OpenURL("https://discordnotifier.com/")
		case <-c.menu["wiki"].Clicked():
			ui.OpenURL("https://trash-guides.info/Misc/Discord-Notifier-Basic-Setup/")
		case <-c.menu["disc1"].Clicked():
			ui.OpenURL("https://discord.gg/AURf8Yz")
		case <-c.menu["disc2"].Clicked():
			ui.OpenURL("https://golift.io/discord")
		case <-c.menu["view"].Clicked():
			ui.Info(Title+": Configuration", c.displayConfig())
		case <-c.menu["edit"].Clicked():
			c.Print("User Editing Config File:", c.Flags.ConfigFile)
			ui.OpenFile(c.Flags.ConfigFile)
		case <-c.menu["load"].Clicked():
			c.reloadConfiguration("UI requested")
		case <-c.menu["key"].Clicked():
			c.changeKey()
		case <-c.menu["logs_view"].Clicked():
			c.Print("User Viewing Log File:", c.Config.LogFile)
			ui.OpenLog(c.Config.LogFile)
		case <-c.menu["logs_http"].Clicked():
			c.Print("User Viewing Log File:", c.Config.HTTPLog)
			ui.OpenLog(c.Config.HTTPLog)
		case <-c.menu["logs_rotate"].Clicked():
			c.rotateLogs()
		case <-c.menu["snap_log"].Clicked():
			c.logSnaps()
		case <-c.menu["plex_test"].Clicked():
			c.sendPlexSessions(notifiarr.TestURL)
		case <-c.menu["snap_test"].Clicked():
			c.sendSystemSnapshot(notifiarr.TestURL)
		case <-c.menu["plex_dev"].Clicked():
			c.sendPlexSessions(notifiarr.DevURL)
		case <-c.menu["snap_dev"].Clicked():
			c.sendSystemSnapshot(notifiarr.DevURL)
		case <-c.menu["plex_prod"].Clicked():
			c.sendPlexSessions(notifiarr.ProdURL)
		case <-c.menu["snap_prod"].Clicked():
			c.sendSystemSnapshot(notifiarr.ProdURL)
		case <-c.menu["update"].Clicked():
			c.checkForUpdate()
		case <-c.menu["dninfo"].Clicked():
			c.menu["dninfo"].Hide()
			ui.Info(Title, "INFO: "+c.info)
		}
	}
}