package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
)

// 请用本机 MySQL 实现下面这 4 个 test case (MySQL 版本 8.0), 注意操作都必须是存盘操作, 不能全部都在内存中返回数据. 数据一定是要经过了硬盘的
// testCase0 : 实现高性能的大量数据的序列化写入
// testCase1 : 在上述基础上可以读取某批 device 的 TimeZone
// testCase2 : 在上述基础上可以读取某批 device 的具体数据, 并且单个 device 的消息是顺序的
// testCase3 : 删除某批 device 写入的数据
// 请注意标记有 FIXME 的地方是一定需要实现代码的地方, 标记有 TODO 的地方是起到提示作用 . 79 行之前的代码结构都可以修改, 只要代码结构和数据总量不变, 表达的意思是一样的就行
// TODO 评分标准为 : testCase 时间尽可能短(时间长短会决定最终分数), 并且保证正确性
func main() {
	runtime.GOMAXPROCS(4)
	startTime := time.Now()
	globalDb = mustGetMysqlDb()
	defer globalDb.Close()
	err := globalDb.Ping()
	if err != nil {
		panic(err)
	}
	testCase0()
	testCase1()
	testCase2()
	testCase3()
	endTime := time.Now()
	fmt.Println("total delay", endTime.Sub(startTime).String())
}

const (
	TABLE_DEVICE = "device"
	BATCH_SIZE   = 20000
	THREAD_SIZE  = 4
)

// FIXME 请完善这个函数, 函数需要把 genStatisticsDataCb 返回的数据写入到数据库里面, 要求写入时间尽可能的短 (使用 MySQL 8.0实现, 一定要存盘)
func testCase0() {
	// 表存在就先删除，方便反复测试
	_, err := globalDb.Exec(
		"drop table if exists `" + TABLE_DEVICE + "`;",
	)
	if err != nil {
		panic(err)
	}
	_, err = globalDb.Exec("set GLOBAL innodb_flush_log_at_trx_commit=0")
	if err != nil {
		panic(err)
	}
	_, err = globalDb.Exec("SET GLOBAL innodb_buffer_pool_size=100663296")
	if err != nil {
		panic(err)
	}
	_, err = globalDb.Exec("SET GLOBAL sync_binlog=0")
	if err != nil {
		panic(err)
	}

	// 创建数据表
	_, err = globalDb.Exec(
		"CREATE TABLE `test`.`" + TABLE_DEVICE + "` (" +
			//"Id int(10) NOT NULL AUTO_INCREMENT," +  // mysql如果不指定主键有默认的隐藏主键
			"device_id varchar(40)," +
			"event_rand_id varchar(10)," +
			"platform varchar(10)," +
			"client_version_detail varchar(10)," +
			"create_time bigint," +
			"device_time_zone varchar(30)" +
			//"INDEX index_name (device_id, device_time_zone)" + // 索引会影响插入效率，这个看具体场景，目前这里不加
			" );",
	)
	if err != nil {
		panic(err)
	}
	/*_, err = globalDb.Exec("ALTER TABLE `test`.`" + TABLE_DEVICE + "`ADD INDEX idx_d_d(device_id, device_time_zone)")
	if err != nil {
		panic(err)
	}*/

	// 写数据
	platformList := []string{
		"ios",
		"android",
		"mac",
		"amazon",
	}
	const deviceNum = 1 << 18
	startTime := time.Now()
	// FIXME 下面部分都可能会修改到 (不能修改数据量 platformList,deviceNum)
	// FIXME 需要实现把下面的数据写入到 MySQL 8.0 并且存盘, 这个函数执行完之后保证所有数据全部存盘

	baseInsert := "insert into `test`.`" + TABLE_DEVICE +
		"` (device_id,event_rand_id,platform,client_version_detail,create_time," +
		"device_time_zone) VALUES "

	// 每次批量插入500条
	// 还想速度更快可以起4个携程，对应不同platform，通过chan 获取批次数据进行写入，直到数据构造完成chan关闭

	var wg sync.WaitGroup

	ch := make(chan string)

	for _, platform := range platformList {
		wg.Add(1)
		platformLocal := platform
		go func() {
			var sb strings.Builder
			for i := 0; i < deviceNum; i++ {
				deviceId := platformLocal + "_" + RandStringBytesMaskImprSrcUnsafe(32)
				genStatisticsDataCb(
					deviceId, platformLocal, 5, func(msg *StatisticsMessage) {
						// FIXME 这里会有大量数据返回, 需要处理最终存盘 (MySQL 8.0)
						sb.WriteString("(\"" + msg.DeviceId + "\",\"" + msg.EventRandId + "\"," +
							"\"" + msg.Platform + "\",\"" + msg.ClientVersionDetail + "\"," +
							"\"" + strconv.FormatInt(msg.Time.UnixNano(), 10) + "\",\"" + msg.DeviceTimeZone + "\"),")
					},
				)

				if (i+1)%BATCH_SIZE == 0 {
					subSql := sb.String()
					subSql = subSql[0 : len(subSql)-1]
					ch <- subSql
					sb.Reset()
				}
			}

			if sb.Cap() > 0 {
				subSql := sb.String()
				subSql = subSql[0 : len(subSql)-1]
				ch <- subSql
				sb.Reset()
			}
			wg.Done()
		}()
	}

	go func() {
		for {
			subSql, ok := <-ch
			if !ok {
				break
			}
			_, err = globalDb.Exec(
				baseInsert + subSql,
			)
			if err != nil {
				panic(err)
			}
		}
	}()

	wg.Wait()
	close(ch)

	// 刷新数据落盘
	_, err = globalDb.Exec("FLUSH ENGINE LOGS;")
	if err != nil {
		panic(err)
	}
	_, err = globalDb.Exec("FLUSH TABLES;")
	if err != nil {
		panic(err)
	}

	endTime := time.Now()

	// 需要保证下面这个 log 执行时所有数据均写入 MySQL (本地硬盘)
	fmt.Println("delay", "testCase0", endTime.Sub(startTime).String())

	resetConn()
}

