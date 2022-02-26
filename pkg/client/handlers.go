package client

import (
	"bytes"
	"expvar"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/Notifiarr/notifiarr/pkg/bindata"
	"github.com/Notifiarr/notifiarr/pkg/exp"
	"github.com/Notifiarr/notifiarr/pkg/notifiarr"
	"github.com/Notifiarr/notifiarr/pkg/ui"
	"github.com/Notifiarr/notifiarr/pkg/update"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golift.io/starr"
)

// httpHandlers initializes GUI HTTP routes.
func (c *Client) httpHandlers() {
	c.httpAPIHandlers() // Init API handlers up front.

	base := path.Join("/", c.Config.URLBase)

	c.Config.Router.Handle("/favicon.ico", http.HandlerFunc(c.favIcon))
	c.Config.Router.Handle(strings.TrimSuffix(base, "/")+"/", handlers.MethodHandler{
		"GET":  http.HandlerFunc(c.slash),
		"POST": http.HandlerFunc(c.loginHandler),
	})

	if !strings.EqualFold(base, "/") {
		// Handle the same URLs as above on the different base URL too.
		c.Config.Router.Handle(path.Join(base, "favicon.ico"), http.HandlerFunc(c.favIcon))
		c.Config.Router.Handle(base, handlers.MethodHandler{
			"GET":  http.HandlerFunc(c.slash),
			"POST": http.HandlerFunc(c.loginHandler),
		})
	}

	//nolint:lll
	if c.Config.UIPassword != "" {
		c.Config.Router.PathPrefix(path.Join(base, "/files/")).
			Handler(http.StripPrefix(strings.TrimSuffix(base, "/"), handlers.MethodHandler{"GET": http.HandlerFunc(c.handleStaticAssets)}))
		c.Config.Router.Handle(path.Join(base, "/login"), handlers.MethodHandler{
			"GET":  http.HandlerFunc(c.loginHandler),
			"POST": http.HandlerFunc(c.loginHandler),
		})
		c.Config.Router.Handle(path.Join(base, "/logout"), handlers.MethodHandler{
			"GET":  http.HandlerFunc(c.logoutHandler),
			"POST": http.HandlerFunc(c.logoutHandler),
		})
		c.Config.Router.Handle(path.Join(base, "/profile"), handlers.MethodHandler{"POST": c.checkAuthorized(c.handleProfilePost)})
		c.Config.Router.Handle(path.Join(base, "/trigger/{action}"), c.checkAuthorized(c.handleGUITrigger)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/trigger/{action}/{content}"), c.checkAuthorized(c.handleGUITrigger)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/ps"), c.checkAuthorized(c.handleProcessList)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/template/{template}"), c.checkAuthorized(c.getTemplatePageHandler)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/getFile/{source}/{id}/{lines}/{skip}"), c.checkAuthorized(c.getFileHandler)).Methods("GET").Queries("sort", "{sort}")
		c.Config.Router.Handle(path.Join(base, "/getFile/{source}/{id}/{lines}"), c.checkAuthorized(c.getFileHandler)).Methods("GET").Queries("sort", "{sort}")
		c.Config.Router.Handle(path.Join(base, "/getFile/{source}/{id}"), c.checkAuthorized(c.getFileHandler)).Methods("GET").Queries("sort", "{sort}")
		c.Config.Router.Handle(path.Join(base, "/getFile/{source}/{id}/{lines}/{skip}"), c.checkAuthorized(c.getFileHandler)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/getFile/{source}/{id}/{lines}"), c.checkAuthorized(c.getFileHandler)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/getFile/{source}/{id}"), c.checkAuthorized(c.getFileHandler)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/downloadFile/{source}/{id}"), c.checkAuthorized(c.getFileDownloadHandler)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/deleteFile/{source}/{id}"), c.checkAuthorized(c.getFileDeleteHandler)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/services/{action:stop|start}"), c.checkAuthorized(c.handleServicesStopStart)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/services/check/{service}"), c.checkAuthorized(c.handleServicesCheck)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/shutdown"), c.checkAuthorized(c.handleShutdown)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/reload"), c.checkAuthorized(c.handleReload)).Methods("GET")
		c.Config.Router.Handle(path.Join(base, "/reconfig"), c.checkAuthorized(c.handleConfigPost)).Methods("POST")
		c.Config.Router.Handle(path.Join(base, "/debug/vars"), c.checkAuthorized(expvar.Handler().ServeHTTP)).Methods("GET")
		c.Config.Router.Handle(path.Join(c.Config.URLBase, "/ws"),
			handlers.MethodHandler{"GET": c.checkAuthorized(c.handleWebSockets)}).
			Queries("source", "{source}", "fileId", "{fileId}")
	}

	// 404 (or redirect to base path) everything else
	c.Config.Router.PathPrefix("/").Handler(http.HandlerFunc(c.notFound))
}

