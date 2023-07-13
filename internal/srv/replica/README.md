# M/S replica
like redis/mysql Master/Slaver replica mode, just simple classic impl (no learner to start, no linear read, slaver replica pull bluk RESP with binlog[startLogID,LastLogID]/snapshot file from master);
1. start with replicaof or exec `replicaof` cmd
   0. when start service, latest current commit log load to commitID
   1. then slaver send `replicaaof` cmd to start connect master; between replica(master/slaver), master start replica goroutine, create connect
   2. slaver send `replconf` to master register slavers
   3. if use restart slaver send `fullsync` cmd to master, master full sync from snapshot compress file to slaver, more detail see `3`
   4. then start sync loop, send `sync/psync` cmd with slaver current latest logId(syncId), master send [lastLogID+binlog] from log store to slaver, util send ack[lastLogID] sync ok

2. RESP `w` op cmd commit to save log, wait quorum slaves to ack(sync pull binlog ok), save commitID(logID) to latest current commit log, if reach the snapshot threshold to save it

3. slaver send `fullsync` cmd to master
   1. master if don't exists snapshot file, or snapshot file has expired , create new snapshot (one connect per goroutine)
   2. if not, use lastest snapshot file (init ticker job to purge expired snapshot)
   3. then lock snapshot, create new snapshot use data kvstore (FSM) lock write to gen snapshot and iter it save to snapshot file (format: [len(compress key) | compress key | len(compress value) | compress value ...])
   4. from snapshot file read snapshot send bluk([]byte) RESP to slaver
   5. slaver receive the bluk RESP to save the dump file(reply log)
   6. then lock write to load, clear all data and load dump file write(put) to data kvstore (FSM)  

slave replica connect state:
* replConnectState: slaver needs to connect to its master
* replConnectingState: slaver-master connection is in progress
* replSyncState: perform the synchronization
* replConnectedState: slaver is online
