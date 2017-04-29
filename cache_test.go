package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCaching(t *testing.T) {
	assert := assert.New(t)
	AddItemToCache("item1", []byte("value1"), time.Minute)
	AddItemToCache("item2", []byte("value2"), time.Minute)
	AddItemToCache("item3", []byte("value3"), time.Minute)
	assert.Equal(string(GetItemFromCache("item1").Value), "value1", "should be equal")
	assert.Equal(string(GetItemFromCache("item2").Value), "value2", "should be equal")
	assert.Equal(string(GetItemFromCache("item3").Value), "value3", "should be equal")
	removeItem("item1")
	assert.True(GetItemFromCache("item1") == nil, "should be equal")
}

func TestExpire(t *testing.T) {
	assert := assert.New(t)
	AddItemToCache("item", []byte("value"), time.Second)
	assert.Equal(string(GetItemFromCache("item").Value), "value", "should be equal")
	time.Sleep(time.Second * 2)
	assert.True(GetItemFromCache("item") == nil, "should be equal")
}

func TestCaching10000(t *testing.T) {
	t0 := time.Now()
	assert := assert.New(t)
	i := 0
	for i < 10000 {
		AddItemToCache("item"+strconv.Itoa(i), []byte("value"+strconv.Itoa(i)), time.Minute)
		assert.Equal(string(GetItemFromCache("item"+strconv.Itoa(i)).Value), "value"+strconv.Itoa(i), "should be equal")
		AddItemToCache("item"+strconv.Itoa(i), []byte("valuechanged"+strconv.Itoa(i)), time.Minute)
		assert.Equal(string(GetItemFromCache("item"+strconv.Itoa(i)).Value), "valuechanged"+strconv.Itoa(i), "should be equal")
		i++
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func TestCacheSizeLimit(t *testing.T) {
	t0 := time.Now()
	assert := assert.New(t)
	i := 0
	for i < 10002 {
		if i < 10001 {
			AddItemToCache("item"+strconv.Itoa(i), []byte("value"+strconv.Itoa(i)), time.Minute)
		} else {
			assert.True(
				assert.Panics(func() {
					AddItemToCache("item"+strconv.Itoa(i), []byte("value"+strconv.Itoa(i)), time.Minute)
				}, "should panic"),
			)
		}
		i++
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}
