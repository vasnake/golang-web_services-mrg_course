package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/jpeg"
	"os"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
)

type ImgResizeTask struct {
	Name string
	MD5  string
}

const (
	ImageResizeQueueName = "image_resize"
	StoragePrefix        = "./images" // it's a no-no
)

var (
	rabbitAddr = flag.String("addr", "amqp://guest:guest@localhost:5672/", "rabbit addr")
	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel
	sizes      = []uint{80, 160, 320} // pic size enum
)

func main() {
	flag.Parse()
	var err error
	rabbitConn, err = amqp.Dial(*rabbitAddr)
	panicOnError("can't connect to rabbit", err)

	rabbitChan, err = rabbitConn.Channel()
	panicOnError("can't open AMQP chan", err)
	defer rabbitChan.Close()

	_, err = rabbitChan.QueueDeclare(
		ImageResizeQueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	panicOnError("can't init queue", err)

	err = rabbitChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	panicOnError("can't set QoS", err)

	// get go chan
	tasksChan, err := rabbitChan.Consume(
		ImageResizeQueueName, // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	panicOnError("can't register consumer", err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// start 10 concurrent resizers
	for i := 0; i <= 10; i++ {
		go ResizeWorker(tasksChan)
	}
	fmt.Println("worker started")
	wg.Wait() // forewer
}

func ResizeWorker(messagesChan <-chan amqp.Delivery) {
	var ack = func(t amqp.Delivery) error { return t.Ack(false) }

	for msg := range messagesChan {
		fmt.Printf("incoming task %+v\n", msg)
		ack(msg) // good or bad, ack anyway

		task := &ImgResizeTask{}
		err := json.Unmarshal(msg.Body, task)
		if err != nil {
			fmt.Println("can't unpack task json", err)
			continue
		}

		originalPath := fmt.Sprintf("%s/%s.jpg", StoragePrefix, task.MD5)
		for _, size := range sizes {
			time.Sleep(3 * time.Second) // not so fast please
			resizedPath := fmt.Sprintf("%s/%s_%d.jpg", StoragePrefix, task.MD5, size)
			err := ResizeImage(originalPath, resizedPath, size)
			if err != nil {
				fmt.Println("resize failed", err)
			}
		} // end of each size
	} // end of channel
}

func ResizeImage(originalPath string, resizedPath string, size uint) error {
	inFile, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("can't open file %s: %s", originalPath, err)
	}
	defer inFile.Close()

	img, err := jpeg.Decode(inFile)
	if err != nil {
		return fmt.Errorf("can't decode jpeg file %s", err)
	}

	resizeImage := resize.Resize(size, 0, img, resize.Lanczos3)

	outFile, err := os.Create(resizedPath)
	if err != nil {
		return fmt.Errorf("can't create file %s: %s", resizedPath, err)
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, resizeImage, nil)
}

// не используйте такой код в prod // ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
