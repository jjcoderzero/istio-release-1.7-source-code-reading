package crdclient

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	"istio.io/pkg/log"

	"istio.io/pkg/ledger"

	"istio.io/istio/pilot/pkg/model"
)

func (cl *Client) Version() string {
	return cl.configLedger.RootHash()
}

func (cl *Client) GetResourceAtVersion(version string, key string) (resourceVersion string, err error) {
	return cl.configLedger.GetPreviousValue(version, key)
}

func (cl *Client) GetLedger() ledger.Ledger {
	return cl.configLedger
}

func (cl *Client) SetLedger(l ledger.Ledger) error {
	cl.configLedger = l
	return nil
}

func castToObject(obj interface{}) (metav1.Object, error) {
	iobj, ok := obj.(metav1.Object)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return nil, fmt.Errorf("couldn't get object from tombstone %#v", obj)
		}
		iobj, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			return nil, fmt.Errorf("tombstone contained object that is not a object %#v", obj)
		}
	}
	return iobj, nil
}

func (cl *Client) tryLedgerPut(obj interface{}, kind string) {
	iobj, err := castToObject(obj)
	if err != nil {
		log.Errora(err)
		return
	}
	key := model.Key(kind, iobj.GetName(), iobj.GetNamespace())
	if _, err := cl.configLedger.Put(key, iobj.GetResourceVersion()); err != nil {
		scope.Errorf("Failed to update %s in ledger, status will be out of date.", key)
	}
}

func (cl *Client) tryLedgerDelete(obj interface{}, kind string) {
	iobj, err := castToObject(obj)
	if err != nil {
		log.Errora(err)
		return
	}
	key := model.Key(kind, iobj.GetName(), iobj.GetNamespace())
	if err := cl.configLedger.Delete(key); err != nil {
		scope.Errorf("Failed to delete %s in ledger, status will be out of date.", key)
	}
}
