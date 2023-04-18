// Copyright 2015 Sorint.lab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"github.com/mitchellh/copystructure"
	"github.com/sorintlab/stolon/cmd"
	"github.com/sorintlab/stolon/internal/cluster"
	"github.com/sorintlab/stolon/internal/common"
	"github.com/sorintlab/stolon/internal/dbmgr"
	"github.com/sorintlab/stolon/internal/flagutil"
	slog "github.com/sorintlab/stolon/internal/log"
	"github.com/sorintlab/stolon/internal/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var log = slog.S()

//var CmdKeeper = &cobra.Command{
//	Use:     "stolon-keeper",
//	Run:     keeper,
//	Version: cmd.Version,
//}

var RootCmd = &cobra.Command{
	Use:   util.StolonBinNameKeeper,
	Short: "stolon keeper",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// do nothing
	},
}

const (
	maxPostgresTimelinesHistory = 2
	minWalKeepSegments          = 8
)

type KeeperIF interface {
	GetInSyncStandbys() ([]string, error)
	GetPGState(pctx context.Context) (*cluster.DBMSState, error)
	Start(ctx context.Context)
}

type KeeperLocalState struct {
	UID        string
	ClusterUID string
}

type DBLocalState struct {
	UID        string
	Generation int64
	// Initializing registers when the db is initializing. Needed to detect
	// when the initialization has failed.
	Initializing bool
	// InitPGParameters contains the postgres parameter after the
	// initialization
	InitPGParameters common.Parameters
}

func (s *DBLocalState) DeepCopy() *DBLocalState {
	if s == nil {
		return nil
	}
	ns, err := copystructure.Copy(s)
	if err != nil {
		panic(err)
	}
	// paranoid test
	if !reflect.DeepEqual(s, ns) {
		panic("not equal")
	}
	return ns.(*DBLocalState)
}

//func init() {
//	cmd.AddCommonFlags(CmdKeeper, &cfg.CommonConfig)
//
//	CmdKeeper.PersistentFlags().StringVar(&cfg.uid, "id", "", "keeper uid (must be unique in the cluster and can contain only lower-case letters, numbers and the underscore character). If not provided a random uid will be generated.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.uid, "uid", "", "keeper uid (must be unique in the cluster and can contain only lower-case letters, numbers and the underscore character). If not provided a random uid will be generated.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.dataDir, "data-dir", "", "data directory")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgListenAddress, "pg-listen-address", "", "postgresql instance listening address, local address used for the postgres instance. For all network interface, you can set the value to '*'.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgAdvertiseAddress, "pg-advertise-address", "", "postgresql instance address from outside. Use it to expose ip different than local ip with a NAT networking config")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgPort, "pg-port", "5432", "postgresql instance listening port")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgAdvertisePort, "pg-advertise-port", "", "postgresql instance port from outside. Use it to expose port different than local port with a PAT networking config")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgBinPath, "pg-bin-path", "", "absolute path to postgresql binaries. If empty they will be searched in the current PATH")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgReplAuthMethod, "pg-repl-auth-method", "md5", "postgres replication user auth method. Default is md5.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgReplUsername, "pg-repl-username", "", "postgres replication user name. Required. It'll be created on db initialization. Must be the same for all keepers.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgReplPassword, "pg-repl-password", "", "postgres replication user password. Only one of --pg-repl-password or --pg-repl-passwordfile must be provided. Must be the same for all keepers.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgReplPasswordFile, "pg-repl-passwordfile", "", "postgres replication user password file. Only one of --pg-repl-password or --pg-repl-passwordfile must be provided. Must be the same for all keepers.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgSUAuthMethod, "pg-su-auth-method", "md5", "postgres superuser auth method. Default is md5.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgSUUsername, "pg-su-username", "", "postgres superuser user name. Used for keeper managed instance access and pg_rewind based synchronization. It'll be created on db initialization. Defaults to the name of the effective user running stolon-keeper. Must be the same for all keepers.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgSUPassword, "pg-su-password", "", "postgres superuser password. Only one of --pg-su-password or --pg-su-passwordfile must be provided. Must be the same for all keepers.")
//	CmdKeeper.PersistentFlags().StringVar(&cfg.pgSUPasswordFile, "pg-su-passwordfile", "", "postgres superuser password file. Only one of --pg-su-password or --pg-su-passwordfile must be provided. Must be the same for all keepers)")
//	CmdKeeper.PersistentFlags().BoolVar(&cfg.debug, "debug", false, "enable debug logging")
//
//	CmdKeeper.PersistentFlags().BoolVar(&cfg.canBeMaster, "can-be-master", true, "prevent keeper from being elected as master")
//	CmdKeeper.PersistentFlags().BoolVar(&cfg.canBeSynchronousReplica, "can-be-synchronous-replica", true, "prevent keeper from being chosen as synchronous replica")
//	CmdKeeper.PersistentFlags().BoolVar(&cfg.disableDataDirLocking, "disable-data-dir-locking", false, "disable locking on data dir. Warning! It'll cause data corruptions if two keepers are concurrently running with the same data dir.")
//
//	if err := CmdKeeper.PersistentFlags().MarkDeprecated("id", "please use --uid"); err != nil {
//		log.Fatal(err)
//	}
//	if err := CmdKeeper.PersistentFlags().MarkDeprecated("debug", "use --log-level=debug instead"); err != nil {
//		log.Fatal(err)
//	}
//}

