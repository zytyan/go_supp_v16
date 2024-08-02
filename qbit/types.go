package qbit

type Torrent struct {
	AddedOn                  int     `json:"added_on"`
	AmountLeft               int     `json:"amount_left"`
	AutoTmm                  bool    `json:"auto_tmm"`
	Availability             float64 `json:"availability"`
	Category                 string  `json:"category"`
	Completed                int     `json:"completed"`
	CompletionOn             int     `json:"completion_on"`
	ContentPath              string  `json:"content_path"`
	DlLimit                  int     `json:"dl_limit"`
	DlSpeed                  int     `json:"dlspeed"`
	DownloadPath             string  `json:"download_path"`
	Downloaded               int     `json:"downloaded"`
	DownloadedSession        int     `json:"downloaded_session"`
	Eta                      int     `json:"eta"`
	FirstLastPiecePriority   bool    `json:"f_l_piece_prio"`
	ForceStart               bool    `json:"force_start"`
	Hash                     string  `json:"hash"`
	InactiveSeedingTimeLimit int     `json:"inactive_seeding_time_limit"`
	InfoHashV1               string  `json:"infohash_v1"`
	InfoHashV2               string  `json:"infohash_v2"`
	LastActivity             int     `json:"last_activity"`
	MagnetURI                string  `json:"magnet_uri"`
	MaxInactiveSeedingTime   int     `json:"max_inactive_seeding_time"`
	MaxRatio                 float64 `json:"max_ratio"`
	MaxSeedingTime           int     `json:"max_seeding_time"`
	Name                     string  `json:"name"`
	NumComplete              int     `json:"num_complete"`
	NumIncomplete            int     `json:"num_incomplete"`
	NumLeechers              int     `json:"num_leechs"`
	NumSeeds                 int     `json:"num_seeds"`
	Priority                 int     `json:"priority"`
	Progress                 float64 `json:"progress"`
	Ratio                    float64 `json:"ratio"`
	RatioLimit               float64 `json:"ratio_limit"`
	SavePath                 string  `json:"save_path"`
	SeedingTime              int     `json:"seeding_time"`
	SeedingTimeLimit         int     `json:"seeding_time_limit"`
	SeenComplete             int     `json:"seen_complete"`
	SeqDl                    bool    `json:"seq_dl"`
	Size                     int     `json:"size"`
	State                    string  `json:"state"`
	SuperSeeding             bool    `json:"super_seeding"`
	Tags                     string  `json:"tags"`
	TimeActive               int     `json:"time_active"`
	TotalSize                int     `json:"total_size"`
	Tracker                  string  `json:"tracker"`
	TrackersCount            int     `json:"trackers_count"`
	UpLimit                  int     `json:"up_limit"`
	Uploaded                 int     `json:"uploaded"`
	UploadedSession          int     `json:"uploaded_session"`
	UpSpeed                  int     `json:"upspeed"`
}

const (
	StatusError              = "error"
	StatusMissingFiles       = "missingFiles"
	StatusUploading          = "uploading"
	StatusPausedUp           = "pausedUP"
	StatusQueuedUp           = "queuedUP"
	StatusStalledUp          = "stalledUP"
	StatusCheckingUp         = "checkingUP"
	StatusForcedUp           = "forcedUP"
	StatusAllocating         = "allocating"
	StatusDownloading        = "downloading"
	StatusMetaDL             = "metaDL"
	StatusPausedDL           = "pausedDL"
	StatusQueuedDL           = "queuedDL"
	StatusStalledDL          = "stalledDL"
	StatusCheckingDL         = "checkingDL"
	StatusForcedDL           = "forcedDL"
	StatusCheckingResumeData = "checkingResumeData"
	StatusMoving             = "moving"
	StatusUnknown            = "unknown"
)
