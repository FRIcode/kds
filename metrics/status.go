package metrics

type Status string

const (
	StatusRunning Status = "running"
	StatusFailed  Status = "failed"
	StatusSuccess Status = "success"
)

type StatusEntry struct {
	ID     string `json:"id"`
	Status Status `json:"status"`
}

var statusArrayIndex = 0
var statusArrayLen = 10
var statusArray = make([]StatusEntry, statusArrayLen)

func AddStatusEntry(entry StatusEntry) {
	statusArray[statusArrayIndex] = entry
	statusArrayIndex = (statusArrayIndex + 1) % statusArrayLen
}

func UpdateStatusEntry(id string, status Status) {
	for i := 0; i < statusArrayLen; i++ {
		if statusArray[i].ID == id {
			statusArray[i].Status = status
			return
		}
	}
}

func GetStatusEntry(id string) *StatusEntry {
	for i := 0; i < statusArrayLen; i++ {
		if statusArray[i].ID == id {
			return &statusArray[i]
		}
	}
	return nil
}