var managedPGParameters = []string{
	"unix_socket_directories",
	"wal_keep_segments",
	"wal_keep_size",
	"hot_standby",
	"listen_addresses",
	"port",
	"max_replication_slots",
	"max_wal_senders",
	"wal_log_hints",
	"synchronous_standby_names",

	// parameters moved from recovery.conf to postgresql.conf in PostgresSQL 12
	"primary_conninfo",
	"primary_slot_name",
	"recovery_min_apply_delay",
	"restore_command",
	"recovery_target_timeline",
	"recovery_target",
	"recovery_target_lsn",
	"recovery_target_name",
	"recovery_target_time",
	"recovery_target_xid",
	"recovery_target_timeline",
	"recovery_target_action",
}

func readPasswordFromFile(filepath string) (string, error) {
	fi, err := os.Lstat(filepath)
	if err != nil {
		return "", fmt.Errorf("unable to read password from file %s: %v", filepath, err)
	}

	if fi.Mode() > 0600 {
		//TODO: enforce this by exiting with an error. Kubernetes makes this file too open today.
		log.Warnw("password file permissions are too open. This file should only be readable to the user executing stolon! Continuing...", "file", filepath, "mode", fmt.Sprintf("%#o", fi.Mode()))
	}

	pwBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("unable to read password from file %s: %v", filepath, err)
	}
	return string(pwBytes), nil
}

func sigHandler(sigs chan os.Signal, cancel context.CancelFunc) {
	s := <-sigs
	log.Debugw("got signal", "signal", s)
	shutdownSeconds.SetToCurrentTime()
	cancel()
}

func Execute() {

	cmd.RegisterVersionInfo(RootCmd)
	registerPostgresCmd(RootCmd)
	registerGreenplumCmd(RootCmd)

	if err := flagutil.SetFlagsFromEnv(RootCmd.PersistentFlags(), "STKEEPER"); err != nil {
		log.Fatal(err)
	}

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// parseSynchronousStandbyNames extracts the standby names from the
// "synchronous_standby_names" postgres parameter.
//
// Since postgres 9.6 (https://www.postgresql.org/docs/9.6/static/runtime-config-replication.html)
// `synchronous_standby_names` can be in one of two formats:
//
//	num_sync ( standby_name [, ...] )
//	standby_name [, ...]
//
// two examples for this:
//
//	2 (node1,node2)
//	node1,node2
//
// TODO(sgotti) since postgres 10 (https://www.postgresql.org/docs/10/static/runtime-config-replication.html)
// `synchronous_standby_names` can be in one of three formats:
//
//	[FIRST] num_sync ( standby_name [, ...] )
//	ANY num_sync ( standby_name [, ...] )
//	standby_name [, ...]
//
// since we are writing ourself the synchronous_standby_names we don't handle this case.
// If needed, to better handle all the cases with also a better validation of
// standby names we could use something like the parser used by postgres
func parseSynchronousStandbyNames(s string) ([]string, error) {
	var spacesSplit []string = strings.Split(s, " ")
	var entries []string
	if len(spacesSplit) < 2 {
		// We're parsing format: standby_name [, ...]
		entries = strings.Split(s, ",")
	} else {
		// We don't know yet which of the 2 formats we're parsing
		_, err := strconv.Atoi(spacesSplit[0])
		if err == nil {
			// We're parsing format: num_sync ( standby_name [, ...] )
			rest := strings.Join(spacesSplit[1:], " ")
			inBrackets := strings.TrimSpace(rest)
			if !strings.HasPrefix(inBrackets, "(") || !strings.HasSuffix(inBrackets, ")") {
				return nil, fmt.Errorf("synchronous standby string has number but lacks brackets")
			}
			withoutBrackets := strings.TrimRight(strings.TrimLeft(inBrackets, "("), ")")
			entries = strings.Split(withoutBrackets, ",")
		} else {
			// We're parsing format: standby_name [, ...]
			entries = strings.Split(s, ",")
		}
	}
	for i, e := range entries {
		entries[i] = strings.TrimSpace(e)
	}
	return entries, nil
}

func getTimeLinesHistory(dbState *cluster.DBMSState, pgm dbmgr.Manager, maxPostgresTimelinesHistory int) (cluster.PostgresTimelinesHistory, error) {
	ctlsh := cluster.PostgresTimelinesHistory{}
	// if timeline <= 1 then no timeline history file exists.
	if dbState.TimelineID > 1 {
		var tlsh []*dbmgr.TimelineHistory
		tlsh, err := pgm.GetTimelinesHistory(dbState.TimelineID)
		if err != nil {
			log.Errorw("error getting timeline history", zap.Error(err))
			return ctlsh, err
		}
		if len(tlsh) > maxPostgresTimelinesHistory {
			tlsh = tlsh[len(tlsh)-maxPostgresTimelinesHistory:]
		}
		for _, tlh := range tlsh {
			ctlh := &cluster.PostgresTimelineHistory{
				TimelineID:  tlh.TimelineID,
				SwitchPoint: tlh.SwitchPoint,
				Reason:      tlh.Reason,
			}
			ctlsh = append(ctlsh, ctlh)
		}
	}
	return ctlsh, nil
}

// IsMaster return if the db is the cluster master db.
// A master is a db that:
// * Has a master db role
// or
// * Has a standby db role with followtype external
func IsMaster(db *cluster.DB) bool {
	switch db.Spec.Role {
	case common.RoleMaster:
		return true
	case common.RoleStandby:
		if db.Spec.FollowConfig.Type == cluster.FollowTypeExternal {
			return true
		}
		return false
	default:
		panic("invalid db role in db Spec")
	}
}
