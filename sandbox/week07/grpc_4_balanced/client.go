package grpc_4_balanced

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	// "gws/7/microservices/grpc/session"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	// consulAddr = flag.String("addr", "127.0.0.1:8500", "consul addr (8500 in original consul)")
	consul *consulapi.Client
)

func MainClient() {
	flag.Parse()
	var err error
	config := consulapi.DefaultConfig()
	config.Address = *consulAddr

	consul, err = consulapi.NewClient(config)
	if err != nil {
		log.Fatalf("client can't connect to consul")
	}

	health, _, err := consul.Health().Service("session-api", "", false, nil)
	if err != nil {
		log.Fatalf("client can't get alive services")
	}

	// build resolver
	servers := make([]resolver.Address, 0, len(health))
	currAddrs := []string{}
	for _, item := range health {
		addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
		currAddrs = append(currAddrs, addr)
		servers = append(servers, resolver.Address{Addr: addr})
	}

	nameResolver := manual.NewBuilderWithScheme("session-api")
	nameResolver.InitialState(resolver.State{
		Addresses: servers,
	})

	// connect
	grcpConn, err := grpc.Dial(
		nameResolver.Scheme()+":///",
		grpc.WithResolvers(nameResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
		grpc.WithBlock(),
		// grpc.WithInsecure(), deprecated
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("client can't connect to grpc")
	}
	fmt.Println("client created grpc conn, servers", servers)
	defer grcpConn.Close()

	// тут мы будем периодически опрашивать консул на предмет изменений
	go runOnlineServiceDiscovery(nameResolver, currAddrs)

	sessManager := NewAuthCheckerClient(grcpConn)
	ctx := context.Background()
	step := 1
	for {
		// проверяем несуществуюущую сессию
		// потому что сейчас между сервисами нет общения
		// получаем загшулку
		sess, err := sessManager.Check(ctx, &SessionID{ID: "not_exist_" + strconv.Itoa(step)})
		fmt.Println("client get sess", step, sess, err)

		time.Sleep(999 * time.Millisecond) // query every ~1 sec
		step++
	}
}

// async resolver updater
func runOnlineServiceDiscovery(nameResolver *manual.Resolver, currServers []string) {
	currAddrs := make(map[string]struct{}, len(currServers))
	for _, addr := range currServers {
		currAddrs[addr] = struct{}{}
	}

	ticker := time.Tick(1700 * time.Millisecond) // update every ~2 sec.
	for _ = range ticker {                       // infinity
		health, _, err := consul.Health().Service("session-api", "", false, nil)
		if err != nil {
			log.Fatalf("client can't get alive services")
		}

		newAddrs := make(map[string]struct{}, len(health))
		newServers := make([]resolver.Address, 0, len(health))

		for _, item := range health {
			addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
			// fmt.Println("DEBUG: health addr, port", item.Service.Address, item.Service.Port)
			newAddrs[addr] = struct{}{}
			newServers = append(newServers, resolver.Address{Addr: addr})
		}
		fmt.Println("client, servers alive", len(health), newServers)

		updates := 0
		// проверяем что удалилось
		for addr := range currAddrs {
			if _, exist := newAddrs[addr]; !exist {
				updates++
				// delete(currAddrs, addr)
				fmt.Println("client remove server", addr)
			}
		}
		// проверяем что добавилось
		for addr := range newAddrs {
			if _, exist := currAddrs[addr]; !exist {
				updates++
				// currAddrs[addr] = struct{}{}
				fmt.Println("client add server", addr)
			}
		}
		if updates > 0 {
			// nameResolver.CC.NewAddress(servers) deprecated
			nameResolver.UpdateState(resolver.State{
				Addresses: newServers,
			})
		}
		currAddrs = newAddrs
	}
}
