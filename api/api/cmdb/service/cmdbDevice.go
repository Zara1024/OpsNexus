package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/util"
	"errors"
	"fmt"
	"net"
	neturl "net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	cmdbDeviceDefaultPageSize        = 10
	cmdbDeviceMaxPageSize            = 100
	cmdbDeviceConnectivityBatchLimit = 200
	cmdbDeviceConnectivityWorkers    = 10
)

var errInvalidCmdbDeviceParams = errors.New("invalid cmdb device params")

type CmdbDeviceService struct {
	dao *dao.CmdbDeviceDao
}

type cmdbDeviceConnectivityTarget struct {
	Protocol       string
	Address        string
	DisplayAddress string
}

func NewCmdbDeviceService(dao *dao.CmdbDeviceDao) *CmdbDeviceService {
	return &CmdbDeviceService{dao: dao}
}

func IsInvalidCmdbDeviceParams(err error) bool {
	return errors.Is(err, errInvalidCmdbDeviceParams)
}

func (s *CmdbDeviceService) CreateDevice(dto model.CreateCmdbDeviceDto) (*model.CmdbDevice, error) {
	device := normalizeCreateCmdbDeviceRequest(dto)
	if err := validateCmdbDevice(device); err != nil {
		return nil, err
	}

	now := util.HTime{Time: time.Now()}
	device.CreateTime = now
	device.UpdateTime = now

	if err := s.dao.Create(&device); err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *CmdbDeviceService) UpdateDevice(dto model.UpdateCmdbDeviceDto) (*model.CmdbDevice, error) {
	existing, err := s.dao.GetByID(dto.ID)
	if err != nil {
		return nil, err
	}

	device := normalizeUpdateCmdbDeviceRequest(*existing, dto)
	if err := validateCmdbDevice(device); err != nil {
		return nil, err
	}

	device.CreateTime = existing.CreateTime
	device.UpdateTime = util.HTime{Time: time.Now()}

	if err := s.dao.Update(&device); err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *CmdbDeviceService) DeleteDevice(id uint) error {
	if _, err := s.dao.GetByID(id); err != nil {
		return err
	}
	return s.dao.Delete(id)
}

func (s *CmdbDeviceService) ListDevices(page, pageSize int) ([]model.CmdbDevice, int64, error) {
	if page <= 0 {
		page = 1
	}
	pageSize = normalizeCmdbDevicePageSize(pageSize)
	return s.dao.List(page, pageSize)
}

func (s *CmdbDeviceService) GetDevice(id uint) (*model.CmdbDevice, error) {
	return s.dao.GetByID(id)
}

func (s *CmdbDeviceService) BatchTestDeviceConnectivity(dto model.BatchCmdbDeviceConnectivityDto) (*model.BatchCmdbDeviceConnectivityResult, error) {
	if len(dto.DeviceIDs) == 0 {
		return nil, invalidCmdbDeviceParams("deviceIds cannot be empty")
	}
	deviceIDs := normalizeCmdbDeviceConnectivityIDs(dto.DeviceIDs)
	if len(deviceIDs) == 0 {
		return nil, invalidCmdbDeviceParams("deviceIds cannot be empty")
	}
	if err := validateCmdbDeviceConnectivityBatchSize(deviceIDs); err != nil {
		return nil, err
	}

	devices, err := s.dao.GetByIDs(deviceIDs)
	if err != nil {
		return nil, err
	}

	deviceMap := make(map[uint]model.CmdbDevice, len(devices))
	for _, device := range devices {
		deviceMap[device.ID] = device
	}

	items := make([]model.CmdbDeviceConnectivityCheckItem, len(deviceIDs))
	sem := make(chan struct{}, cmdbDeviceConnectivityWorkers)
	var wg sync.WaitGroup

	for index, deviceID := range deviceIDs {
		wg.Add(1)
		go func(idx int, id uint) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			items[idx] = buildCmdbDeviceConnectivityItem(id, deviceMap)
		}(index, deviceID)
	}
	wg.Wait()

	summary := &model.BatchCmdbDeviceConnectivityResult{}
	summary.Total = len(items)
	summary.Items = items
	for _, item := range items {
		if item.Reachable {
			summary.Reachable++
		}
	}
	summary.Unreachable = summary.Total - summary.Reachable
	return summary, nil
}

