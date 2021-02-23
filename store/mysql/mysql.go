package mysql

import (
	"fmt"
	"log"
	"time"
	"unicode"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"

	"github.com/CharlesBases/common/store"
)

var (
	// DefaultDatabase is the database that the sql store will use if no database is provided.
	DefaultDatabase = "store"
	// DefaultTable is the table that the sql store will use if no table is provided.
	DefaultTable = "store"
)

// New returns a new micro Store backed by sql
func NewStore(opts ...store.Option) store.Store {
	var options store.Options
	for _, o := range opts {
		o(&options)
	}

	// new store
	orm := new(ormStore)
	// set the options
	orm.options = options

	// configure the store
	if err := orm.configure(); err != nil {
		log.Fatal(err)
	}

	// return store
	return orm
}

// ormStore .
type ormStore struct {
	options store.Options

	db              *gorm.DB
	database, table string
}

// storeRecord .
type storeRecord struct {
	ID     int64  `xorm:"id"`
	Key    string `xorm:"key"`
	Value  []byte `xorm:"value"`
	Expiry int64  `xorm:"expiry"`
}

// set update or install
func (orm *ormStore) set(r *storeRecord) error {
	if err := orm.db.Exec("SELECT id FROM ? WHERE key = ?", orm.table, r.Key).First(&r.ID).Error; err != nil {
		return fmt.Errorf(`Couldn't insert record %s: %v`, r.Key, err)
	}

	switch r.ID {
	case 0:
		if err := orm.db.Create(r).Error; err != nil {
			return fmt.Errorf(`Couldn't insert record %s: %v`, r.Key, err)
		}
	default:
		err := orm.db.Transaction(func(tx *gorm.DB) error {
			return tx.Exec("SELECT id FROM ? WHERE id = ? LIMIT 1 FOR UPDATE; UPDATE ? SET value = ?, expiry = ? WHERE id = ?",
				orm.table, r.ID,
				orm.table, r.Value, r.Expiry, r.ID,
			).Error
		})
		if err != nil {
			return fmt.Errorf(`Couldn't insert record %s: %v`, r.Key, err)
		}
	}
	return nil
}

func (orm *ormStore) initDB() error {
	// Create the namespace's database
	if err := orm.db.Exec("CREATE DATABASE IF NOT EXISTS ? DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci", orm.database).Error; err != nil {
		return fmt.Errorf(`Couldn't create database %s: %v`, orm.database, err)
	}

	if err := orm.db.Exec("USE ?", orm.database).Error; err != nil {
		return fmt.Errorf(`Couldn't use database %s: %v`, orm.database, err)
	}

	// Create a table of the storeRecord
	if err := orm.db.Exec("CREATE TABLE IF NOT EXISTS ? (`id` bigint auto_increment, `key` varchar(255) not null, `value` blob not null, `expiry` bigint not null, constraint store_pk primary key (id)); CREATE UNIQUE INDEX uniq_key on ? (`key`)", orm.table, orm.table).Error; err != nil {
		return fmt.Errorf(`Couldn't create table %s.%s: %v`, orm.database, orm.table, err)
	}

	// setting
	orm.db.LogMode(false)
	orm.db.DB().SetMaxIdleConns(2000)
	orm.db.DB().SetMaxOpenConns(1000)
	orm.db.DB().SetConnMaxLifetime(time.Second * 10)

	return nil
}

func (orm *ormStore) configure() error {
	database := orm.options.Database
	if len(database) == 0 {
		database = DefaultDatabase
	}

	table := orm.options.Table
	if len(table) == 0 {
		table = DefaultTable
	}

	addrs := orm.options.Addresses
	if len(addrs) == 0 {
		addrs = []string{fmt.Sprintf(`root:@tcp("127.0.0.1:3306")/%s?charset=utf8mb4&parseTime=True&loc=Local`, database)}
	}

	for _, r := range database {
		if !unicode.IsLetter(r) {
			return errors.New("store.namespace must only contain letters")
		}
	}

	source := addrs[0]
	// create source from first node
	db, err := gorm.Open("mysql", source)
	if err != nil {
		return err
	}

	if err := db.DB().Ping(); err != nil {
		return err
	}

	if orm.db != nil {
		orm.db.Close()
	}

	// save the values
	orm.db = db
	orm.database = database
	orm.table = table

	// initialise the database
	return orm.initDB()
}

