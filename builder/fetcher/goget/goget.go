package goget

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/go-spirit/spirit-builder/builder/fetcher"
	"github.com/go-spirit/spirit-builder/utils"
	"github.com/gogap/config"
)

type GoGetFetcher struct {
	conf config.Configuration
}

func init() {
	fetcher.RegisterFetcher("goget", NewGoGetFetcher)
}

func NewGoGetFetcher(conf config.Configuration) (f fetcher.Fetcher, err error) {

	if conf == nil {
		conf = config.NewConfig()
	}

	ft := &GoGetFetcher{
		conf: conf,
	}

	f = ft

	return
}

func (p *GoGetFetcher) Fetch(url, revision string, update bool, repoConf config.Configuration) (err error) {

	args := repoConf.GetStringList("args")
	gopath := p.conf.GetString("gopath", os.Getenv("GOPATH"))
	repoDir := path.Join(gopath, "src", url)

	exist, err := utils.DirExist(repoDir)
	if err != nil {
		return
	}

	if !exist {

		err = utils.GoGet(url, args...)
		if err != nil {
			return
		}

		update = false

		logrus.WithField("fetcher", "goget").WithField("url", url).WithField("revision", revision).Infoln("fetched")
	}

	err = utils.GitCheckout(repoDir, revision)
	if err != nil {
		return
	}

	logrus.WithField("fetcher", "goget").WithField("url", url).WithField("revision", revision).Infoln("checked out")

	if update {

		var deteched bool
		deteched, err = utils.GitDetached(repoDir)
		if err != nil {
			return
		}

		if !deteched {
			err = utils.GitPull(repoDir)
			if err != nil {
				return
			}
			logrus.WithField("fetcher", "goget").WithField("url", url).WithField("revision", revision).Infoln("updated")
		} else {
			logrus.WithField("fetcher", "goget").WithField("url", url).WithField("revision", revision).Warnln("git detetched, update skipped")
		}
	}

	return
}
