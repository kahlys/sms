// Package sms provides a generic interface around API-SMS Wrappers.
//
// The sms package must be used in conjuction with a api-sms driver.
// See https://godoc.org/github.com/kahlys/sms/driver for a list of drivers.
package sms

import (
	"fmt"
	"sort"
	"sync"
)

// Driver is the interface that must be implemented by a sms sender driver.
type Driver interface {
	Init(param map[string]string) error
	Send(msg, num string) error
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a sms driver available by the provided name. If Register is called
// twice with the same name or if driver is nil, it panics.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sms: Register driver is nil")
	}
	if _, ok := drivers[name]; ok {
		panic("sms: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// Sender is a sms sender handler.
type Sender struct {
	driver Driver
}

// Init inits a sms sender specified by its sms sender driver name and validates
// its parameters.
func Init(name string, params map[string]string) (*Sender, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()
	driveri, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("sms: unknown driver %q (forgotten import ?)", name)
	}
	if err := driveri.Init(params); err != nil {
		return nil, err
	}
	sender := &Sender{driver: driveri}
	return sender, nil
}

// Send sends a sms.
func (s *Sender) Send(msg, num string) error {
	return s.driver.Send(msg, num)
}