func normalizeCreateCmdbDeviceRequest(dto model.CreateCmdbDeviceDto) model.CmdbDevice {
	return model.CmdbDevice{
		Name:          strings.TrimSpace(dto.Name),
		Address:       strings.TrimSpace(dto.Address),
		Platform:      strings.TrimSpace(dto.Platform),
		GroupID:       dto.GroupID,
		AccountID:     dto.AccountID,
		ProtocolGroup: normalizeCmdbDeviceProtocolGroup(dto.ProtocolGroup),
		Tags:          strings.TrimSpace(dto.Tags),
		IsActive:      normalizeCmdbDeviceActive(dto.IsActive),
		Remark:        strings.TrimSpace(dto.Remark),
		DeviceType:    normalizeCmdbDeviceType(dto.DeviceType),
		SSHPort:       normalizeCmdbDeviceSSHPort(dto.SSHPort),
		TelnetPort:    normalizeCmdbDeviceTelnetPort(dto.TelnetPort),
		WebURL:        strings.TrimSpace(dto.WebURL),
	}
}

func normalizeUpdateCmdbDeviceRequest(existing model.CmdbDevice, dto model.UpdateCmdbDeviceDto) model.CmdbDevice {
	device := normalizeCreateCmdbDeviceRequest(model.CreateCmdbDeviceDto{
		Name:          dto.Name,
		Address:       dto.Address,
		Platform:      dto.Platform,
		GroupID:       dto.GroupID,
		AccountID:     dto.AccountID,
		ProtocolGroup: dto.ProtocolGroup,
		Tags:          dto.Tags,
		IsActive:      dto.IsActive,
		Remark:        dto.Remark,
		DeviceType:    dto.DeviceType,
		SSHPort:       dto.SSHPort,
		TelnetPort:    dto.TelnetPort,
		WebURL:        dto.WebURL,
	})
	device.ID = existing.ID
	device.CreateTime = existing.CreateTime
	if dto.IsActive == nil {
		device.IsActive = existing.IsActive
	}
	return device
}

func normalizeCmdbDeviceProtocolGroup(value string) string {
	seen := map[string]bool{
		"ssh":    false,
		"telnet": false,
		"web":    false,
	}

	for _, item := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(strings.ToLower(item))
		if _, ok := seen[trimmed]; ok {
			seen[trimmed] = true
		}
	}

	ordered := make([]string, 0, len(seen))
	for _, protocol := range []string{"ssh", "telnet", "web"} {
		if seen[protocol] {
			ordered = append(ordered, protocol)
		}
	}
	if len(ordered) == 0 {
		return "ssh"
	}
	return strings.Join(ordered, ",")
}

func normalizeCmdbDeviceType(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return "network"
	}
	return trimmed
}

func normalizeCmdbDeviceActive(value *bool) bool {
	if value == nil {
		return true
	}
	return *value
}

func normalizeCmdbDeviceSSHPort(port int) int {
	if port > 0 {
		return port
	}
	return 22
}

func normalizeCmdbDeviceTelnetPort(port int) int {
	if port > 0 {
		return port
	}
	return 23
}

func validateCmdbDevice(device model.CmdbDevice) error {
	if device.Name == "" {
		return invalidCmdbDeviceParams("name cannot be empty")
	}
	if device.Address == "" {
		return invalidCmdbDeviceParams("address cannot be empty")
	}
	if device.GroupID == 0 {
		return invalidCmdbDeviceParams("groupId cannot be empty")
	}
	if device.AccountID == 0 {
		return invalidCmdbDeviceParams("accountId cannot be empty")
	}
	if device.SSHPort < 1 || device.SSHPort > 65535 {
		return invalidCmdbDeviceParams("sshPort must be between 1 and 65535")
	}
	if device.TelnetPort < 1 || device.TelnetPort > 65535 {
		return invalidCmdbDeviceParams("telnetPort must be between 1 and 65535")
	}
	return nil
}

