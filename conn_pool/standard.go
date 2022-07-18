package conn_pool

//
//import (
//	"context"
//	"database/sql/driver"
//	"io"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//// 参考sql和redis连接池设计的连接池
//// db,err := sql.Open("mysql", "xxx")
//
//type DB struct {
//	waitDuration int64        // Total time waited for new connections.
//	mu           sync.RWMutex // protects following fields
//	freeConn     []*driverConn  // slice for conn
//	connRequests map[uint64]chan connRequest
//	nextRequest  uint64 // Next key to use in connRequests.
//	numOpen      int    // number of opened and pending open connections
//	// Used to signal the need for new connections
//	// a goroutine running connectionOpener() reads on this chan and
//	// maybeOpenNewConnections sends on the chan (one send per needed connection)
//	// It is closed during db.Close(). The close tells the connectionOpener
//	// goroutine to exit.
//	openerCh          chan struct{}
//	closed            bool
//	maxIdle           int           // zero means defaultMaxIdleConns; negative means 0
//	maxOpen           int           // <= 0 means unlimited
//	maxLifetime       time.Duration // maximum amount of time a connection may be reused
//	cleanerCh         chan struct{}
//	waitCount         int64 // Total number of connections waited for.
//	maxIdleClosed     int64 // Total number of connections closed due to idle.
//	maxLifetimeClosed int64 // Total number of connections closed due to max free limit.
//}
//
//// conn returns a newly-opened or cached *driverConn.
//func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
//	// 先判断db是否已经关闭。
//	db.mu.Lock()
//	if db.closed {
//		db.mu.Unlock()
//		return nil, errDBClosed
//	}
//	// 注意检测context是否已经被超时等原因被取消。
//	select {
//	default:
//	case <-ctx.Done():
//		db.mu.Unlock()
//		return nil, ctx.Err()
//	}
//	lifetime := db.maxLifetime
//
//	// 这边如果在freeConn这个切片有空闲连接的话，就left pop一个出列。注意的是，这边因为是切片操作，所以需要前面需要加锁且获取后进行解锁操作。同时判断返回的连接是否已经过期。
//	numFree := len(db.freeConn)
//	if strategy == cachedOrNewConn && numFree > 0 {
//		conn := db.freeConn[0]
//		copy(db.freeConn, db.freeConn[1:])
//		db.freeConn = db.freeConn[:numFree-1]
//		conn.inUse = true
//		db.mu.Unlock()
//		if conn.expired(lifetime) {
//			conn.Close()
//			return nil, driver.ErrBadConn
//		}
//		// Lock around reading lastErr to ensure the session resetter finished.
//		conn.Lock()
//		err := conn.lastErr
//		conn.Unlock()
//		if err == driver.ErrBadConn {
//			conn.Close()
//			return nil, driver.ErrBadConn
//		}
//		return conn, nil
//	}
//
//	// 这边就是等候获取连接的重点了。当空闲的连接为空的时候，这边将会新建一个request（的等待连接 的请求）并且开始等待
//	if db.maxOpen > 0 && db.numOpen >= db.maxOpen {
//		// 下面的动作相当于往connRequests这个map插入自己的号码牌。
//		// 插入号码牌之后这边就不需要阻塞等待继续往下走逻辑。
//		req := make(chan connRequest, 1)
//		reqKey := db.nextRequestKeyLocked()
//		db.connRequests[reqKey] = req
//		db.waitCount++
//		db.mu.Unlock()
//
//		waitStart := time.Now()
//
//		// Timeout the connection request with the context.
//		select {
//		case <-ctx.Done():
//			// context取消操作的时候，记得从connRequests这个map取走自己的号码牌。
//			db.mu.Lock()
//			delete(db.connRequests, reqKey)
//			db.mu.Unlock()
//
//			atomic.AddInt64(&db.waitDuration, int64(time.Since(waitStart)))
//
//			select {
//			default:
//			case ret, ok := <-req:
//				// 这边值得注意了，因为现在已经被context取消了。但是刚刚放了自己的号码牌进去排队里面。意思是说不定已经发了连接了，所以得注意归还！
//				if ok && ret.conn != nil {
//					db.putConn(ret.conn, ret.err, false)
//				}
//			}
//			return nil, ctx.Err()
//		case ret, ok := <-req:
//			// 下面是已经获得连接后的操作了。检测一下获得连接的状况。因为有可能已经过期了等等。
//			atomic.AddInt64(&db.waitDuration, int64(time.Since(waitStart)))
//
//			if !ok {
//				return nil, errDBClosed
//			}
//			if ret.err == nil && ret.conn.expired(lifetime) {
//				ret.conn.Close()
//				return nil, driver.ErrBadConn
//			}
//			if ret.conn == nil {
//				return nil, ret.err
//			}
//			ret.conn.Lock()
//			err := ret.conn.lastErr
//			ret.conn.Unlock()
//			if err == driver.ErrBadConn {
//				ret.conn.Close()
//				return nil, driver.ErrBadConn
//			}
//			return ret.conn, ret.err
//		}
//	}
//	// 下面就是如果上面说的限制情况不存在，可以创建先连接时候，要做的创建连接操作了。
//	db.numOpen++ // optimistically
//	db.mu.Unlock()
//	ci, err := db.connector.Connect(ctx)
//	if err != nil {
//		db.mu.Lock()
//		db.numOpen-- // correct for earlier optimism
//		db.maybeOpenNewConnections()
//		db.mu.Unlock()
//		return nil, err
//	}
//	db.mu.Lock()
//	dc := &driverConn{
//		db:        db,
//		createdAt: nowFunc(),
//		ci:        ci,
//		inUse:     true,
//	}
//	db.addDepLocked(dc, dc)
//	db.mu.Unlock()
//	return dc, nil
//}
//复制代码
