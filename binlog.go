package main

import (
	"context"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"os"
)

func main() {
	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "snan",
		Password: "19990928",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// Start sync with specified binlog file and position
	//streamer, _ := syncer.StartSync(mysql.Position{"mysql-bin.000017", 197})

	gtidSetStr := "78f4ca0e-ce5c-11ef-ac8e-000c2906c90a:1-38"
	gtidSet, _ := mysql.ParseGTIDSet("mysql", gtidSetStr)
	streamer, _ := syncer.StartSyncGTID(gtidSet)
	for {
		ev, _ := streamer.GetEvent(context.Background())
		// Dump event
		ev.Dump(os.Stdout)
	}

}
