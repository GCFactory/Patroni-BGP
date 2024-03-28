package patroni_bgp

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

type PatroniWatcher struct {
	Url        string
	httpClient *http.Client

	urlToCheck map[string]int

	resultChan chan int
	stopChan   chan int
}

func NewPatroniWatcher(url string) (w *PatroniWatcher) {
	return &PatroniWatcher{
		Url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 60 * 60 * 600,
		},
		urlToCheck: map[string]int{},
		resultChan: make(chan int),
		stopChan:   make(chan int),
	}
}

func (w *PatroniWatcher) EnablePrimaryAddress() {
	w.urlToCheck["/primary"] = PatroniStateMaster
}

func (w *PatroniWatcher) EnableAsyncReplicaAddress() {
	w.urlToCheck["/asynchronous"] = PatroniStateAsyncReplica
}
func (w *PatroniWatcher) EnableSyncReplicaAddress() {
	w.urlToCheck["/synchronous"] = PatroniStateSyncReplica
}

func (w *PatroniWatcher) EnableReplicaAddress() {
	w.urlToCheck["/replica"] = PatroniStateReplica
}

func (w *PatroniWatcher) ResultChan() *chan int {
	return &w.resultChan
}

func (w *PatroniWatcher) Stop() {
	w.stopChan <- 1
}

func (w *PatroniWatcher) Start() {
	go func() {
		defer func() {
			close(w.resultChan)
		}()
		r := rand.New(rand.NewSource(99))
		c := time.Tick(10 * time.Second)

		for range c {
			select {
			case <-w.stopChan:
				close(w.stopChan)
				return
			default:
				break
			}
			for path, state := range w.urlToCheck {
				log.Infof("called for %s", path)
				url := w.Url + path
				req, err := http.NewRequest("GET", url, nil)
				res, err := w.httpClient.Do(req)
				log.Infof("called for %v response %v", path, res.StatusCode)
				if err != nil {
					log.Warnf("error on patroni check: %s", err)
					w.resultChan <- PatroniStateError
				} else {
					if res.StatusCode == http.StatusOK {
						w.resultChan <- state
						break
					}
				}
			}
			jitter := time.Duration(r.Int31n(5000)) * time.Millisecond
			time.Sleep(jitter)
		}
	}()
	return
}