// FIXME 通过 DB 获取某批 DeviceId 的 TimeZone, 要求读取时间尽可能的短, 禁止读取内存里面已有的那个 map (MySQL 8.0 从磁盘里面获取数据)
func getDeviceIdTimeTimeZoneMap(deviceIdList []string) map[string]string {
	if len(deviceIdList) == 0 {
		return nil
	}
	ret := make(map[string]string, len(deviceIdList))

	n := len(deviceIdList) / THREAD_SIZE

	ch := make(chan string)
	go func() {
		for {
			device, ok := <-ch
			if !ok {
				return
			}
			s := strings.Split(device, ":")
			_, exist := ret[s[0]]
			if exist {
				continue
			}
			ret[s[0]] = s[1]
		}
	}()

	var wg sync.WaitGroup
	wg.Add(THREAD_SIZE)
	// 启动与核心数一致的协程避免多次申请协程，虽然go GPM本身协程有复用的情况，但是不可控
	for r := 0; r < THREAD_SIZE; r++ {
		go func(seg int) {
			sb := ""
			if seg == THREAD_SIZE-1 {
				sb = buildWhere(deviceIdList[seg*n:])
			} else {
				sb = buildWhere(deviceIdList[seg*n : (seg+1)*n+1])
			}
			result, err := globalDb.Query("SELECT DISTINCT device_id, device_time_zone FROM `test`.`" +
				TABLE_DEVICE + "` WHERE device_id IN " + sb)
			if err != nil {
				panic(err)
			}

			for result.Next() {
				deviceId, deviceTime := "", ""
				err = result.Scan(&deviceId, &deviceTime)
				if err != nil {
					panic(err)
				}
				ch <- deviceId + ":" + deviceTime
			}

			err = result.Close()
			if err != nil {
				panic(err)
			}

			wg.Done()
		}(r)
	}

	wg.Wait()
	close(ch)

	resetConn()

	return ret
}

func resetConn() {
	//大查询后断开连接释放mysql用到的临时内存
	err := globalDb.Close()
	if err != nil {
		panic(err)
	}
	// 重新打开一个新的连接
	globalDb = mustGetMysqlDb()
}

