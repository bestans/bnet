package com

import "bnet/inet"

type CoreStartUp struct {

}

func (self *CoreStartUp) StartUp(coms... inet.IStartup) error {
	for _, com := range coms {
		if err := com.Startup(); err != nil {
			return err
		}
	}
	return nil
}
