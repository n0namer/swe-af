package main

import (
	"context"
	"log"
	"time"

	"github.com/Agent-Field/SWE-AF/go/internal/node"
	"github.com/Agent-Field/SWE-AF/go/internal/pro"
)

func main() {
	n, err := node.BuildAgent("swe-factory-30", "8010", "SWE Factory executor (Full by default, Fast optional)")
	if err != nil {
		log.Fatalf("swe-factory-30: build agent: %v", err)
	}
	n.RegisterFactory()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var sup *pro.Supervisor
	if pro.Enabled() {
		sup = pro.Start(ctx, pro.Options{Server: n.AgentFieldServer, Token: n.Token})
	}
	runErr := n.App.Run(ctx)
	cancel()
	sup.Wait(10 * time.Second)
	if runErr != nil {
		log.Fatalf("swe-factory-30: run: %v", runErr)
	}
}