func resolveCmdbDeviceConnectivityTarget(device model.CmdbDevice) (cmdbDeviceConnectivityTarget, error) {
	if canUseCmdbDeviceProtocol(device.ProtocolGroup, "ssh") && strings.TrimSpace(device.Address) != "" && device.SSHPort > 0 {
		address := net.JoinHostPort(strings.TrimSpace(device.Address), strconv.Itoa(device.SSHPort))
		return cmdbDeviceConnectivityTarget{
			Protocol:       "ssh",
			Address:        address,
			DisplayAddress: address,
		}, nil
	}

	if canUseCmdbDeviceProtocol(device.ProtocolGroup, "telnet") && strings.TrimSpace(device.Address) != "" && device.TelnetPort > 0 {
		address := net.JoinHostPort(strings.TrimSpace(device.Address), strconv.Itoa(device.TelnetPort))
		return cmdbDeviceConnectivityTarget{
			Protocol:       "telnet",
			Address:        address,
			DisplayAddress: address,
		}, nil
	}

	if canUseCmdbDeviceProtocol(device.ProtocolGroup, "web") && strings.TrimSpace(device.WebURL) != "" {
		parsed, err := parseCmdbDeviceWebURL(device.WebURL)
		if err != nil {
			return cmdbDeviceConnectivityTarget{}, err
		}

		hostPort := parsed.Host
		if parsed.Port() == "" {
			port := "80"
			if strings.EqualFold(parsed.Scheme, "https") {
				port = "443"
			}
			hostPort = net.JoinHostPort(parsed.Hostname(), port)
		}

		return cmdbDeviceConnectivityTarget{
			Protocol:       "web",
			Address:        hostPort,
			DisplayAddress: parsed.String(),
		}, nil
	}

	return cmdbDeviceConnectivityTarget{}, errors.New("device has no reachable connectivity target configured")
}

func canUseCmdbDeviceProtocol(protocolGroup, protocol string) bool {
	group := strings.TrimSpace(protocolGroup)
	if group == "" {
		return true
	}

	for _, item := range strings.Split(group, ",") {
		if strings.TrimSpace(strings.ToLower(item)) == protocol {
			return true
		}
	}
	return false
}

func parseCmdbDeviceWebURL(raw string) (*neturl.URL, error) {
	value := strings.TrimSpace(raw)
	if !strings.Contains(value, "://") {
		value = "http://" + value
	}

	parsed, err := neturl.Parse(value)
	if err != nil {
		return nil, err
	}
	if parsed.Host == "" {
		return nil, errors.New("device has no reachable connectivity target configured")
	}
	return parsed, nil
}

func normalizeCmdbDevicePageSize(pageSize int) int {
	if pageSize <= 0 {
		return cmdbDeviceDefaultPageSize
	}
	if pageSize > cmdbDeviceMaxPageSize {
		return cmdbDeviceMaxPageSize
	}
	return pageSize
}

func validateCmdbDeviceConnectivityBatchSize(deviceIDs []uint) error {
	if len(deviceIDs) > cmdbDeviceConnectivityBatchLimit {
		return invalidCmdbDeviceParams(fmt.Sprintf("deviceIds exceeds limit %d", cmdbDeviceConnectivityBatchLimit))
	}
	return nil
}

func invalidCmdbDeviceParams(message string) error {
	return fmt.Errorf("%w: %s", errInvalidCmdbDeviceParams, message)
}

func normalizeCmdbDeviceConnectivityIDs(deviceIDs []uint) []uint {
	normalized := make([]uint, 0, len(deviceIDs))
	seen := make(map[uint]struct{}, len(deviceIDs))

	for _, deviceID := range deviceIDs {
		if deviceID == 0 {
			continue
		}
		if _, ok := seen[deviceID]; ok {
			continue
		}
		seen[deviceID] = struct{}{}
		normalized = append(normalized, deviceID)
	}

	return normalized
}

func buildCmdbDeviceConnectivityItem(deviceID uint, deviceMap map[uint]model.CmdbDevice) model.CmdbDeviceConnectivityCheckItem {
	item := model.CmdbDeviceConnectivityCheckItem{
		DeviceID: deviceID,
		Status:   "disconnected",
	}

	device, ok := deviceMap[deviceID]
	if !ok {
		item.Reason = "device not found"
		return item
	}

	item.Name = device.Name

	target, err := resolveCmdbDeviceConnectivityTarget(device)
	if err != nil {
		item.Reason = err.Error()
		return item
	}

	item.Protocol = target.Protocol
	item.Address = target.Address
	item.DisplayAddress = target.DisplayAddress

	conn, dialErr := net.DialTimeout("tcp", target.Address, 3*time.Second)
	if dialErr != nil {
		item.Reason = strings.TrimSpace(dialErr.Error())
		return item
	}
	_ = conn.Close()

	item.Reachable = true
	item.Status = "connected"
	return item
}
