package container

import (
	"html/template"
	"regexp"
	"sync"
	"time"

	"strings"

	"github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/filters"
	"github.com/containrrr/watchtower/pkg/types"
	"github.com/shemanaev/updtr/internal/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var reGithub = regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)/?`)

type Client struct {
	mapping      []config.ImagesChangelog
	stale        []Info
	lastUpdate   time.Time
	staleMutex   sync.RWMutex
	refreshMutex sync.Mutex
	client       container.Client
	filter       types.Filter
	timeout      time.Duration
}

func NewClient(cfg config.Settings, mapping []config.ImagesChangelog) *Client {
	cl := container.NewClient(
		!cfg.NoPull,
		cfg.IncludeStopped,
		cfg.ReviveStopped,
		cfg.RemoveVolumes,
		cfg.IncludeRestarting,
		cfg.WarnOnHeadFailure,
	)

	filter, _ := filters.BuildFilter([]string{}, cfg.EnableLabel, cfg.Scope)

	client := Client{
		mapping:    mapping,
		stale:      []Info{},
		lastUpdate: time.Now(),
		client:     cl,
		filter:     filter,
		timeout:    cfg.StopTimeout,
	}
	return &client
}

func (c *Client) GetLastUpdateTime() time.Time {
	return c.lastUpdate
}

func (c *Client) GetStaleContainers() []Info {
	c.staleMutex.RLock()
	defer c.staleMutex.RUnlock()

	return c.stale
}

func (c *Client) Update(ids []string) {
	c.staleMutex.Lock()
	defer c.staleMutex.Unlock()

	for _, id := range ids {
		idx := slices.IndexFunc(c.stale, func(v Info) bool { return v.ID == id })
		if idx == -1 {
			log.Warnf("container with id %s not found", id)
			continue
		}

		if c.stale[idx].State != RequestedUpdate {
			c.stale[idx].State = RequestedUpdate
			c.lastUpdate = time.Now()
			log.Tracef("Requested container update for %s", id)
		}
	}
}

func (c *Client) RefreshContainers() {
	if !c.refreshMutex.TryLock() {
		log.Trace("refresh mutex already locked. ignoring refresh")
		return
	}
	defer c.refreshMutex.Unlock()

	log.Info("Refreshing containers...")
	start := time.Now()
	containers, err := c.client.ListContainers(c.filter)
	if err != nil {
		log.Error(err)
		return
	}

	var wg sync.WaitGroup
	ci := []Info{}
	for _, targetContainer := range containers {
		stale, latestImage, err := c.client.IsContainerStale(targetContainer)
		if err != nil {
			log.Infof("Unable to update container %q: %v. Proceeding to next.", targetContainer.Name(), err)
			stale = false
		}

		if !stale {
			continue
		}

		c.staleMutex.RLock()
		idx := slices.IndexFunc(c.stale, func(c Info) bool { return c.ID == targetContainer.ContainerInfo().ID })
		if idx != -1 && c.stale[idx].LatestImage == string(latestImage) {
			// reuse info if image wasn't updated since last check
			ci = append(ci, c.stale[idx])
			c.staleMutex.RUnlock()
			continue
		}
		c.staleMutex.RUnlock()

		imageName := targetContainer.ImageName()
		containerName := strings.TrimPrefix(targetContainer.Name(), "/")
		container := Info{
			ID:          targetContainer.ContainerInfo().ID,
			Name:        containerName,
			ImageName:   imageName,
			LatestImage: string(latestImage),
			State:       Stale,
			Container:   targetContainer,
		}
		ci = append(ci, container)

		idx = slices.IndexFunc(c.mapping, func(i config.ImagesChangelog) bool { return slices.Contains(i.Names, imageName) })
		if idx == -1 {
			// try with untagged name
			imageNameWithoutTag := strings.Split(imageName, ":")[0]
			idx = slices.IndexFunc(c.mapping, func(i config.ImagesChangelog) bool { return slices.Contains(i.Names, imageNameWithoutTag) })
		}

		if idx != -1 {
			containerIdx := len(ci) - 1
			changelogType := c.mapping[idx].Type
			changelogUrl := c.mapping[idx].Url

			wg.Add(1)
			go func() {
				defer wg.Done()

				changelog, err := c.getChangelog(changelogType, changelogUrl)
				if err != nil {
					log.Errorf("can't retrieve log: %v", err)
					changelog = "<pre>error retrieving changelog</pre>"
				}
				ci[containerIdx].Changelog = template.HTML(changelog)
			}()
		}
	}

	wg.Wait()

	c.staleMutex.Lock()
	c.stale = ci
	c.staleMutex.Unlock()

	c.lastUpdate = time.Now()
	elapsed := time.Since(start)
	log.Infof("Refreshing took %s", elapsed)
}

func (c *Client) DoUpdate() {
	// TODO: self update
	for i, container := range c.stale {
		if container.State != RequestedUpdate {
			continue
		}

		c.staleMutex.Lock()
		c.stale[i].State = Updating
		c.lastUpdate = time.Now()
		c.staleMutex.Unlock()

		if err := c.client.StopContainer(container.Container, c.timeout); err != nil {
			log.Error(err)
			c.staleMutex.Lock()
			c.stale[i].State = Error
			c.lastUpdate = time.Now()
			c.staleMutex.Unlock()
			continue
		}

		if _, err := c.client.StartContainer(container.Container); err != nil {
			log.Error(err)
			c.staleMutex.Lock()
			c.stale[i].State = Error
			c.lastUpdate = time.Now()
			c.staleMutex.Unlock()
			continue
		}

		c.staleMutex.Lock()
		c.stale[i].State = Fresh
		c.lastUpdate = time.Now()
		c.staleMutex.Unlock()

		log.Tracef("Container %s successfully updated", container.Name)
	}
}
