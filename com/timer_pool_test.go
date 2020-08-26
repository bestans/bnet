package com

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimer(t *testing.T)  {
	countTimer := int32(0)
	for j := 0; j < 100; j++ {
		for i := 0; i < 100; i++ {
			go func() {
				timer := NewTimer(time.Millisecond)
				defer timer.Recycle()

				select {
				case <-timer.Get():
				case <- time.After(time.Millisecond):
				}
				atomic.AddInt32(&countTimer, 1)
			}()
		}
		time.Sleep(time.Millisecond * 2)
	}
	time.Sleep(time.Second * 1)
	fmt.Println("count", countTimer)
	if countTimer != 10000 {
		t.Error("failed")
	}
}

func TestTimerFunc(t *testing.T)  {

	countTimer := int32(0)
	for j := 0; j < 100; j++ {
		for i := 0; i < 100; i++ {
			go func() {
					timer := GetTimerWithAfterFunc(time.Millisecond, func() {
					})
					defer PutTimerWithAfterFunc(timer)
					select {
					case <-time.After(time.Millisecond):
					}
			}()
		}
		for i := 0; i < 100; i++ {
			go func() {
				exit := make(chan bool)
				timer := GetTimerWithAfterFunc(time.Millisecond, func() {
					atomic.AddInt32(&countTimer, 1)
					exit <- true
				})
				defer PutTimerWithAfterFunc(timer)
				<-exit
			}()
		}
		time.Sleep(time.Millisecond * 2)
	}
	time.Sleep(time.Second * 1)
	fmt.Println("count", countTimer)
	if countTimer != 10000 {
		t.Error("failed")
	}
}