func buildWhere(deviceIdList []string) string {
	var sb strings.Builder
	sb.WriteByte('(')
	sb.WriteByte('"')
	sb.WriteString(deviceIdList[0])
	sb.WriteByte('"')

	for idx := 1; idx < len(deviceIdList); idx++ {
		sb.WriteByte(',')
		sb.WriteByte(' ')
		sb.WriteByte('"')
		sb.WriteString(deviceIdList[idx])
		sb.WriteByte('"')
		sb.WriteByte(' ')
	}

	sb.WriteByte(')')
	return sb.String()
}

// FIXME 通过 DB 获取某批 DeviceId 的 StatisticsMessage List 的 map,并且内部的 List 是顺序的, 要求读取时间尽可能的短 (MySQL 8.0 从磁盘里面获取数据)
func getDeviceIdStatisticDataListMap(deviceIdList []string) (deviceIdMsgMap map[string][]*StatisticsMessage) {
	// 此处传入的也可能是其他已知的 deviceId
	if len(deviceIdList) == 0 {
		return nil
	}
	ret := make(map[string][]*StatisticsMessage, len(deviceIdList))

	sb := buildWhere(deviceIdList)
	result, err := globalDb.Query("SELECT device_id, event_rand_id, platform, client_version_detail, create_time, device_time_zone FROM `test`.`" + TABLE_DEVICE + "` WHERE device_id IN " + sb + " ORDER BY device_id, event_rand_id, create_time")
	if err != nil {
		panic(err)
	}
	defer result.Close()

	for result.Next() {
		sm := &StatisticsMessage{}
		var timestam int64
		err = result.Scan(&sm.DeviceId, &sm.EventRandId, &sm.Platform, &sm.ClientVersionDetail, &timestam, &sm.DeviceTimeZone)
		if err != nil {
			panic(err)
		}

		sm.Time = time.Unix(timestam/1e9, timestam%1e9)
		ret[sm.DeviceId] = append(ret[sm.DeviceId], sm)
	}
	resetConn()
	return ret
}

// FIXME 通过 DB 删除某批 DeviceId 的 StatisticsMessage (MySQL 8.0 从磁盘里面删除)
func deleteDeviceIdStatisticData(deviceIdList []string) {
	if len(deviceIdList) == 0 {
		return
	}

	// 如果开启并行删除需要考虑到目前我们没有对应行的标识，会有人为锁冲突的情况，这个时候串行才是有效的删除策略
	// 这里可根据不同db设置，如tidb有5w的限制，如果db没有限制最好是给一个结合业务相对较小的值，避免大事务
	n := 1 << 10

	subSql := buildWhere(deviceIdList)

	// 启动与核心数一致的协程避免多次申请协程，虽然go GPM本身协程有复用的情况，但是不可控
	for {
		rsp, err := globalDb.Exec("DELETE FROM `test`.`" + TABLE_DEVICE + "` WHERE device_id IN " + subSql + " LIMIT " + strconv.Itoa(n))
		if err != nil {
			panic(err)
		}
		affected, err := rsp.RowsAffected()
		if err != nil {
			panic(err)
		}
		if int(affected) < n {
			break
		}
		fmt.Println("delete num: ", affected)
		_, err = globalDb.Exec("commit")
		if err != nil {
			panic(err)
		}
	}

	// 复原, 事务提交，也可以不用，最后通过重建表清空磁盘数据，防止磁盘空洞占用磁盘
	_, err := globalDb.Exec("alter table `test`.`" + TABLE_DEVICE + "` ENGINE=InnoDB, ALGORITHM=INPLACE, LOCK=NONE")
	if err != nil {
		panic(err)
	}

}

/* ================================================================================== 代码请写到上面,重要的事情说 3 遍 ================================================================================== */
/* ================================================================================== 代码请写到上面,重要的事情说 3 遍 ================================================================================== */
/* ================================================================================== 代码请写到上面,重要的事情说 3 遍 ================================================================================== */
/* ================================================================================== 下面都是不用修改的代码,重要的事情说 3 遍 ================================================================================== */
/* ================================================================================== 下面都是不用修改的代码,重要的事情说 3 遍 ================================================================================== */
/* ================================================================================== 下面都是不用修改的代码,重要的事情说 3 遍 ================================================================================== */

