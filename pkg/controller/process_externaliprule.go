package controller

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	snat "github.com/kubeovn/ces-controller/pkg/apis/snat/v1alpha1"
)

func (c *Controller) processNextExternalIPRuleWorkItem() bool {
	obj, shutdown := c.externalIPRuleWorkQueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.externalIPRuleWorkQueue.Done(obj)

		key, err := cache.MetaNamespaceKeyFunc(obj)
		if err != nil {
			c.externalIPRuleWorkQueue.Forget(obj)
			utilruntime.HandleError(err)
			return err
		}

		var eipRule *snat.ExternalIPRule
		var ok bool
		if eipRule, ok = obj.(*snat.ExternalIPRule); !ok {
			c.externalIPRuleWorkQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected ExternalIPRule in workqueue but got %#v", obj))
			return nil
		}

		if err := c.externalIPRuleSyncHandler(key, eipRule); err != nil {
			c.externalIPRuleWorkQueue.AddRateLimited(eipRule)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.externalIPRuleWorkQueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) externalIPRuleSyncHandler(key string, eipRule *snat.ExternalIPRule) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	klog.Infof("===============================>start sync externalIPRule[%s]", name)
	defer klog.Infof("===============================>end sync externalIPRule[%s]", name)

	var isDelete bool
	var rule *snat.ExternalIPRule
	if rule, err = c.externalIPRuleLister.ExternalIPRules(namespace).Get(name); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		isDelete = true
		err = nil
	} else {
		eipRule = rule
	}

	defer func() {
		if err != nil {
			c.recorder.Event(eipRule, corev1.EventTypeWarning, err.Error(), MessageResourceFailedSynced)
		}
	}()

	_ = isDelete
	// todo
	//if len(serviceEgressRuleList.Items) == 0 && len(namespaceEgressRuleList.Items) == 0 && len(clusterEgressruleList.Items) == 0 {
	//	klog.Info("not found Associated rules，don,t neet sync!!")
	//	return nil
	//}
	//err = c.as3Client.As3Request(&serviceEgressRuleList, &namespaceEgressRuleList, &clusterEgressruleList, &externalServicesList, &endpointList, &namespaceList,
	//	tntcfg, ruleType, isDelete)
	//if err != nil {
	//	klog.Error(err)
	//	return err
	//}
	c.recorder.Event(eipRule, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}
