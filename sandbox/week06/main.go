package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	ioutil "io" // "io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/go-sql-driver/mysql" // database/sql implementation
	"github.com/gomodule/redigo/redis" // "github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func main() {
	// mysqlSimple()
	// gormCRUD()
	// sqlInjection()
	// memcacheSimple()
	// taggedMemCache()
	// redisSession()
	rabbitPicResize()
}

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func rabbitPicResize() {
	show("rabbitPicResize: program started ...")
	show("HINT: worker program should be running by now: `go run rabbit_pic_resize_worker/worker.go&`")

	flag.Parse()
	var err error
	rabbitConn, err = amqp.Dial(*rabbitAddr)
	panicOnError("can't connect to rabbit", err)

	rabbitChan, err = rabbitConn.Channel()
	panicOnError("can't open AMQP channel", err)
	defer rabbitChan.Close()

	q, err := rabbitChan.QueueDeclare(
		ImageResizeQueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	panicOnError("can't init queue", err)
	fmt.Printf("queue %s have %d msg and %d consumers\n", q.Name, q.Messages, q.Consumers)

	http.HandleFunc("/", imageResizeMainPage)
	http.HandleFunc("/upload", imageResizeUploadPage)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

const (
	ImageResizeQueueName     = "image_resize"
	ImageResizeStoragePrefix = "./images" // it's a no-no
)

var (
	rabbitAddr = flag.String("addr", "amqp://guest:guest@localhost:5672/", "rabbit addr")
	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel
)

type ImgResizeTask struct {
	Name string
	MD5  string
}

func imageResizeMainPage(w http.ResponseWriter, r *http.Request) {
	var uploadFormTmpl = []byte(`
<html>
	<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload">
	</form>
	</body>
</html>
`)

	w.Write(uploadFormTmpl)
}
func imageResizeUploadPage(w http.ResponseWriter, r *http.Request) {
	uploadData, handler, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer uploadData.Close()

	fmt.Fprintf(w, "handler.Filename %v\n", handler.Filename)
	fmt.Fprintf(w, "handler.Header %#v\n", handler.Header)

	tmpName := RandStringRunes(32)
	tmpFile := ImageResizeStoragePrefix + "/" + tmpName + ".jpg"
	fileRef, err := os.Create(tmpFile) // possible drop of existing file
	if err != nil {
		http.Error(w, "can't create file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hasher := md5.New()
	writtenBytesCnt, err := ioutil.Copy(fileRef, ioutil.TeeReader(uploadData, hasher))
	if err != nil {
		http.Error(w, "can't save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fileRef.Sync()
	fileRef.Close()

	md5Sum := hex.EncodeToString(hasher.Sum(nil))
	realFile := ImageResizeStoragePrefix + "/" + md5Sum + ".jpg"
	err = os.Rename(tmpFile, realFile)
	if err != nil {
		http.Error(w, "can't rename file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	msgBytes, _ := json.Marshal(ImgResizeTask{handler.Filename, md5Sum}) // skip error

	fmt.Println("sending AMQP message ", string(msgBytes))
	err = rabbitChan.Publish(
		"",                   // exchange
		ImageResizeQueueName, // routing key
		false,                // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgBytes,
		})
	panicOnError("can't publish task", err)

	fmt.Fprintf(w, "Upload %d bytes successful\n", writtenBytesCnt)
}

func redisSession() {
	show("redisSession: program started ...")

	var redisAddr = flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")
	flag.Parse()

	var err error
	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		log.Fatalf("can't connect to redis")
	}

	redisSessManager = NewRedisDemoSessionManager(redisConn)

	http.HandleFunc("/", redisDemoInnerPage)
	http.HandleFunc("/login", redisDemoLoginPage)
	http.HandleFunc("/logout", redisDemoLogoutPage)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

var (
	redisSessManager  *RedisDemoSessionManager
	usersDB_redisDemo = map[string]string{
		"foo": "bar",
		"baz": "todolo",
	}
)

func redisDemoInnerPage(w http.ResponseWriter, r *http.Request) {
	var loginFormTmpl = []byte(`
	<html>
		<body>
		<form action="/login" method="post">
			Login: <input type="text" name="login">
			Password: <input type="password" name="password">
			<input type="submit" value="Login">
		</form>
		</body>
	</html>
	`)

	sess, err := redisCheckSession(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.Write(loginFormTmpl)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Welcome, "+sess.Login+" <br />")
	fmt.Fprintln(w, "Session ua: "+sess.Useragent+" <br />")
	fmt.Fprintln(w, `<a href="/logout">logout</a>`)
}

func redisDemoLoginPage(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue("login")
	inputPass := r.FormValue("password")
	expiration := time.Now().Add(24 * time.Hour)

	passwd, exist := usersDB_redisDemo[inputLogin]
	if !exist || passwd != inputPass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sess, err := redisSessManager.SaveSession(&RedisDemoSession{
		Login:     inputLogin,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		log.Println("can't create session:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func redisDemoLogoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redisSessManager.DropSession(&RedisDemoSessionID{
		ID: session.Value,
	})

	session.Expires = time.Now().AddDate(0, 0, -1)

	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}

func redisCheckSession(r *http.Request) (*RedisDemoSession, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sess := redisSessManager.LoadSession(&RedisDemoSessionID{
		ID: cookieSessionID.Value,
	})
	return sess, nil
}

type RedisDemoSessionID struct {
	ID string
}
type RedisDemoSession struct {
	Login     string
	Useragent string
}
type RedisDemoSessionManager struct {
	redisConn redis.Conn
}

func NewRedisDemoSessionManager(conn redis.Conn) *RedisDemoSessionManager {
	return &RedisDemoSessionManager{
		redisConn: conn,
	}
}

func (sm *RedisDemoSessionManager) SaveSession(sess *RedisDemoSession) (*RedisDemoSessionID, error) {
	const sessKeyLen = 10

	id := RedisDemoSessionID{ID: RandStringRunes(sessKeyLen)}

	sessBytes, _ := json.Marshal(sess)
	mkey := "sessions:" + id.ID
	reply, err := sm.redisConn.Do("SET", mkey, sessBytes, "EX", 86400)
	result, err := redis.String(reply, err)
	if err != nil {
		return nil, err
	}
	if result != "OK" {
		return nil, fmt.Errorf("result not OK")
	}

	return &id, nil
}

func (sm *RedisDemoSessionManager) LoadSession(in *RedisDemoSessionID) *RedisDemoSession {
	mkey := "sessions:" + in.ID
	sessBytes, err := redis.Bytes(sm.redisConn.Do("GET", mkey))
	if err != nil {
		log.Println("can't get data:", err)
		return nil
	}

	sess := &RedisDemoSession{}
	err = json.Unmarshal(sessBytes, sess)
	if err != nil {
		log.Println("can't unpack session data:", err)
		return nil
	}

	return sess
}

func (sm *RedisDemoSessionManager) DropSession(in *RedisDemoSessionID) {
	mkey := "sessions:" + in.ID
	_, err := redis.Int(sm.redisConn.Do("DEL", mkey))
	if err != nil {
		log.Println("redis error:", err)
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func taggedMemCache() {
	show("taggedMemCache: program started ...")

	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	tc := &TaggedCacheWrapper{memcacheClient} // just embedded MC, with custom methods

	mkey := "habrposts"
	tc.Delete(mkey)

	// fetch real data from remote server: rss, tags
	fetchDataFromRemoteSrv := func() (interface{}, []string, error) {
		rssRef, err := GetHabrPosts()
		if err != nil {
			return nil, nil, err
		}
		return rssRef, []string{"habrTag", "geektimes"}, nil
	}

	// lets ask for cached data a few times ...

	fmt.Println("\nTGet call #1") // cache empty, fetch data, put in cache
	posts := RSS_TaggedCache{}
	err := tc.TGet(mkey, 30, &posts, fetchDataFromRemoteSrv) // key, ttl, buf, data_getter
	fmt.Println("#1", len(posts.Items), "err:", err)
	/*
		TGet call #1
		Record not found in memcache
		fetching https://habrahabr.ru/rss/best/
		#1 40 err: <nil>
	*/

	fmt.Println("\nTGet call #2") // get from cache
	posts = RSS_TaggedCache{}
	err = tc.TGet(mkey, 30, &posts, fetchDataFromRemoteSrv)
	fmt.Println("#2", len(posts.Items), "err:", err)
	/*
		TGet call #2
		#2 40 err: <nil>
	*/

	fmt.Println("\ninc tag habrTag") // current version (of habrTag) updated
	tc.Increment("habrTag", 1)       // you should consider data under tag expired, fetch again

	go func() { // fetch data async: n.b. how two requests intertwine
		// time.Sleep(time.Millisecond)
		fmt.Println("\nTGet call #async")
		posts = RSS_TaggedCache{}
		err = tc.TGet(mkey, 30, &posts, fetchDataFromRemoteSrv)
		fmt.Println("#async", len(posts.Items), "err:", err)
		/*
			TGet call #async
			fetching https://habrahabr.ru/rss/best/
			get lock try 0
			get lock try 1
			get lock try 2
			get lock try 3
			fetching https://habrahabr.ru/rss/best/
			#3 40 err: <nil>
			#async 40 err: <nil>
		*/
	}()

	fmt.Println("\nTGet call #3")
	posts = RSS_TaggedCache{}
	err = tc.TGet(mkey, 30, &posts, fetchDataFromRemoteSrv)
	fmt.Println("#3", len(posts.Items), "err:", err)
	/*
		TGet call #3
		TGet call #async
		fetching https://habrahabr.ru/rss/best/
		get lock try 0
		get lock try 1
		get lock try 2
		get lock try 3
		fetching https://habrahabr.ru/rss/best/
		#3 40 err: <nil>
	*/

	time.Sleep(1000 * time.Millisecond)
	show("end of program. ", err)
}

type RSS_TaggedCache struct {
	Items []RSSItem_TaggedCache `xml:"channel>item"`
}
type RSSItem_TaggedCache struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

// GetHabrPosts fetch data from remote server
func GetHabrPosts() (*RSS_TaggedCache, error) {
	fmt.Println("fetching https://habrahabr.ru/rss/best/")
	resp, err := http.Get("https://habrahabr.ru/rss/best/")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	rss := new(RSS_TaggedCache)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

// TaggedCacheRebuildFuncT is a type for GetHabrPosts wrapper
type TaggedCacheRebuildFuncT func() (interface{}, []string, error)

type TaggedCacheItemRaw struct { // first layer of marshalling
	Data json.RawMessage // no unmarshal before tags are checked
	Tags map[string]int  // key:version
}
type TaggedCacheItem struct { // second layer of marshalling
	Data interface{}
	Tags map[string]int
}

type TaggedCacheWrapper struct {
	*memcache.Client
}

// TGet is main method for cached data fetching, wrapper for get, check, fetch-from-remote, store
func (tc *TaggedCacheWrapper) TGet(
	mkey string,
	ttl int32,
	buf interface{},
	rebuildFunc TaggedCacheRebuildFuncT,
) (err error) {
	inKind := reflect.ValueOf(buf).Kind()
	if inKind != reflect.Ptr {
		return fmt.Errorf("in must be ptr, got %s", inKind)
	}

	tc.checkLock(mkey)           // we are dictributed, check if key is un-locked (should set read-lock here)
	itemRaw, err := tc.Get(mkey) // and here we see possible bug: imagine that key updated in this moment
	if err == memcache.ErrCacheMiss {
		fmt.Println("Record not found in memcache")
		return tc.rebuild(mkey, ttl, buf, rebuildFunc) // fetch, store
	} else if err != nil {
		return err
	}

	item := &TaggedCacheItemRaw{} // don't touch data, only tags
	err = json.Unmarshal(itemRaw.Value, &item)
	if err != nil {
		return err
	}

	tagsValid, err := tc.isTagsValid(item.Tags)
	if err != nil {
		return fmt.Errorf("isTagsValid error %s", err)
	}
	if tagsValid {
		err = json.Unmarshal(item.Data, &buf) // now we need data
		return err                            // data in buf or err
	} else { // if tags are invalid:
		return tc.rebuild(mkey, ttl, buf, rebuildFunc)
	}
}

func (tc *TaggedCacheWrapper) rebuild(
	mkey string,
	ttl int32,
	buf interface{},
	rebuildFunc TaggedCacheRebuildFuncT,
) error {
	tc.lockRebuild(mkey) // should lock just before update (later, not now)
	defer tc.unlockRebuild(mkey)

	result, tags, err := rebuildFunc()

	// ожидаем и возвращаем одинаковые типы
	if reflect.TypeOf(result) != reflect.TypeOf(buf) {
		return fmt.Errorf(
			"data type mismatch, expected %s, got %s", reflect.TypeOf(buf),
			reflect.TypeOf(result),
		)
	}

	currTags, err := tc.getCurrentItemTags(tags, ttl)
	if err != nil {
		return err
	}

	cacheData := TaggedCacheItem{result, currTags}
	rawData, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}

	err = tc.Set(&memcache.Item{
		Key:        mkey,
		Value:      rawData,
		Expiration: int32(ttl),
	})

	// write data to reference under interface wrapper
	resultValRef := reflect.ValueOf(result)
	resultVal := reflect.Indirect(resultValRef)
	bufValRef := reflect.ValueOf(buf)
	bufVal := reflect.Indirect(bufValRef)

	bufVal.Set(resultVal)
	return nil
}

// isTagsValid compare given tags versions with stored in cache
func (tc *TaggedCacheWrapper) isTagsValid(itemTags map[string]int) (bool, error) {
	cacheKeys := make([]string, 0, len(itemTags))
	for tagKey := range itemTags {
		cacheKeys = append(cacheKeys, tagKey)
	}

	cachedTags, err := tc.GetMulti(cacheKeys)
	if err != nil {
		return false, err
	}

	currentTagsMap := make(map[string]int, len(cachedTags))
	for tagKey, tagItem := range cachedTags {
		tagVersion, err := strconv.Atoi(string(tagItem.Value))
		if err != nil {
			return false, err
		}

		currentTagsMap[tagKey] = tagVersion
	}

	return reflect.DeepEqual(itemTags, currentTagsMap), nil
}

// getCurrentItemTags read/write tags from/to cache: if no tag found (expired?): set a new tag with current time as value
func (tc *TaggedCacheWrapper) getCurrentItemTags(tags []string, ttl int32) (map[string]int, error) {
	currTags, err := tc.GetMulti(tags)
	if err != nil {
		return nil, err
	}

	resultTags := make(map[string]int, len(tags))
	now := int(time.Now().Unix())
	nowBytes := []byte(fmt.Sprint(now))

	for _, tagKey := range tags {
		tagItem, tagExist := currTags[tagKey]
		if !tagExist { // set new tag value
			err := tc.Set(&memcache.Item{
				Key:        tagKey,
				Value:      nowBytes,
				Expiration: int32(ttl),
			})
			if err != nil {
				return nil, err
			}
			resultTags[tagKey] = now
		} else { // got tag value
			i, err := strconv.Atoi(string(tagItem.Value))
			if err != nil {
				return nil, err
			}
			resultTags[tagKey] = i
		}
	}

	return resultTags, nil
}

func (tc *TaggedCacheWrapper) checkLock(mkey string) error {
	for i := 0; i < 4; i++ {
		_, err := tc.Get("lock_" + mkey)
		if err == memcache.ErrCacheMiss {
			return nil // no lock
		}
		if err != nil {
			return err // error
		}
		// lock exist, wait
		time.Sleep(10 * time.Millisecond)
	}
	return nil // timeout, should return error
}

func (tc *TaggedCacheWrapper) unlockRebuild(mkey string) {
	tc.Delete("lock_" + mkey)
}
func (tc *TaggedCacheWrapper) lockRebuild(mkey string) (bool, error) {
	// пытаемся взять лок на перестроение кеша // чтобы все не ломанулись его перестраивать

	// параметры надо тюнить
	lockKey := "lock_" + mkey
	lockAccuired := false
	for i := 0; i < 4; i++ {
		// add добавляет запись если её ещё нету
		err := tc.Add(&memcache.Item{
			Key:        lockKey,
			Value:      []byte("1"),
			Expiration: int32(3),
		})
		if err == memcache.ErrNotStored {
			fmt.Println("get lock, try #", i)
			time.Sleep(time.Millisecond * 10)
			continue
		} else if err != nil {
			return false, err
		}

		lockAccuired = true
		break
	}
	if !lockAccuired {
		return false, fmt.Errorf("Can't get lock")
	}
	return true, nil
}

func memcacheSimple() {
	show("memcacheSimple: program started ...")

	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	mkey := "coursera"

	err := memcacheClient.Set(&memcache.Item{
		Key:        mkey,
		Value:      []byte("1"),
		Expiration: 3, // seconds
	})
	show("mc.Set item under key: ", mkey, err)

	newVal, err := memcacheClient.Increment("habrTag", 1)
	show("mc.Increment by 1 under tag habrTag; new value: ", newVal, err)

	item, err := memcacheClient.Get(mkey)
	if err != nil && err != memcache.ErrCacheMiss {
		show("MC error ", err)
	}
	show("mc.Get value under tag", item, mkey)

	err = memcacheClient.Delete(mkey)
	show("mc.Delete key ", mkey)

	item, err = memcacheClient.Get(mkey)
	if err == memcache.ErrCacheMiss {
		show("record not found in MC, key ", mkey)
	} else {
		show("mc.Get ", item, mkey)
	}

	show("end of program. ", err)
	/*
		2024-04-30T09:02:39.945Z: memcacheSimple: program started ...
		2024-04-30T09:02:39.948Z: mc.Set item under key: string(coursera); <nil>(<nil>);
		2024-04-30T09:02:39.949Z: mc.Increment by 1 under tag habrTag; new value: uint64(0); *errors.errorString(memcache: cache miss);
		2024-04-30T09:02:39.950Z: mc.Get value under tag*memcache.Item(&{coursera [49] 0 0 1}); string(coursera);
		2024-04-30T09:02:39.950Z: mc.Delete key string(coursera);
		2024-04-30T09:02:39.951Z: record not found in MC, key string(coursera);
		2024-04-30T09:02:39.951Z: end of program. *errors.errorString(memcache: cache miss);
	*/
}

func sqlInjection() {
	show("sqlInjection: program started ...")
	/*
		-- setup_db.sql
		DROP TABLE IF EXISTS `users`;
		CREATE TABLE `users` (
		  `id` int(11) NOT NULL,
		  `login` varchar(200) NOT NULL,
		  `name` varchar(200) NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

		INSERT INTO `users` (`id`, `login`, `name`) VALUES
		(1,	'user',	''),
		(2,	'admin',	'');
	*/
	var loginFormTmpl = `
	<html>
		<body>
		<form action="/login" method="post">
			Login: <input type="text" name="login">
			Password: <input type="password" name="password">
			<input type="submit" value="Login">
		</form>
		</body>
	</html>
	`

	// основные настройки к базе
	dsn := "root@tcp(localhost:3306)/coursera?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prapared statements // параметры подставляются сразу
	dsn += "&interpolateParams=true"

	var err error
	// создаём структуру базы // но соединение происходит только при первом запросе
	db, err := sql.Open("mysql", dsn)
	__err_panic(err)
	err = db.Ping() // вот тут будет первое подключение к базе
	__err_panic(err)
	show("connected to DB")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(loginFormTmpl))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var (
			id          int
			login, body string
		)

		inputLogin := r.FormValue("login")
		body += fmt.Sprintln("inputLogin:", inputLogin)

		// ПЛОХО! НЕ ДЕЛАЙТЕ ТАК! // параметры не экранированы должным образом // мы подставляем в запрос параметр как есть
		query := fmt.Sprintf("SELECT id, login FROM users WHERE login = '%s' LIMIT 1", inputLogin)
		// try this input: 404' or login = 'admin
		body += fmt.Sprintln("Sprint query:", query)

		row := db.QueryRow(query)
		err := row.Scan(&id, &login)
		if err == sql.ErrNoRows {
			body += fmt.Sprintln("Sprint case: NOT FOUND")
		} else {
			__err_panic(err)
			body += fmt.Sprintln("Sprint case: FOUND id:", id, "login:", login)
		}

		// ПРАВИЛЬНО // Мы используем плейсхолдеры, там параметры будет экранирован должным образом
		row = db.QueryRow("SELECT id, login FROM users WHERE login = ? LIMIT 1", inputLogin)
		err = row.Scan(&id, &login)
		if err == sql.ErrNoRows {
			body += fmt.Sprintln("Placeholders case: NOT FOUND")
		} else {
			__err_panic(err)
			body += fmt.Sprintln("Placeholders case: FOUND id:", id, "login:", login)
		}

		w.Write([]byte(body))
	})

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func gormCRUD() {
	show("gormCRUD: program started ...")

	// основные настройки к базе
	dsn := "root@tcp(localhost:3306)/coursera?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prapared statements // параметры подставляются сразу
	dsn += "&interpolateParams=true"

	db, err := gorm.Open("mysql", dsn)
	db.DB()
	db.DB().Ping()
	__err_panic(err)
	// defer db.Close() // have no effect?
	show("connected to DB")

	srv := &GormSimpleHttpHandlers{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("./week06/gorm_templates/*")),
	}
	show("loaded templates")

	// в целях упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/", srv.List).Methods("GET")
	r.HandleFunc("/items", srv.List).Methods("GET")
	r.HandleFunc("/items/new", srv.ShowCreateForm).Methods("GET")
	r.HandleFunc("/items/new", srv.Create).Methods("POST")
	r.HandleFunc("/items/{id}", srv.ShowUpdateForm).Methods("GET")
	r.HandleFunc("/items/{id}", srv.Update).Methods("POST")
	r.HandleFunc("/items/{id}", srv.Delete).Methods("DELETE")

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, r)
	show("end of program. ", err)
}

type GormSimplePostItem struct { // added tags, nice
	Id          int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title       string
	Description string
	Updated     string `sql:"null"`
}

func (i *GormSimplePostItem) TableName() string { // gorm hook
	return "items"
}

func (i *GormSimplePostItem) BeforeSave() (err error) { // gorm hook
	fmt.Println("trigger on before save")
	return
}

type GormSimpleHttpHandlers struct {
	DB   *gorm.DB
	Tmpl *template.Template
}

func (h *GormSimpleHttpHandlers) List(w http.ResponseWriter, r *http.Request) {
	items := []*GormSimplePostItem{} // slice of references

	db := h.DB.Find(&items)
	err := db.Error
	__err_panic(err)

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*GormSimplePostItem
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GormSimpleHttpHandlers) ShowCreateForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GormSimpleHttpHandlers) Create(w http.ResponseWriter, r *http.Request) {
	newItem := &GormSimplePostItem{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	db := h.DB.Create(&newItem)
	err := db.Error
	__err_panic(err)
	affected := db.RowsAffected

	fmt.Println("Insert: RowsAffected", affected, "LastInsertId: ", newItem.Id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *GormSimpleHttpHandlers) ShowUpdateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	post := &GormSimplePostItem{}

	db := h.DB.Find(post, id)
	err = db.Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found", id)
	} else {
		__err_panic(err)
	}

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GormSimpleHttpHandlers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	post := &GormSimplePostItem{}
	h.DB.Find(post, id)

	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Updated = "by gorm"

	db := h.DB.Save(post)
	err = db.Error
	__err_panic(err)
	affected := db.RowsAffected

	fmt.Println("Update: RowsAffected", affected)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *GormSimpleHttpHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	db := h.DB.Delete(&GormSimplePostItem{Id: id})
	err = db.Error
	__err_panic(err)
	affected := db.RowsAffected

	fmt.Println("Delete: RowsAffected", affected)

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

func mysqlSimple() {
	/*
		-- items.sql // sandbox/week06/mysql_items.sql
		SET NAMES utf8;
		SET time_zone = '+00:00';
		SET foreign_key_checks = 0;
		SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

		DROP TABLE IF EXISTS `items`;
		CREATE TABLE `items` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `title` varchar(255) NOT NULL,
		  `description` text NOT NULL,
		  `updated` varchar(255) DEFAULT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

		INSERT INTO `items` (`id`, `title`, `description`, `updated`) VALUES
		(1,	'database/sql',	'Рассказать про базы данных',	'foo'),
		(2,	'memcache',	'Рассказать про мемкеш с примером использования',	NULL);
	*/

	show("mysqlItems: program started ...")

	// основные настройки к базе
	dsn := "root@tcp(localhost:3306)/coursera?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prapared statements // параметры подставляются сразу
	dsn += "&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	err = db.Ping() // вот тут будет первое подключение к базе
	__err_panic(err)
	show("connected to DB ", db)

	tmpl := template.Must(template.ParseGlob("./week06/crud_templates/*")) // sandbox\week06\crud_templates\
	show("loaded templates ", tmpl)

	srv := &MysqlSimpleHttpHandlers{
		DB:   db,
		Tmpl: tmpl,
	}

	// в целях упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter() // gorilla/mux
	r.HandleFunc("/", srv.ListAll).Methods("GET")
	r.HandleFunc("/items", srv.ListAll).Methods("GET")
	r.HandleFunc("/items/new", srv.ShowCreateForm).Methods("GET")
	r.HandleFunc("/items/new", srv.CreateItem).Methods("POST")
	r.HandleFunc("/items/{id}", srv.ShowUpdateForm).Methods("GET")
	r.HandleFunc("/items/{id}", srv.UpdateItem).Methods("POST")
	r.HandleFunc("/items/{id}", srv.DeleteItem).Methods("DELETE")

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, r)
	show("end of program. ", err)
}

// sql table repr.
type MysqlSimplePostItem struct {
	Id          int
	Title       string
	Description string
	Updated     sql.NullString
}

type MysqlSimpleHttpHandlers struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *MysqlSimpleHttpHandlers) ListAll(w http.ResponseWriter, r *http.Request) {
	items := []*MysqlSimplePostItem{}

	rows, err := h.DB.Query("SELECT id, title, updated FROM items")
	__err_panic(err)
	for rows.Next() {
		post := &MysqlSimplePostItem{}
		err = rows.Scan(&post.Id, &post.Title, &post.Updated)
		__err_panic(err)
		items = append(items, post)
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*MysqlSimplePostItem
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MysqlSimpleHttpHandlers) ShowCreateForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MysqlSimpleHttpHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	// в целях упрощения примера пропущена валидация
	result, err := h.DB.Exec(
		"INSERT INTO items (`title`, `description`) VALUES (?, ?)",
		r.FormValue("title"),
		r.FormValue("description"),
	)
	__err_panic(err)

	affected, err := result.RowsAffected()
	__err_panic(err)
	lastID, err := result.LastInsertId()
	__err_panic(err)

	fmt.Println("Insert: RowsAffected ", affected, "; LastInsertId ", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *MysqlSimpleHttpHandlers) ShowUpdateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	post := &MysqlSimplePostItem{}
	// QueryRow сам закрывает коннект
	row := h.DB.QueryRow("SELECT id, title, updated, description FROM items WHERE id = ?", id)

	err = row.Scan(&post.Id, &post.Title, &post.Updated, &post.Description)
	__err_panic(err)

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MysqlSimpleHttpHandlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	// в целях упрощения примера пропущена валидация
	result, err := h.DB.Exec(
		"UPDATE items SET"+
			"`title` = ?"+
			",`description` = ?"+
			",`updated` = ?"+
			"WHERE id = ?",
		r.FormValue("title"),
		r.FormValue("description"),
		"foo",
		id,
	)
	__err_panic(err)

	affected, err := result.RowsAffected()
	__err_panic(err)

	fmt.Println("Update: RowsAffected", affected, "; item id: ", id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *MysqlSimpleHttpHandlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	result, err := h.DB.Exec(
		"DELETE FROM items WHERE id = ?",
		id,
	)
	__err_panic(err)

	affected, err := result.RowsAffected()
	__err_panic(err)

	fmt.Println("Delete: RowsAffected", affected)

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
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