// httpAPIHandlers initializes API routes.
func (c *Client) httpAPIHandlers() {
	c.Config.HandleAPIpath("", "version", c.versionHandler, "GET", "HEAD")
	c.Config.HandleAPIpath("", "trigger/{trigger:[0-9a-z-]+}", c.handleTrigger, "GET")
	c.Config.HandleAPIpath("", "trigger/{trigger:[0-9a-z-]+}/{content}", c.handleTrigger, "GET")
	// Aggregate handlers. Non-app specific.
	c.Config.HandleAPIpath("", "/trash/{app}", c.aggregateTrash, "POST")

	if c.Config.Plex.Configured() {
		c.Config.HandleAPIpath(starr.Plex, "sessions", c.Config.Plex.HandleSessions, "GET")
		c.Config.HandleAPIpath(starr.Plex, "kill", c.Config.Plex.HandleKillSession, "GET").
			Queries("reason", "{reason:.*}", "sessionId", "{sessionId:[0-9a-z-]+}")

		tokens := fmt.Sprintf("{token:%s|%s}", c.Config.Plex.Token, c.Config.Apps.APIKey)
		c.Config.Router.Handle("/plex",
			http.HandlerFunc(c.website.PlexHandler)).Methods("POST").Queries("token", tokens)
		c.Config.Router.Handle("/",
			http.HandlerFunc(c.website.PlexHandler)).Methods("POST").Queries("token", tokens)

		if c.Config.URLBase != "/" {
			// Allow plex to use the base url too.
			c.Config.Router.Handle(path.Join(c.Config.URLBase, "plex"),
				http.HandlerFunc(c.website.PlexHandler)).Methods("POST").Queries("token", tokens)
		}
	}
}

// notFound is the handler for paths that are not found: 404s.
func (c *Client) notFound(response http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.URL.Path, c.Config.URLBase) {
		// If the request did not have the base url, redirect.
		http.Redirect(response, request, path.Join(c.Config.URLBase, request.URL.Path), http.StatusPermanentRedirect)
		return
	}

	response.WriteHeader(http.StatusNotFound)

	if err := c.templat.ExecuteTemplate(response, "404.html", nil); err != nil {
		c.Errorf("Sending HTTP Reply: %v", err)
	}
}

// slash is the GET handler for /.
func (c *Client) slash(response http.ResponseWriter, request *http.Request) {
	c.indexPage(response, request, "")
}

func (c *Client) favIcon(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
	ico, err := bindata.Asset("files/images/favicon.ico")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, r.URL.Path, time.Now(), bytes.NewReader(ico))
}

// stripSecrets runs first to save a redacted URI in a special request header.
// The logger uses this special value to save a redacted URI in the log file.
func (c *Client) stripSecrets(next http.Handler) http.Handler {
	secrets := []string{c.Config.Apps.APIKey}
	secrets = append(secrets, c.Config.ExKeys...)
	// gather configured/known secrets.
	if c.Config.Plex != nil {
		secrets = append(secrets, c.Config.Plex.Token)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
		uri := r.RequestURI
		// then redact secrets from request.
		for _, s := range secrets {
			if s != "" {
				uri = strings.ReplaceAll(uri, s, "<redacted>")
			}
		}

		// save into a request header for the logger.
		r.Header.Set("X-Redacted-URI", uri)
		next.ServeHTTP(w, r)
	})
}

