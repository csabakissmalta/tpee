package store

import "github.com/csabakissmalta/tpee/exec"

type Store interface {
	// SaveData saves data to the target storage
	SaveData(interface{}, *exec.ExecRequestsElemDataPersistenceDataOutElem)

	// RetrieveData gets data from the storage
	RetrieveData(string) interface{}
}
