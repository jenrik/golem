package golem

import (
	"database/sql"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Task struct {
	Link  string
	JobId string
}

var connected bool

var queueConn *amqp.Connection
var queueChan *amqp.Channel
var queueDec amqp.Queue
var queue <-chan amqp.Delivery

var storageConn *sql.DB
var stmtCreateJob, stmtSubmitTask, stmtSubmitData, stmtCompleteTask, stmtGetDefinitions *sql.Stmt

func StartWorkers(amount int) sync.WaitGroup {
	if !connected {
		panic("Tried to start workers before connection to backend services were established")
	}

	var wg sync.WaitGroup
	wg.Add(amount)

	for i := 0; i < amount; i++ {
		go func() {
			defer wg.Done()
			for d := range queue {
				var task *Task
				task = new(Task)
				err := json.Unmarshal(d.Body, task)
				if err != nil {
					log.WithFields(log.Fields{
						"delivery": d,
					}).WithError(err).Error("Received invalid task")
					d.Ack(false)
					continue
				}

				rows, err := stmtGetDefinitions.Query(task.JobId)
				if err != nil {
					rows.Close()
					panic(err)
				}
				var defsJson string
				rows.Next()
				rows.Scan(&defsJson)
				if rows.Err() != nil {
					rows.Close()
					log.WithFields(log.Fields{
						"task": task,
					}).WithError(err).Error("Error while retrieving definitions from storage")
					continue
				}
				rows.Close()
				var defs []Definition
				err = json.Unmarshal([]byte(defsJson), &defs)
				if err != nil {
					log.WithFields(log.Fields{
						"task":     task,
						"defsJson": defsJson,
					}).WithError(err).Error("Received invalid definitions json from storage")
					continue
				}

				if data, err := ScrapePage(task, &defs); err != nil {
					// Scrape errors
					log.WithFields(log.Fields{
						"task": task,
					}).WithError(err).Warn("Error during scrape")
				} else {
					if jm, err := json.Marshal(data); err != nil {
						log.WithFields(log.Fields{
							"task": task,
							// "data": data,
						}).WithError(err).Error("Failed to marshal scrape result data")
					} else {
						_, err = stmtSubmitData.Exec(
							task.JobId,
							task.Link,
							string(jm),
							nil,
						)
						if err != nil {
							log.WithError(err).Error("Failed to submit scrape result to storage")
						}
					}
				}

				stmtCompleteTask.Exec(task.JobId, task.Link)
				d.Ack(false)
			}
		}()
	}

	return wg
}

func SubmitJob(defs *[]Definition, links *[]string) (string, error) {
	if !connected {
		panic("Tried to submit a job before connection to backend services were established")
	}

	jobId := uuid.NewV4().String()
	jd, err := json.Marshal(defs)
	if err != nil {
		return "", err
	}
	_, err = stmtCreateJob.Exec(jobId, string(jd))
	if err != nil {
		return "", err
	}

	return SubmitTask(jobId, links)
}

func SubmitTask(jobId string, links *[]string) (string, error) {
	if !connected {
		panic("Tried to submit a task before connection to backend services were established")
	}

	task := new(Task)
	task.JobId = jobId
	for _, link := range *links {
		result, err := stmtSubmitTask.Exec(jobId, link)
		if err != nil {
			log.WithFields(log.Fields{
				"jobid": jobId,
			}).WithError(err).Error("Failed to submit task to storage")
			continue
		}
		c, err := result.RowsAffected()
		if err != nil {
			log.WithFields(log.Fields{
				"jobid": jobId,
			}).WithError(err).Error("Failed to get affected rows after link submission to storage")
			continue
		}
		if c == 0 {
			continue
		}
		task.Link = link
		jd, err := json.Marshal(task)
		if err != nil {
			// Note: we have a lose task in storage that doesn't exist in the queue
			log.WithError(err).Error("Failed to marshal task into json during task submission")
			continue
		}
		queueChan.Publish("", queueDec.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "golem/task",
			Body:         jd,
		})
	}

	return jobId, nil
}

func Disconnect() {
	if !connected {
		panic("Called Disconnect with being connected")
	}
	if queueConn != nil {
		queueChan.Close()
		queueConn.Close()
	}
	if storageConn != nil {
		storageConn.Close()
		stmtCreateJob.Close()
		stmtSubmitTask.Close()
		stmtSubmitData.Close()
		stmtCompleteTask.Close()
	}
}

func Connect(queueAdr string, storageAdr string) error {
	var err error

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, syscall.SIGTERM)
	go func() {
		<-signalChannel
		Disconnect()
		os.Exit(1)
	}()

	// Queue
	queueConn, err = amqp.Dial(queueAdr)
	if err != nil {
		return err
	}

	queueChan, err = queueConn.Channel()
	if err != nil {
		return err
	}

	err = queueChan.Qos(1, 0, false)
	if err != nil {
		return err
	}

	queueDec, err = queueChan.QueueDeclare("golem", true, false, false, false, nil)
	if err != nil {
		return err
	}

	queue, err = queueChan.Consume(queueDec.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// Storage
	storageConn, err := sql.Open("postgres", storageAdr)
	if err != nil {
		return err
	}

	stmtCreateJob, err = storageConn.Prepare("INSERT INTO jobs(jobId, defs) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	stmtSubmitTask, err = storageConn.Prepare("INSERT INTO pages(jobId, link, done) VALUES ($1, $2, false) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}
	stmtCompleteTask, err = storageConn.Prepare("UPDATE pages SET done=true WHERE jobId=$1 AND link=$2")
	if err != nil {
		return err
	}
	stmtSubmitData, err = storageConn.Prepare("INSERT INTO data(jobId, link, data, error) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	stmtGetDefinitions, err = storageConn.Prepare("SELECT defs FROM jobs WHERE jobId=$1")
	if err != nil {
		return err
	}

	connected = true
	return nil
}