func (orm *ormStore) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&orm.options)
	}

	// reconfigure
	return orm.configure()
}

// Read all records with keys
func (orm *ormStore) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	var options store.ReadOptions
	for _, o := range opts {
		o(&options)
	}

	var records = make([]*store.Record, 0)
	var storeRecords = make([]*storeRecord, 0)

	session := orm.db.Table(orm.table)

	if options.Prefix {
		session = session.Where("key LIKE ?", fmt.Sprintf(`%s%%`, key))
		if options.Suffix {
			session = session.Or("key LIKE ?", fmt.Sprintf(`%%%s`, key))
		}
	} else if options.Suffix {
		session = session.Where("key LIKE ?", fmt.Sprintf(`%%%s`, key))
	} else {
		session = session.Where("key = ?", key)
	}

	if options.Limit != 0 {
		session = session.Limit(options.Limit)
		if options.Offset != 0 {
			session.Offset(options.Limit * options.Offset)
		}
	}
	if err := session.Order("id").Scan(&storeRecords).Error; err != nil {
		return records, store.ErrNotFound
	}

	for _, storeRecord := range storeRecords {
		if storeRecord.Expiry != 0 && time.Unix(storeRecord.Expiry, 0).Before(time.Now()) {
			go orm.Delete(storeRecord.Key)
			continue
		}

		records = append(
			records,
			&store.Record{
				Key:    storeRecord.Key,
				Value:  storeRecord.Value,
				Expiry: time.Unix(storeRecord.Expiry, 0),
			})
	}

	return records, nil
}

// Write records
func (orm *ormStore) Write(r *store.Record, opts ...store.WriteOption) error {
	var options store.WriteOptions
	for _, o := range opts {
		o(&options)
	}

	var storeRecord = new(storeRecord)
	storeRecord.Key = r.Key
	storeRecord.Value = r.Value

	if !r.Expiry.IsZero() {
		storeRecord.Expiry = r.Expiry.Unix()
	} else if r.TTL != 0 {
		storeRecord.Expiry = time.Now().Add(r.TTL).Unix()
	} else if !options.Expiry.IsZero() {
		storeRecord.Expiry = options.Expiry.Unix()
	} else if options.TTL != 0 {
		storeRecord.Expiry = time.Now().Add(options.TTL).Unix()
	}

	return orm.set(storeRecord)
}

// Delete records with keys
func (orm *ormStore) Delete(key string, opts ...store.DeleteOption) error {
	return orm.db.Exec("DELETE FROM ? WHERE key = ? LIMIT 1", orm.table, key).Error
}

// List all the known records
func (orm *ormStore) List(opts ...store.ListOption) ([]string, error) {
	var options store.ListOptions
	for _, o := range opts {
		o(&options)
	}

	var storeRecords = make([]*storeRecord, 0)

	session := orm.db.Table(orm.table)
	if options.Prefix != "" {
		session = session.Where("key LIKE ?", fmt.Sprintf(`%s%%`, options.Prefix))
		if options.Suffix != "" {
			session = session.Or("key LIKE ?", fmt.Sprintf(`%%%s`, options.Suffix))
		}
	} else if options.Suffix != "" {
		session = session.Where("key LIKE ?", fmt.Sprintf(`%%%s`, options.Suffix))
	}
	if options.Limit != 0 {
		session = session.Limit(options.Limit)
		if options.Offset != 0 {
			session = session.Offset(options.Limit * options.Offset)
		}
	}

	if err := session.Or("id").Find(&storeRecords).Error; err != nil {
		return nil, fmt.Errorf("Couldn't find keys %s.%s: %v", orm.database, orm.table, err)
	}

	var records = make([]string, 0)
	for _, storeRecord := range storeRecords {
		if storeRecord.Expiry != 0 && time.Unix(storeRecord.Expiry, 0).Before(time.Now()) {
			go orm.Delete(storeRecord.Key)
			continue
		}

		records = append(records, storeRecord.Key)
	}

	return records, nil
}

func (orm *ormStore) Close() error {
	return orm.db.Close()
}

func (orm *ormStore) String() string {
	return "mysql"
}

func (orm *ormStore) Options() store.Options {
	return orm.options
}