func (c *Client) countRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, path.Join(c.Config.URLBase, "api")) ||
			strings.HasPrefix(req.RequestURI, path.Join(c.Config.URLBase, "plex")) ||
			strings.HasPrefix(req.RequestURI, "/plex") {
			exp.HTTPRequests.Add("/api and /plex Requests", 1)
		} else {
			exp.HTTPRequests.Add("Non /api Requests", 1)
		}

		exp.HTTPRequests.Add("Total Requests", 1)
		next.ServeHTTP(response, req)
	})
}

// fixForwardedFor sets the X-Forwarded-For header to the client IP
// under specific circumstances.
func (c *Client) fixForwardedFor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
		ip := strings.Trim(r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")], "[]")
		if x := r.Header.Get("X-Forwarded-For"); x == "" || !c.Config.Allow.Contains(ip) {
			r.Header.Set("X-Forwarded-For", ip)
		} else if l := strings.LastIndexAny(x, ", "); l != -1 {
			r.Header.Set("X-Forwarded-For", strings.Trim(x[l:len(x)-1], ", "))
		}

		next.ServeHTTP(w, r)
	})
}

func (c *Client) handleTrigger(r *http.Request) (int, interface{}) {
	return c.runTrigger(notifiarr.EventAPI, mux.Vars(r)["trigger"], mux.Vars(r)["content"])
}

func (c *Client) runTrigger(source notifiarr.EventType, trigger, content string) (int, string) { //nolint:cyclop,funlen
	if content != "" {
		c.Debugf("Incoming API Trigger: %s (%s)", trigger, content)
	} else {
		c.Debugf("Incoming API Trigger: %s", trigger)
	}

	switch trigger {
	case "cfsync":
		c.website.Trigger.SyncCF(source)
		return http.StatusOK, "TRaSH Custom Formats and Release Profile Sync initiated."
	case "services":
		c.Config.Services.RunChecks(source)
		return http.StatusOK, "All service checks rescheduled for immediate exeution."
	case "sessions":
		if !c.Config.Plex.Configured() {
			return http.StatusNotImplemented, "Plex Sessions are not enabled."
		}

		c.website.Trigger.SendPlexSessions(source)

		return http.StatusOK, "Plex sessions triggered."
	case "stuckitems":
		c.website.Trigger.SendStuckQueueItems(source)
		return http.StatusOK, "Stuck Queue Items triggered."
	case "dashboard":
		c.website.Trigger.SendDashboardState(source)
		return http.StatusOK, "Dashboard states triggered."
	case "snapshot":
		c.website.Trigger.SendSnapshot(source)
		return http.StatusOK, "System Snapshot triggered."
	case "gaps":
		c.website.Trigger.SendGaps(source)
		return http.StatusOK, "Radarr Collections Gaps initiated."
	case "corrupt":
		err := c.website.Trigger.Corruption(source, starr.App(strings.Title(content)))
		if err != nil {
			return http.StatusBadRequest, "Corruption trigger failed: " + err.Error()
		}

		return http.StatusOK, strings.Title(content) + " corruption checks initiated."
	case "backup":
		err := c.website.Trigger.Backup(source, starr.App(strings.Title(content)))
		if err != nil {
			return http.StatusBadRequest, "Backup trigger failed: " + err.Error()
		}

		return http.StatusOK, strings.Title(content) + " backups check initiated."
	case "reload":
		c.sighup <- &update.Signal{Text: "reload http triggered"}
		return http.StatusOK, "Application reload initiated."
	case "notification":
		if content != "" {
			ui.Notify("Notification: %s", content) //nolint:errcheck
			c.Printf("NOTIFICATION: %s", content)

			return http.StatusOK, "Local Nntification sent."
		}

		return http.StatusBadRequest, "Missing notification content."
	default:
		return http.StatusBadRequest, "Unknown trigger provided:'" + trigger + "'"
	}
}
