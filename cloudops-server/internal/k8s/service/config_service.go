package service

import (
	"context"
	"encoding/base64"
	"sort"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigService 提供配置资源能力。
type ConfigService struct {
	clusterSvc *ClusterService
}

// NewConfigService 创建配置服务。
func NewConfigService(clusterSvc *ClusterService) *ConfigService {
	return &ConfigService{clusterSvc: clusterSvc}
}

// ListConfigMaps 查询 ConfigMap 列表。
func (s *ConfigService) ListConfigMaps(clusterID int64, namespace string) ([]v1.ConfigMap, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

// CreateOrUpdateConfigMap 创建或更新 ConfigMap。
func (s *ConfigService) CreateOrUpdateConfigMap(clusterID int64, namespace string, item *v1.ConfigMap) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	old, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, item.Name, metav1.GetOptions{})
	if err != nil {
		_, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, item, metav1.CreateOptions{})
		return err
	}
	item.ResourceVersion = old.ResourceVersion
	_, err = client.CoreV1().ConfigMaps(namespace).Update(ctx, item, metav1.UpdateOptions{})
	return err
}

// DeleteConfigMap 删除 ConfigMap。
func (s *ConfigService) DeleteConfigMap(clusterID int64, namespace, name string) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// ListSecrets 查询 Secret 列表（值脱敏）。
func (s *ConfigService) ListSecrets(clusterID int64, namespace string, decode bool) ([]map[string]interface{}, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, 0, len(res.Items))
	for _, item := range res.Items {
		entry := map[string]interface{}{
			"name": item.Name,
			"type": string(item.Type),
			"data": map[string]string{},
		}
		data := map[string]string{}
		for k, v := range item.Data {
			if decode {
				data[k] = string(v)
			} else {
				data[k] = "******"
			}
		}
		entry["data"] = data
		list = append(list, entry)
	}
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i]["name"].(string), list[j]["name"].(string)) < 0
	})
	return list, nil
}

// CreateOrUpdateSecret 创建或更新 Secret。
func (s *ConfigService) CreateOrUpdateSecret(clusterID int64, namespace string, item *v1.Secret, base64Encoded bool) error {
	if !base64Encoded {
		if item.StringData == nil {
			item.StringData = map[string]string{}
		}
		for k, v := range item.Data {
			item.StringData[k] = base64.StdEncoding.EncodeToString(v)
		}
	}
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	old, err := client.CoreV1().Secrets(namespace).Get(ctx, item.Name, metav1.GetOptions{})
	if err != nil {
		_, err = client.CoreV1().Secrets(namespace).Create(ctx, item, metav1.CreateOptions{})
		return err
	}
	item.ResourceVersion = old.ResourceVersion
	_, err = client.CoreV1().Secrets(namespace).Update(ctx, item, metav1.UpdateOptions{})
	return err
}

// DeleteSecret 删除 Secret。
func (s *ConfigService) DeleteSecret(clusterID int64, namespace, name string) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
