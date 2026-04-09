package service

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type hostConnectivityTarget struct {
	Protocol       string
	Address        string
	DisplayAddress string
}

func resolveHostConnectivityTarget(host model.CmdbHost) (hostConnectivityTarget, error) {
	if normalizeHostDeviceType(host.DeviceType) == "windows" {
		targetHost := firstNonEmptyConnectivity(host.RemoteDomain, host.SSHIP, host.PrivateIP, host.PublicIP)
		if targetHost == "" {
			return hostConnectivityTarget{}, errors.New("Windows 主机未配置 RDP 地址")
		}
		port := normalizeRDPPort(host.RDPPort)
		address := net.JoinHostPort(targetHost, strconv.Itoa(port))
		return hostConnectivityTarget{
			Protocol:       "rdp",
			Address:        address,
			DisplayAddress: address,
		}, nil
	}

	targetHost := firstNonEmptyConnectivity(host.SSHIP)
	if targetHost == "" || host.SSHPort <= 0 {
		return hostConnectivityTarget{}, errors.New("主机未配置 SSH 连接信息")
	}
	sshPort := host.SSHPort
	if sshPort <= 0 {
		sshPort = 22
	}
	address := net.JoinHostPort(targetHost, strconv.Itoa(sshPort))
	return hostConnectivityTarget{
		Protocol:       "ssh",
		Address:        address,
		DisplayAddress: address,
	}, nil
}

func (s *CmdbHostServiceImpl) BatchTestHostConnectivity(c *gin.Context, dto *model.BatchHostConnectivityDto) {
	if len(dto.HostIDs) == 0 {
		result.Failed(c, constant.INVALID_PARAMS, "请选择要测试的主机")
		return
	}

	items := make([]model.HostConnectivityCheckItem, 0, len(dto.HostIDs))
	summary := model.BatchHostConnectivityResult{}
	seen := make(map[uint]struct{}, len(dto.HostIDs))

	for _, hostID := range dto.HostIDs {
		if hostID == 0 {
			continue
		}
		if _, ok := seen[hostID]; ok {
			continue
		}
		seen[hostID] = struct{}{}

		item := model.HostConnectivityCheckItem{
			HostID: hostID,
			Status: "disconnected",
		}

		host, err := s.dao.GetCmdbHostById(hostID)
		if err != nil {
			item.Reason = "主机不存在或已被删除"
			items = append(items, item)
			continue
		}

		item.HostName = host.HostName
		item.DeviceType = normalizeHostDeviceType(host.DeviceType)

		target, err := resolveHostConnectivityTarget(host)
		if err != nil {
			item.Reason = err.Error()
			items = append(items, item)
			continue
		}

		item.Protocol = target.Protocol
		item.Address = target.Address
		item.DisplayAddress = target.DisplayAddress

		dialer := net.Dialer{Timeout: 3 * time.Second}
		conn, dialErr := dialer.DialContext(c.Request.Context(), "tcp", target.Address)
		if dialErr != nil {
			item.Reason = strings.TrimSpace(dialErr.Error())
			items = append(items, item)
			continue
		}
		_ = conn.Close()

		item.Reachable = true
		item.Status = "connected"
		item.Reason = ""
		summary.Reachable++
		items = append(items, item)
	}

	summary.Total = len(items)
	summary.Unreachable = summary.Total - summary.Reachable
	summary.Items = items
	result.Success(c, summary)
}

func firstNonEmptyConnectivity(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
