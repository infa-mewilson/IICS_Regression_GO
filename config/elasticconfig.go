package config

import "time"

type Response struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			Index  string `json:"_index"`
			ID     string `json:"_id"`
			Source struct {
				Message string `json:"message"`
				Log     struct {
					File struct {
						Path string `json:"path"`
					} `json:"file"`
					Offset int `json:"offset"`
				} `json:"log"`
				Min                 int      `json:"Min"`
				Throughput          float64  `json:"Throughput"`
				Median              int      `json:"Median"`
				Tags                []string `json:"tags"`
				Error               string   `json:"Error"`
				ExecutionEndTime    string   `json:"Execution_End_Time"`
				ScriptName          string   `json:"ScriptName"`
				TestDurationSeconds int      `json:"Test_Duration_Seconds"`
				ExecutionStartTime  string   `json:"Execution_Start_Time"`
				Ecs                 struct {
					Version string `json:"version"`
				} `json:"ecs"`
				Version      string    `json:"@version"`
				Iteration    int       `json:"Iteration"`
				Nine0Th      int       `json:"90th"`
				ReceivedKBps float64   `json:"ReceivedKBps"`
				Timestamp    time.Time `json:"@timestamp"`
				Event        struct {
					Original string `json:"original"`
				} `json:"event"`
				Nine9Th int `json:"99th"`
				Average int `json:"Average"`
				Agent   struct {
					Version     string `json:"version"`
					Type        string `json:"type"`
					EphemeralID string `json:"ephemeral_id"`
					ID          string `json:"id"`
					Name        string `json:"name"`
				} `json:"agent"`
				Nine5Th       int     `json:"95th"`
				ReleaseNumber float64 `json:"ReleaseNumber"`
				Host          struct {
					IP  []string `json:"ip"`
					Mac []string `json:"mac"`
					Os  struct {
						Platform string `json:"platform"`
						Family   string `json:"family"`
						Name     string `json:"name"`
						Kernel   string `json:"kernel"`
						Codename string `json:"codename"`
						Type     string `json:"type"`
						Version  string `json:"version"`
					} `json:"os"`
					Name          string `json:"name"`
					Architecture  string `json:"architecture"`
					Containerized bool   `json:"containerized"`
					ID            string `json:"id"`
					Hostname      string `json:"hostname"`
				} `json:"host"`
				Input struct {
					Type string `json:"type"`
				} `json:"input"`
				Samples           int     `json:"Samples"`
				BuildNumber       int     `json:"BuildNumber"`
				Release_Iteration float64 `json:"Release_Iteration"`
				Max               int     `json:"Max"`
				UserLoadThreads   int     `json:"UserLoad_Threads"`
				Label             string  `json:"Label"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