// ////////////////////////////////////////////
// /////// 接下来为测试的检测代码,不需要修改 ////////
// ////////////////////////////////////////////

// TODO 注意后面 3 个 testCase 都不需要修改, 需要保证响应速度和正确性, 检测磁盘返回的数据是否正常
// TODO testCase1 的 TODO 都不用写代码, 只是正确性检测(MySQL 8.0 从磁盘读取出来的数据)
func testCase1() {
	startTime := time.Now()
	deviceIdList := getRandHalfDeviceIdList()
	m := getDeviceIdTimeTimeZoneMap(deviceIdList)
	if len(m) != len(deviceIdList) { // TODO 检测数据是否正常
		panic("tz data len error")
	}
	originDeviceIdTimeZoneMap := getDeviceIdTimeZoneMap()
	for deviceId, timeZone := range m {
		if originDeviceIdTimeZoneMap[deviceId] != timeZone { // TODO 检测数据是否正常
			panic("check deviceId timeZone failed")
		}
	}
	endTime := time.Now()
	fmt.Println("delay", "testCase1", endTime.Sub(startTime).String())
}

// TODO testCase2 的 TODO 都不用写代码, 只是正确性检测(MySQL 8.0 从磁盘读取出来的数据)
func testCase2() {
	startTime := time.Now()
	// 基于 testCase0 , 把写入的某批设备的 StatisticsMessageList 取出来
	deviceIdList := getRandHalfDeviceIdList()
	m := getDeviceIdStatisticDataListMap(deviceIdList)
	if len(m) != len(deviceIdList) { // TODO 检测数据是否正常
		panic("map data len error")
	}
	var tagDeviceId string
	for k := range m {
		tagDeviceId = k
		break
	}
	list := m[tagDeviceId]
	if len(list) != 5 { // TODO 检测数据是否正常
		panic("data len error")
	}
	getLastInt := func(s string) int {
		lastValue := s[len(s)-1:]
		out, err := strconv.Atoi(lastValue)
		if err != nil {
			panic(err)
		}
		return out
	}
	originDeviceIdTimeZoneMap := getDeviceIdTimeZoneMap()
	var lastMsg *StatisticsMessage
	for i, msg := range list {
		if msg.DeviceId != tagDeviceId { // TODO 检测数据是否正常
			panic("deviceId error")
		}
		if msg.DeviceTimeZone != originDeviceIdTimeZoneMap[msg.DeviceId] { // TODO 检测数据是否正常
			panic("DeviceTimeZone error")
		}
		if msg.ClientVersionDetail != "1.1.1" { // TODO 检测数据是否正常
			panic("ClientVersionDetail error")
		}
		if msg.Platform != strings.Split(msg.DeviceId, "_")[0] { // TODO 检测数据是否正常
			panic("Platform error")
		}
		if i > 0 {
			if msg.Time.Sub(lastMsg.Time) != 1 { // TODO 检测数据是否正常
				panic("data Time sort error")
			}
			if getLastInt(msg.EventRandId)-getLastInt(lastMsg.EventRandId) != 1 { // TODO 检测数据是否正常
				panic("data EventRandId sort error")
			}
		}
		lastMsg = msg
	}
	endTime := time.Now()
	fmt.Println("delay", "testCase2", endTime.Sub(startTime).String())
}

// TODO testCase3 的 TODO 都不用写代码, 只是正确性检测(MySQL 8.0 从磁盘读取出来的数据)
func testCase3() {
	startTime := time.Now()
	// TODO 把 testCase0 写入的数全部删除, 需要在磁盘上全部抹掉
	deviceIdList := getRandHalfDeviceIdList()
	deleteDeviceIdStatisticData(deviceIdList)
	list := getDeviceIdStatisticDataListMap(deviceIdList)
	if len(list) != 0 { // TODO 检测数据是否从磁盘上删除
		panic("not delete datasource")
	}
	endTime := time.Now()
	fmt.Println("delay", "testCase3", endTime.Sub(startTime).String())
}

// ////////////////////////////////////////////
// /////// 接下来的代码为模版代码，请勿修改 /////////
// ////////////////////////////////////////////

