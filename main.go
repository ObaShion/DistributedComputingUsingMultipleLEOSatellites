package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	ID    int
	Input int
}

type Result struct {
	JobID   int
	Output  int
	Handled []string // 経由した衛星名
}

// 衛星ノード（ワーカー）構造体
type Satellite struct {
	Name string
}

func (s *Satellite) process(job Job, results chan<- Result) {
	// 処理のシミュレーション
	fmt.Printf("[%s] processing job %d\n", s.Name, job.ID)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)+50))

	// 次の衛星に転送（ここでは2段階目と仮定）
	relaySatellite := Satellite{Name: s.Name + "-relay"}
	fmt.Printf("[%s] relaying job %d to %s\n", s.Name, job.ID, relaySatellite.Name)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)+50))

	// 地上局への結果送信
	output := job.Input * job.Input
	results <- Result{
		JobID:   job.ID,
		Output:  output,
		Handled: []string{s.Name, relaySatellite.Name, "Ground Station"},
	}
}

func main() {
	const numJobs = 2
	satellites := []Satellite{
		{Name: "LEO-1"},
		{Name: "LEO-2"},
		{Name: "LEO-3"},
	}

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// 衛星ネットワークで処理を分担
	for _, sat := range satellites {
		s := sat // ローカルコピー
		go func() {
			for job := range jobs {
				s.process(job, results)
			}
		}()
	}

	// ジョブの投入
	for i := 1; i <= numJobs; i++ {
		jobs <- Job{ID: i, Input: i}
	}
	close(jobs)

	// 結果を集約（地上局）
	for i := 1; i <= numJobs; i++ {
		res := <-results
		fmt.Printf("[Ground Station] Job %d result: %d | Path: %v\n",
			res.JobID, res.Output, res.Handled)
	}
}
