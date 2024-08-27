/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectbox

/*
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

const (
	// DebugflagsLogTransactionsRead enable read transaction logging
	DebugflagsLogTransactionsRead = 1

	// DebugflagsLogTransactionsWrite enable write transaction logging
	DebugflagsLogTransactionsWrite = 2

	// DebugflagsLogQueries enable query logging
	DebugflagsLogQueries = 4

	// DebugflagsLogQueryParameters enable query parameters logging
	DebugflagsLogQueryParameters = 8

	// DebugflagsLogAsyncQueue enable async operations logging
	DebugflagsLogAsyncQueue = 16
)

const (
	// Standard put ("insert or update")
	cPutModePut = 1

	// Put succeeds only if the entity does not exist yet.
	cPutModeInsert = 2

	// Put succeeds only if the entity already exist.
	cPutModeUpdate = 3

	// Not used yet (does not make sense for asnyc puts)
	// The given ID (non-zero) is guaranteed to be new; don't use unless you know exactly what you are doing!
	// This is primarily used internally. Wrong usage leads to inconsistent data (e.g. index data not updated)!
	cPutModePutIdGuaranteedToBeNew = 4
)

// atomic boolean true & false
const aTrue = 1
const aFalse = 0

// TypeId is a type of an internal ID on model/property/relation/index
type TypeId uint32

// ObjectBox provides super-fast object storage
type ObjectBox struct {
	store          *C.OBX_store
	entitiesById   map[TypeId]*entity
	entitiesByName map[string]*entity
	boxes          map[TypeId]*Box
	boxesMutex     sync.Mutex
	options        options
	syncClient     *SyncClient
}

type options struct {
	asyncTimeout uint
}

// constant during runtime so no need to call this each time it's necessary
// Armv7 linux device bool(C.obx_has_feature(C.OBXFeature_ResultArray)) is false, will case a panic
// var supportsResultArray = bool(C.obx_has_feature(C.OBXFeature_ResultArray))
var supportsResultArray = true

// Close fully closes the database and frees resources
func (ob *ObjectBox) Close() {
	storeToClose := ob.store
	ob.store = nil
	if ob.syncClient != nil {
		_ = ob.syncClient.Close()
	}
	if storeToClose != nil {
		C.obx_store_close(storeToClose)
	}
}

// RunInReadTx executes the given function inside a read transaction.
// The execution of the function `fn` must be sequential and executed in the same thread, which is enforced internally.
// If you launch goroutines inside `fn`, they will be executed on separate threads and not part of the same transaction.
// Multiple read transaction may be executed concurrently.
// The error returned by your callback is passed-through as the output error
func (ob *ObjectBox) RunInReadTx(fn func() error) error {
	return ob.runInTxn(true, fn)
}

// RunInWriteTx executes the given function inside a write transaction.
// The execution of the function `fn` must be sequential and executed in the same thread, which is enforced internally.
// If you launch goroutines inside `fn`, they will be executed on separate threads and not part of the same transaction.
// Only one write transaction may be active at a time (concurrently).
// The error returned by your callback is passed-through as the output error.
// If the resulting error is not nil, the transaction is aborted (rolled-back)
func (ob *ObjectBox) RunInWriteTx(fn func() error) error {
	return ob.runInTxn(false, fn)
}

func (ob *ObjectBox) runInTxn(readOnly bool, fn func() error) (err error) {
	// NOTE if runtime.LockOSThread() is about to be removed, evaluate use of createError() inside transactions
	runtime.LockOSThread()

	var cTxn *C.OBX_txn
	if readOnly {
		cTxn = C.obx_txn_read(ob.store)
	} else {
		cTxn = C.obx_txn_write(ob.store)
	}

	if cTxn == nil {
		err = createError()
		runtime.UnlockOSThread()
		return err
	}

	// Defer to ensure a TX is ALWAYS closed, even in a panic
	defer func() {
		if cTxn != nil {
			if rc := C.obx_txn_close(cTxn); rc != 0 {
				if err == nil {
					err = createError()
				} else {
					err = fmt.Errorf("%s; %s", err, createError())
				}
			}
		}

		runtime.UnlockOSThread()
	}()

	err = fn()

	if !readOnly && err == nil {
		var ptr = cTxn
		cTxn = nil
		if rc := C.obx_txn_success(ptr); rc != 0 {
			err = createError()
		}
	}

	return err
}

func (ob *ObjectBox) getEntityById(id TypeId) *entity {
	entity := ob.entitiesById[id]
	if entity == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no entity registered for entity ID " + strconv.Itoa(int(id)))
	}
	return entity
}

func (ob *ObjectBox) getEntityByName(name string) *entity {
	entity := ob.entitiesByName[name]
	if entity == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no entity registered for entity name " + name)
	}
	return entity
}

// SetDebugFlags configures debug logging of the ObjectBox core.
// See DebugFlags* constants
func (ob *ObjectBox) SetDebugFlags(flags uint) error {
	return cCall(func() C.obx_err {
		return C.obx_store_debug_flags(ob.store, C.uint32_t(flags))
	})
}

// InternalBox returns an Entity Box or panics on error (in case entity with the given ID doesn't exist)
func (ob *ObjectBox) InternalBox(entityId TypeId) *Box {
	box, err := ob.box(entityId)
	if err != nil {
		panic(fmt.Sprintf("Could not create box for entity ID %d: %s", entityId, err))
	}
	return box
}

// Gets an Entity Box which provides CRUD access to objects of the given type
func (ob *ObjectBox) box(entityId TypeId) (*Box, error) {
	ob.boxesMutex.Lock()
	defer ob.boxesMutex.Unlock()

	if box := ob.boxes[entityId]; box != nil {
		return box, nil
	}

	box, err := newBox(ob, entityId)
	if err != nil {
		return nil, err
	}

	ob.boxes[entityId] = box
	return box, nil
}

// AwaitAsyncCompletion blocks until all PutAsync insert have been processed
func (ob *ObjectBox) AwaitAsyncCompletion() error {
	return cCallBool(func() bool {
		return bool(C.obx_store_await_async_completion(ob.store))
	})
}

// SyncClient returns an existing client associated with the store or nil if not available.
// Use NewSyncClient() to create it the first time.
func (ob *ObjectBox) SyncClient() (*SyncClient, error) {
	if ob.syncClient == nil {
		return nil, errors.New("this store doesn't have a SyncClient associated, use NewSyncClient() to create one")
	}
	return ob.syncClient, nil
}
