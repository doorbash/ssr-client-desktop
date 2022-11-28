package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx               context.Context
	githubUrl         string
	baseurl           string
	clientFileName    string
	clientDir         string
	downloadLock      *sync.Mutex
	clientFileHandler *ClientFileHandler
	cancelDownload    context.CancelFunc
	dbName            string
	dbHandler         *DBHandler
	proxyHandler      *ProxyHandler
}

// NewApp creates a new App application struct
func NewApp(githubUrl, baseurl, clientFileName, clientDir, dbName string) *App {
	return &App{
		githubUrl:      githubUrl,
		baseurl:        baseurl,
		clientFileName: clientFileName,
		clientDir:      clientDir,
		dbName:         dbName,
		downloadLock:   &sync.Mutex{},
	}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
	b.dbHandler = NewDBHandler(fmt.Sprintf("%s/%s", b.clientDir, b.dbName), 3*time.Second)
	err := b.dbHandler.Init(ctx)
	if err != nil {
		log.Println(err)
	}
	b.clientFileHandler = NewClientFileHandler(
		ctx,
		b.baseurl,
		b.clientFileName,
		b.clientDir,
	)
	b.proxyHandler = NewProxyHandler(fmt.Sprintf("%s/%s", b.clientDir, b.clientFileName), func(id int64, runStatus string) {
		runtime.EventsEmit(b.ctx, "run-status", id, runStatus)
	}, func(id int64, time int64, _type string, message string) {
		runtime.EventsEmit(b.ctx, "run-log", id, time, _type, message)
	})
	b.proxyHandler.Start(b.ctx, b.dbHandler)
	go func() {
		systray.Run(func() {
			systray.SetIcon(Icon)
			systray.SetTitle(APP_TITLE)
			systray.SetTooltip(APP_TITLE)
			mOpen := systray.AddMenuItem(fmt.Sprintf("Show %s", APP_TITLE), "Show App")
			mQuit := systray.AddMenuItem("Exit", "Exit app")
			go func() {
				for {
					select {
					case <-mQuit.ClickedCh:
						systray.Quit()
					case <-mOpen.ClickedCh:
						runtime.WindowShow(b.ctx)
					}
				}
			}()
		}, func() {
			runtime.Quit(b.ctx)
		})
	}()
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
	runtime.EventsOn(
		ctx,
		"event-from-js",
		func(optionalData ...interface{}) {
			if len(optionalData) == 0 {
				return
			}
			key, ok := optionalData[0].(string)
			if ok {
				b.onEvent(key, nil)
				return
			}
			list := optionalData[0].([]interface{})
			if len(list) == 0 {
				return
			}
			key, ok = list[0].(string)
			if !ok {
				return
			}
			if len(list) == 1 {
				b.onEvent(key, nil)
				return
			}
			b.onEvent(key, list[1])
		})
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	runtime.EventsOff(ctx, "event-from-js")
	// err := KillProcessByName(b.clientFileName) // Kill all child processes by their name!
	b.proxyHandler.Stop()
}

func (b *App) onEvent(key string, value interface{}) {
	fmt.Println("New Event!! key:", key, "value:", value)
}

func (b *App) ClientFileExists() bool {
	if b.clientFileHandler == nil {
		return false
	}
	return b.clientFileHandler.ClientFileExists()
}

func (b *App) DownloadClientFile() {
	if b.clientFileHandler == nil {
		return
	}
	go func() {
		b.downloadLock.Lock()
		defer b.downloadLock.Unlock()
		c, cancel := context.WithCancel(b.ctx)
		b.cancelDownload = cancel
		err := b.clientFileHandler.DownloadClientFile(c, func(progress int) {
			runtime.EventsEmit(b.ctx, "download-report", "progress", int(progress))
		})
		if err != nil {
			if err == context.Canceled {
				runtime.EventsEmit(b.ctx, "download-report", "result", "canceled")
			} else {
				runtime.EventsEmit(b.ctx, "download-report", "result", "error", err.Error())
			}
			return
		}
		runtime.EventsEmit(b.ctx, "download-report", "result", "success")
	}()
}

func (b *App) CancelDownload() {
	if b.cancelDownload != nil {
		b.cancelDownload()
	}
}

func (b *App) OpenGithub() {
	OpenBrowser(b.githubUrl)
}

func (b *App) InsertProxy(p *Proxy) error {
	id, err := b.dbHandler.InsertProxy(b.ctx, p)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("id is 0")
	}
	p.Id = id
	if p.Status == 1 {
		b.proxyHandler.AddProxy(p)
	}
	return nil
}

func (b *App) UpdateProxy(p *Proxy) error {
	b.proxyHandler.RemoveProxy(p.Id)
	err := b.dbHandler.UpdateProxy(b.ctx, p)
	if err != nil {
		return err
	}
	p, err = b.dbHandler.GetProxy(b.ctx, p.Id)
	if err != nil {
		return err
	}
	if p.Status == 1 {
		b.proxyHandler.AddProxy(p)
	}
	return nil
}

func (b *App) DeleteProxy(id int64) error {
	b.proxyHandler.RemoveProxy(id)
	b.proxyHandler.DeleteLogs(id)
	return b.dbHandler.DeleteProxy(b.ctx, id)
}

func (b *App) GetProxies() []Proxy {
	list, err := b.dbHandler.GetProxies(b.ctx, false)
	if err != nil {
		return []Proxy{}
	}
	plist := b.proxyHandler.proxyList.GetAll()
	log.Println("plist:", plist)

	for i := range list {
		for j := range plist {
			if plist[j] != nil && list[i].Id == plist[j].Id {
				list[i].RunStatus = plist[j].RunStatus
				break
			}
		}
	}

	return list
}

func (b *App) RunProxy(id int64) error {
	fmt.Println("run called mother fucker")
	p, err := b.dbHandler.GetProxy(b.ctx, id)
	if err != nil {
		return err
	}
	p.Status = 1
	err = b.dbHandler.UpdateProxy(b.ctx, p)
	if err != nil {
		return err
	}
	b.proxyHandler.AddProxy(p)
	return nil
}

func (b *App) StopProxy(id int64) error {
	p, err := b.dbHandler.GetProxy(b.ctx, id)
	if err != nil {
		return err
	}
	p.Status = 0
	err = b.dbHandler.UpdateProxy(b.ctx, p)
	if err != nil {
		return err
	}
	b.proxyHandler.RemoveProxy(p.Id)
	return nil
}

func (b *App) GetLogs(id int64) []*Log {
	return b.proxyHandler.GetLogs(id)
}