var globalDb *sql.DB

func mustGetMysqlDb() *sql.DB {
	// 连接数据库
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/?parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}
	databaseName := "test"
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("USE " + databaseName)
	if err != nil {
		panic(err)
	}
	return db
}

type StatisticsMessage struct {
	Time                time.Time `json:",omitempty"`
	EventRandId         string    `json:",omitempty"`
	DeviceId            string    `json:",omitempty"`
	DeviceTimeZone      string    `json:",omitempty"`
	ClientVersionDetail string    `json:",omitempty"`
	Platform            string    `json:",omitempty"`
}

var timeZoneList = []string{
	"Asia/Shanghai",
	"Asia/Urumqi",
	"Asia/Hong_Kong",
	"Asia/Taipei",
	"Asia/Singapore",
	"Asia/Qatar",
	"America/Cayman",
	"America/New_York",
	"America/Sao_Paulo",
}

var gTotalDeviceIdList []string
var gTotalDeviceIdListLock sync.Mutex
var gDeviceIdTimeZoneMap = map[string]string{}
var gDeviceIdTimeZoneMapLock sync.Mutex

func totalDeviceIdAddOne(deviceId string) {
	gTotalDeviceIdListLock.Lock()
	gTotalDeviceIdList = append(gTotalDeviceIdList, deviceId)
	gTotalDeviceIdListLock.Unlock()
}

// TODO 随机获取一半的数据
func getRandHalfDeviceIdList() []string {
	gTotalDeviceIdListLock.Lock()
	defer gTotalDeviceIdListLock.Unlock()
	srcDeviceIdList := gTotalDeviceIdList
	dest := make([]string, len(srcDeviceIdList))
	perm := rand.Perm(len(srcDeviceIdList))
	for i, v := range perm {
		dest[v] = srcDeviceIdList[i]
	}
	if len(dest) <= 1 {
		return dest
	}
	return dest[:len(dest)/2]
}

func deviceIdSetTimeZone(deviceId, timeZone string) {
	gDeviceIdTimeZoneMapLock.Lock()
	gDeviceIdTimeZoneMap[deviceId] = timeZone
	gDeviceIdTimeZoneMapLock.Unlock()
}

// TODO 此处只是用来做正确性的检查, 笔试禁止使用这个 map 来读取数据
func getDeviceIdTimeZoneMap() map[string]string {
	gDeviceIdTimeZoneMapLock.Lock()
	out := gDeviceIdTimeZoneMap
	gDeviceIdTimeZoneMapLock.Unlock()
	return out
}

func genStatisticsDataCb(deviceId, platform string, count int, f func(msg *StatisticsMessage)) {
	totalDeviceIdAddOne(deviceId)
	randId := RandStringBytesMaskImprSrcUnsafe(6)
	timeZone := timeZoneList[int(srcInt63())%len(timeZoneList)]
	deviceIdSetTimeZone(deviceId, timeZone)
	t := time.Now()
	for i := 0; i < count; i++ {
		msg := newStatisticsMessage()
		msg.Time = t.Add(time.Duration(i))
		msg.EventRandId = randId + strconv.Itoa(i)
		msg.DeviceId = deviceId
		msg.DeviceTimeZone = timeZone
		msg.ClientVersionDetail = "1.1.1"
		msg.Platform = platform
		f(msg)
		freeStatisticsMessage(msg)
	}
}

var structPool sync.Pool

func newStatisticsMessage() *StatisticsMessage {
	msg := structPool.Get()
	if msg != nil {
		return msg.(*StatisticsMessage)
	}
	return &StatisticsMessage{}
}

func freeStatisticsMessage(msg *StatisticsMessage) {
	msg.reset()
	structPool.Put(msg)
}

func (msg *StatisticsMessage) reset() {
	msg.Time = time.Time{}
	msg.EventRandId = ""
	msg.DeviceId = ""
	msg.DeviceTimeZone = ""
	msg.ClientVersionDetail = ""
	msg.Platform = ""
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())
var randLock sync.Mutex

func srcInt63() int64 {
	randLock.Lock()
	out := src.Int63()
	randLock.Unlock()
	return out
}

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, srcInt63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = srcInt63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
