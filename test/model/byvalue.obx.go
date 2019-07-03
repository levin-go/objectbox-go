// Code generated by ObjectBox; DO NOT EDIT.
// Learn more about defining entities and generating this file - visit https://golang.objectbox.io/entity-annotations

package model

import (
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type entityByValue_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var EntityByValueBinding = entityByValue_EntityInfo{
	Entity: objectbox.Entity{
		Id: 3,
	},
	Uid: 2793387980842421409,
}

// EntityByValue_ contains type-based Property helpers to facilitate some common operations such as Queries.
var EntityByValue_ = struct {
	Id   *objectbox.PropertyUint64
	Text *objectbox.PropertyString
}{
	Id: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &EntityByValueBinding.Entity,
		},
	},
	Text: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &EntityByValueBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (entityByValue_EntityInfo) GeneratorVersion() int {
	return 2
}

// AddToModel is called by ObjectBox during model build
func (entityByValue_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("EntityByValue", 3, 2793387980842421409)
	model.Property("Id", 6, 1, 8853550994304785841)
	model.PropertyFlags(8193)
	model.Property("Text", 9, 2, 6704507893545428268)
	model.EntityLastPropertyId(2, 6704507893545428268)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (entityByValue_EntityInfo) GetId(object interface{}) (uint64, error) {
	if obj, ok := object.(*EntityByValue); ok {
		return obj.Id, nil
	} else {
		return object.(EntityByValue).Id, nil
	}
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (entityByValue_EntityInfo) SetId(object interface{}, id uint64) {
	if obj, ok := object.(*EntityByValue); ok {
		obj.Id = id
	} else {
		// NOTE while this can't update, it will at least behave consistently (panic in case of a wrong type)
		_ = object.(EntityByValue).Id
	}
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (entityByValue_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (entityByValue_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	var obj *EntityByValue
	if objPtr, ok := object.(*EntityByValue); ok {
		obj = objPtr
	} else {
		objVal := object.(EntityByValue)
		obj = &objVal
	}

	var offsetText = fbutils.CreateStringOffset(fbb, obj.Text)

	// build the FlatBuffers object
	fbb.StartObject(2)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetUOffsetTSlot(fbb, 1, offsetText)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (entityByValue_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}
	var id = table.GetUint64Slot(4, 0)

	return &EntityByValue{
		Id:   id,
		Text: fbutils.GetStringSlot(table, 6),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (entityByValue_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]EntityByValue, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (entityByValue_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	return append(slice.([]EntityByValue), *object.(*EntityByValue))
}

// Box provides CRUD access to EntityByValue objects
type EntityByValueBox struct {
	*objectbox.Box
}

// BoxForEntityByValue opens a box of EntityByValue objects
func BoxForEntityByValue(ob *objectbox.ObjectBox) *EntityByValueBox {
	return &EntityByValueBox{
		Box: ob.InternalBox(3),
	}
}

// Put synchronously inserts/updates a single object.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the EntityByValue.Id property on the passed object will be assigned the new ID as well.
func (box *EntityByValueBox) Put(object *EntityByValue) (uint64, error) {
	return box.Box.Put(object)
}

// PutAsync asynchronously inserts/updates a single object.
// When inserting, the EntityByValue.Id property on the passed object will be assigned the new ID as well.
//
// It's executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "Put & Forget:" you gain faster puts as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
//
// In situations with (extremely) high async load, this method may be throttled (~1ms) or delayed (<1s).
// In the unlikely event that the object could not be enqueued after delaying, an error will be returned.
//
// Note that this method does not give you hard durability guarantees like the synchronous Put provides.
// There is a small time window (typically 3 ms) in which the data may not have been committed durably yet.
func (box *EntityByValueBox) PutAsync(object *EntityByValue) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case Ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the EntityByValue.Id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the EntityByValue.Id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *EntityByValueBox) PutMany(objects []EntityByValue) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *EntityByValueBox) Get(id uint64) (*EntityByValue, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*EntityByValue), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is an empty object
func (box *EntityByValueBox) GetMany(ids ...uint64) ([]EntityByValue, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]EntityByValue), nil
}

// GetAll reads all stored objects
func (box *EntityByValueBox) GetAll() ([]EntityByValue, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]EntityByValue), nil
}

// Remove deletes a single object
func (box *EntityByValueBox) Remove(object *EntityByValue) (err error) {
	return box.Box.Remove(object.Id)
}

// Creates a query with the given conditions. Use the fields of the EntityByValue_ struct to create conditions.
// Keep the *EntityByValueQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *EntityByValueBox) Query(conditions ...objectbox.Condition) *EntityByValueQuery {
	return &EntityByValueQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the EntityByValue_ struct to create conditions.
// Keep the *EntityByValueQuery if you intend to execute the query multiple times.
func (box *EntityByValueBox) QueryOrError(conditions ...objectbox.Condition) (*EntityByValueQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &EntityByValueQuery{query}, nil
	}
}

// Query provides a way to search stored objects
//
// For example, you can find all EntityByValue which Id is either 42 or 47:
// 		box.Query(EntityByValue_.Id.In(42, 47)).Find()
type EntityByValueQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *EntityByValueQuery) Find() ([]EntityByValue, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]EntityByValue), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *EntityByValueQuery) Offset(offset uint64) *EntityByValueQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *EntityByValueQuery) Limit(limit uint64) *EntityByValueQuery {
	query.Query.Limit(limit)
	return query
}
