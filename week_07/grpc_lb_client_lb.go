package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"coursera/microservices/grpc/session"

	consulapi "github.com/hashicorp/consul/api"
)

// demo: ask our service using consul registry and load balancing

var (
	// commandline parameters
	consulAddr = flag.String("addr", "192.168.99.100:32769", "consul addr (8500 in original consul)")
)

var (
	// global refs
	consul       *consulapi.Client
	nameResolver *testNameResolver // mutable, updated every 5 sec
)

func main() {
	flag.Parse() // commandline params

	// ask consul: list of services
	var err error
	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err = consulapi.NewClient(config)

	health, _, err := consul.Health().Service("session-api", "", false, nil) // by name, should use health-check procedure
	if err != nil {
		log.Fatalf("cant get alive services")
	}

	servers := []string{} // collect service addr info
	for _, item := range health {
		addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
		servers = append(servers, addr)
	}

	// create global var
	nameResolver = &testNameResolver{
		addr: servers[0],
	}

	// grpc connection with load-balancer
	grcpConn, err := grpc.Dial(
		servers[0],
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithBalancer(grpc.RoundRobin(nameResolver)),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	// update services list
	if len(servers) > 1 {
		var updates []*naming.Update
		for i := 1; i < len(servers); i++ {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: servers[i],
			})
		}
		nameResolver.w.inject(updates)
	}

	// тут мы будем периодически опрашивать консул на предмет изменений
	go runOnlineServiceDiscovery(servers) // every 5 sec update list of services from consul

	// ask service, workload imitation
	sessManager := session.NewAuthCheckerClient(grcpConn)
	ctx := context.Background()
	step := 1
	for {
		// проверяем несуществуюущую сессию
		// потому что сейчас между сервисами нет общения
		// получаем загшулку
		sess, err := sessManager.Check(ctx,
			&session.SessionID{
				ID: "not_exist_" + strconv.Itoa(step),
			})
		fmt.Println("get sess", step, sess, err)

		time.Sleep(1500 * time.Millisecond)
		step++
	}
}

func runOnlineServiceDiscovery(servers []string) {
	// mutate global var nameResolver

	// create map addr => interface // why not set?
	currAddrs := make(map[string]struct{}, len(servers))
	for _, addr := range servers {
		currAddrs[addr] = struct{}{}
	}

	// eternal go-channel
	ticker := time.Tick(5 * time.Second)
	for _ = range ticker {
		// ask consul
		health, _, err := consul.Health().Service("session-api", "", false, nil)
		if err != nil {
			log.Fatalf("cant get alive services")
		}
		// set of addr from consul
		newAddrs := make(map[string]struct{}, len(health))
		for _, item := range health {
			addr := item.Service.Address +
				":" + strconv.Itoa(item.Service.Port)
			newAddrs[addr] = struct{}{}
		}

		// fill list of updates

		var updates []*naming.Update
		// проверяем что удалилось
		for addr := range currAddrs {
			if _, exist := newAddrs[addr]; !exist {
				updates = append(updates, &naming.Update{
					Op:   naming.Delete,
					Addr: addr,
				})
				delete(currAddrs, addr)
				fmt.Println("remove", addr)
			}
		}
		// проверяем что добавилось
		for addr := range newAddrs {
			if _, exist := currAddrs[addr]; !exist {
				updates = append(updates, &naming.Update{
					Op:   naming.Add,
					Addr: addr,
				})
				currAddrs[addr] = struct{}{}
				fmt.Println("add", addr)
			}
		}

		// update global resolver
		if len(updates) > 0 {
			nameResolver.w.inject(updates)
		}
	}
}
