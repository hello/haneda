// Code generated by protoc-gen-go.
// source: morpheus_ble.proto
// DO NOT EDIT!

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ErrorType int32

const (
	ErrorType_TIME_OUT                    ErrorType = 0
	ErrorType_NETWORK_ERROR               ErrorType = 1
	ErrorType_DEVICE_ALREADY_PAIRED       ErrorType = 2
	ErrorType_INTERNAL_DATA_ERROR         ErrorType = 3
	ErrorType_DEVICE_DATABASE_FULL        ErrorType = 4
	ErrorType_DEVICE_NO_MEMORY            ErrorType = 5
	ErrorType_INTERNAL_OPERATION_FAILED   ErrorType = 6
	ErrorType_NO_ENDPOINT_IN_RANGE        ErrorType = 7
	ErrorType_WLAN_CONNECTION_ERROR       ErrorType = 8
	ErrorType_FAIL_TO_OBTAIN_IP           ErrorType = 9
	ErrorType_WLAN_ENDPOINT_DELETE_FAILED ErrorType = 10
	ErrorType_INVALID_ACCOUNT_ID          ErrorType = 11
	ErrorType_FORCE_DATA_PUSH_FAILED      ErrorType = 12
	ErrorType_PROTOBUF_ENCODE_FAILED      ErrorType = 13
	ErrorType_PROTOBUF_DECODE_FAILED      ErrorType = 14
	ErrorType_SERVER_CONNECTION_TIMEOUT   ErrorType = 15
)

var ErrorType_name = map[int32]string{
	0:  "TIME_OUT",
	1:  "NETWORK_ERROR",
	2:  "DEVICE_ALREADY_PAIRED",
	3:  "INTERNAL_DATA_ERROR",
	4:  "DEVICE_DATABASE_FULL",
	5:  "DEVICE_NO_MEMORY",
	6:  "INTERNAL_OPERATION_FAILED",
	7:  "NO_ENDPOINT_IN_RANGE",
	8:  "WLAN_CONNECTION_ERROR",
	9:  "FAIL_TO_OBTAIN_IP",
	10: "WLAN_ENDPOINT_DELETE_FAILED",
	11: "INVALID_ACCOUNT_ID",
	12: "FORCE_DATA_PUSH_FAILED",
	13: "PROTOBUF_ENCODE_FAILED",
	14: "PROTOBUF_DECODE_FAILED",
	15: "SERVER_CONNECTION_TIMEOUT",
}
var ErrorType_value = map[string]int32{
	"TIME_OUT":                    0,
	"NETWORK_ERROR":               1,
	"DEVICE_ALREADY_PAIRED":       2,
	"INTERNAL_DATA_ERROR":         3,
	"DEVICE_DATABASE_FULL":        4,
	"DEVICE_NO_MEMORY":            5,
	"INTERNAL_OPERATION_FAILED":   6,
	"NO_ENDPOINT_IN_RANGE":        7,
	"WLAN_CONNECTION_ERROR":       8,
	"FAIL_TO_OBTAIN_IP":           9,
	"WLAN_ENDPOINT_DELETE_FAILED": 10,
	"INVALID_ACCOUNT_ID":          11,
	"FORCE_DATA_PUSH_FAILED":      12,
	"PROTOBUF_ENCODE_FAILED":      13,
	"PROTOBUF_DECODE_FAILED":      14,
	"SERVER_CONNECTION_TIMEOUT":   15,
}

func (x ErrorType) Enum() *ErrorType {
	p := new(ErrorType)
	*p = x
	return p
}
func (x ErrorType) String() string {
	return proto.EnumName(ErrorType_name, int32(x))
}
func (x *ErrorType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ErrorType_value, data, "ErrorType")
	if err != nil {
		return err
	}
	*x = ErrorType(value)
	return nil
}
func (ErrorType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

type WifiConnectionState int32

const (
	WifiConnectionState_NO_WLAN_CONNECTED WifiConnectionState = 0
	WifiConnectionState_WLAN_CONNECTING   WifiConnectionState = 1
	WifiConnectionState_WLAN_CONNECTED    WifiConnectionState = 2
	WifiConnectionState_IP_RETRIEVED      WifiConnectionState = 3
	WifiConnectionState_DNS_RESOLVED      WifiConnectionState = 4
	WifiConnectionState_SOCKET_CONNECTED  WifiConnectionState = 5
	WifiConnectionState_REQUEST_SENT      WifiConnectionState = 6
	WifiConnectionState_CONNECTED         WifiConnectionState = 7
	WifiConnectionState_SSL_FAIL          WifiConnectionState = 8
	WifiConnectionState_HELLO_KEY_FAIL    WifiConnectionState = 9
	WifiConnectionState_DNS_FAILED        WifiConnectionState = 10
	WifiConnectionState_CONNECT_FAILED    WifiConnectionState = 11
)

var WifiConnectionState_name = map[int32]string{
	0:  "NO_WLAN_CONNECTED",
	1:  "WLAN_CONNECTING",
	2:  "WLAN_CONNECTED",
	3:  "IP_RETRIEVED",
	4:  "DNS_RESOLVED",
	5:  "SOCKET_CONNECTED",
	6:  "REQUEST_SENT",
	7:  "CONNECTED",
	8:  "SSL_FAIL",
	9:  "HELLO_KEY_FAIL",
	10: "DNS_FAILED",
	11: "CONNECT_FAILED",
}
var WifiConnectionState_value = map[string]int32{
	"NO_WLAN_CONNECTED": 0,
	"WLAN_CONNECTING":   1,
	"WLAN_CONNECTED":    2,
	"IP_RETRIEVED":      3,
	"DNS_RESOLVED":      4,
	"SOCKET_CONNECTED":  5,
	"REQUEST_SENT":      6,
	"CONNECTED":         7,
	"SSL_FAIL":          8,
	"HELLO_KEY_FAIL":    9,
	"DNS_FAILED":        10,
	"CONNECT_FAILED":    11,
}

func (x WifiConnectionState) Enum() *WifiConnectionState {
	p := new(WifiConnectionState)
	*p = x
	return p
}
func (x WifiConnectionState) String() string {
	return proto.EnumName(WifiConnectionState_name, int32(x))
}
func (x *WifiConnectionState) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(WifiConnectionState_value, data, "WifiConnectionState")
	if err != nil {
		return err
	}
	*x = WifiConnectionState(value)
	return nil
}
func (WifiConnectionState) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

//
// _u8 ssid[MAXIMAL_SSID_LENGTH];
// _u8 ssid_len;
// _u8 sec_type;
// _u8 bssid[SL_BSSID_LENGTH];
// _i8 rssi;
type WifiEndpointSecType int32

const (
	WifiEndpoint_SL_SCAN_SEC_TYPE_OPEN WifiEndpointSecType = 0
	WifiEndpoint_SL_SCAN_SEC_TYPE_WEP  WifiEndpointSecType = 1
	WifiEndpoint_SL_SCAN_SEC_TYPE_WPA  WifiEndpointSecType = 2
	WifiEndpoint_SL_SCAN_SEC_TYPE_WPA2 WifiEndpointSecType = 3
)

var WifiEndpointSecType_name = map[int32]string{
	0: "SL_SCAN_SEC_TYPE_OPEN",
	1: "SL_SCAN_SEC_TYPE_WEP",
	2: "SL_SCAN_SEC_TYPE_WPA",
	3: "SL_SCAN_SEC_TYPE_WPA2",
}
var WifiEndpointSecType_value = map[string]int32{
	"SL_SCAN_SEC_TYPE_OPEN": 0,
	"SL_SCAN_SEC_TYPE_WEP":  1,
	"SL_SCAN_SEC_TYPE_WPA":  2,
	"SL_SCAN_SEC_TYPE_WPA2": 3,
}

func (x WifiEndpointSecType) Enum() *WifiEndpointSecType {
	p := new(WifiEndpointSecType)
	*p = x
	return p
}
func (x WifiEndpointSecType) String() string {
	return proto.EnumName(WifiEndpointSecType_name, int32(x))
}
func (x *WifiEndpointSecType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(WifiEndpointSecType_value, data, "WifiEndpointSecType")
	if err != nil {
		return err
	}
	*x = WifiEndpointSecType(value)
	return nil
}
func (WifiEndpointSecType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0, 0} }

type MorpheusCommand_CommandType int32

const (
	MorpheusCommand_MORPHEUS_COMMAND_SET_TIME                     MorpheusCommand_CommandType = 0
	MorpheusCommand_MORPHEUS_COMMAND_GET_TIME                     MorpheusCommand_CommandType = 1
	MorpheusCommand_MORPHEUS_COMMAND_SET_WIFI_ENDPOINT            MorpheusCommand_CommandType = 2
	MorpheusCommand_MORPHEUS_COMMAND_GET_WIFI_ENDPOINT            MorpheusCommand_CommandType = 3
	MorpheusCommand_MORPHEUS_COMMAND_SET_ALARMS                   MorpheusCommand_CommandType = 4
	MorpheusCommand_MORPHEUS_COMMAND_GET_ALARMS                   MorpheusCommand_CommandType = 5
	MorpheusCommand_MORPHEUS_COMMAND_SWITCH_TO_PAIRING_MODE       MorpheusCommand_CommandType = 6
	MorpheusCommand_MORPHEUS_COMMAND_SWITCH_TO_NORMAL_MODE        MorpheusCommand_CommandType = 7
	MorpheusCommand_MORPHEUS_COMMAND_START_WIFISCAN               MorpheusCommand_CommandType = 8
	MorpheusCommand_MORPHEUS_COMMAND_STOP_WIFISCAN                MorpheusCommand_CommandType = 9
	MorpheusCommand_MORPHEUS_COMMAND_GET_DEVICE_ID                MorpheusCommand_CommandType = 10
	MorpheusCommand_MORPHEUS_COMMAND_ERASE_PAIRED_PHONE           MorpheusCommand_CommandType = 11
	MorpheusCommand_MORPHEUS_COMMAND_PAIR_PILL                    MorpheusCommand_CommandType = 12
	MorpheusCommand_MORPHEUS_COMMAND_ERROR                        MorpheusCommand_CommandType = 13
	MorpheusCommand_MORPHEUS_COMMAND_PAIR_SENSE                   MorpheusCommand_CommandType = 14
	MorpheusCommand_MORPHEUS_COMMAND_UNPAIR_PILL                  MorpheusCommand_CommandType = 15
	MorpheusCommand_MORPHEUS_COMMAND_MORPHEUS_DFU_BEGIN           MorpheusCommand_CommandType = 16
	MorpheusCommand_MORPHEUS_COMMAND_PILL_DATA                    MorpheusCommand_CommandType = 17
	MorpheusCommand_MORPHEUS_COMMAND_PILL_HEARTBEAT               MorpheusCommand_CommandType = 18
	MorpheusCommand_MORPHEUS_COMMAND_PILL_DFU_BEGIN               MorpheusCommand_CommandType = 19
	MorpheusCommand_MORPHEUS_COMMAND_FACTORY_RESET                MorpheusCommand_CommandType = 20
	MorpheusCommand_MORPHEUS_COMMAND_PILL_SHAKES                  MorpheusCommand_CommandType = 21
	MorpheusCommand_MORPHEUS_COMMAND_SYNC_DEVICE_ID               MorpheusCommand_CommandType = 22
	MorpheusCommand_MORPHEUS_COMMAND_PHONE_BLE_CONNECTED          MorpheusCommand_CommandType = 23
	MorpheusCommand_MORPHEUS_COMMAND_PHONE_BLE_BONDED             MorpheusCommand_CommandType = 24
	MorpheusCommand_MORPHEUS_COMMAND_LED_BUSY                     MorpheusCommand_CommandType = 25
	MorpheusCommand_MORPHEUS_COMMAND_LED_TRIPPY                   MorpheusCommand_CommandType = 26
	MorpheusCommand_MORPHEUS_COMMAND_LED_OPERATION_FAILED         MorpheusCommand_CommandType = 27
	MorpheusCommand_MORPHEUS_COMMAND_SCAN_WIFI                    MorpheusCommand_CommandType = 28
	MorpheusCommand_MORPHEUS_COMMAND_GET_NEXT_WIFI_AP             MorpheusCommand_CommandType = 29
	MorpheusCommand_MORPHEUS_COMMAND_LED_OPERATION_SUCCESS        MorpheusCommand_CommandType = 30
	MorpheusCommand_MORPHEUS_COMMAND_PUSH_DATA_AFTER_SET_TIMEZONE MorpheusCommand_CommandType = 31
	MorpheusCommand_MORPHEUS_COMMAND_SET_COUNTRY_CODE             MorpheusCommand_CommandType = 32
	MorpheusCommand_MORPHEUS_COMMAND_SET_SERVER_IP                MorpheusCommand_CommandType = 33
	MorpheusCommand_MORPHEUS_COMMAND_CONNECTION_STATE             MorpheusCommand_CommandType = 34
	MorpheusCommand_MORPHEUS_COMMAND_PILL_PROX_DATA               MorpheusCommand_CommandType = 35
)

var MorpheusCommand_CommandType_name = map[int32]string{
	0:  "MORPHEUS_COMMAND_SET_TIME",
	1:  "MORPHEUS_COMMAND_GET_TIME",
	2:  "MORPHEUS_COMMAND_SET_WIFI_ENDPOINT",
	3:  "MORPHEUS_COMMAND_GET_WIFI_ENDPOINT",
	4:  "MORPHEUS_COMMAND_SET_ALARMS",
	5:  "MORPHEUS_COMMAND_GET_ALARMS",
	6:  "MORPHEUS_COMMAND_SWITCH_TO_PAIRING_MODE",
	7:  "MORPHEUS_COMMAND_SWITCH_TO_NORMAL_MODE",
	8:  "MORPHEUS_COMMAND_START_WIFISCAN",
	9:  "MORPHEUS_COMMAND_STOP_WIFISCAN",
	10: "MORPHEUS_COMMAND_GET_DEVICE_ID",
	11: "MORPHEUS_COMMAND_ERASE_PAIRED_PHONE",
	12: "MORPHEUS_COMMAND_PAIR_PILL",
	13: "MORPHEUS_COMMAND_ERROR",
	14: "MORPHEUS_COMMAND_PAIR_SENSE",
	15: "MORPHEUS_COMMAND_UNPAIR_PILL",
	16: "MORPHEUS_COMMAND_MORPHEUS_DFU_BEGIN",
	17: "MORPHEUS_COMMAND_PILL_DATA",
	18: "MORPHEUS_COMMAND_PILL_HEARTBEAT",
	19: "MORPHEUS_COMMAND_PILL_DFU_BEGIN",
	20: "MORPHEUS_COMMAND_FACTORY_RESET",
	21: "MORPHEUS_COMMAND_PILL_SHAKES",
	22: "MORPHEUS_COMMAND_SYNC_DEVICE_ID",
	23: "MORPHEUS_COMMAND_PHONE_BLE_CONNECTED",
	24: "MORPHEUS_COMMAND_PHONE_BLE_BONDED",
	25: "MORPHEUS_COMMAND_LED_BUSY",
	26: "MORPHEUS_COMMAND_LED_TRIPPY",
	27: "MORPHEUS_COMMAND_LED_OPERATION_FAILED",
	28: "MORPHEUS_COMMAND_SCAN_WIFI",
	29: "MORPHEUS_COMMAND_GET_NEXT_WIFI_AP",
	30: "MORPHEUS_COMMAND_LED_OPERATION_SUCCESS",
	31: "MORPHEUS_COMMAND_PUSH_DATA_AFTER_SET_TIMEZONE",
	32: "MORPHEUS_COMMAND_SET_COUNTRY_CODE",
	33: "MORPHEUS_COMMAND_SET_SERVER_IP",
	34: "MORPHEUS_COMMAND_CONNECTION_STATE",
	35: "MORPHEUS_COMMAND_PILL_PROX_DATA",
}
var MorpheusCommand_CommandType_value = map[string]int32{
	"MORPHEUS_COMMAND_SET_TIME":                     0,
	"MORPHEUS_COMMAND_GET_TIME":                     1,
	"MORPHEUS_COMMAND_SET_WIFI_ENDPOINT":            2,
	"MORPHEUS_COMMAND_GET_WIFI_ENDPOINT":            3,
	"MORPHEUS_COMMAND_SET_ALARMS":                   4,
	"MORPHEUS_COMMAND_GET_ALARMS":                   5,
	"MORPHEUS_COMMAND_SWITCH_TO_PAIRING_MODE":       6,
	"MORPHEUS_COMMAND_SWITCH_TO_NORMAL_MODE":        7,
	"MORPHEUS_COMMAND_START_WIFISCAN":               8,
	"MORPHEUS_COMMAND_STOP_WIFISCAN":                9,
	"MORPHEUS_COMMAND_GET_DEVICE_ID":                10,
	"MORPHEUS_COMMAND_ERASE_PAIRED_PHONE":           11,
	"MORPHEUS_COMMAND_PAIR_PILL":                    12,
	"MORPHEUS_COMMAND_ERROR":                        13,
	"MORPHEUS_COMMAND_PAIR_SENSE":                   14,
	"MORPHEUS_COMMAND_UNPAIR_PILL":                  15,
	"MORPHEUS_COMMAND_MORPHEUS_DFU_BEGIN":           16,
	"MORPHEUS_COMMAND_PILL_DATA":                    17,
	"MORPHEUS_COMMAND_PILL_HEARTBEAT":               18,
	"MORPHEUS_COMMAND_PILL_DFU_BEGIN":               19,
	"MORPHEUS_COMMAND_FACTORY_RESET":                20,
	"MORPHEUS_COMMAND_PILL_SHAKES":                  21,
	"MORPHEUS_COMMAND_SYNC_DEVICE_ID":               22,
	"MORPHEUS_COMMAND_PHONE_BLE_CONNECTED":          23,
	"MORPHEUS_COMMAND_PHONE_BLE_BONDED":             24,
	"MORPHEUS_COMMAND_LED_BUSY":                     25,
	"MORPHEUS_COMMAND_LED_TRIPPY":                   26,
	"MORPHEUS_COMMAND_LED_OPERATION_FAILED":         27,
	"MORPHEUS_COMMAND_SCAN_WIFI":                    28,
	"MORPHEUS_COMMAND_GET_NEXT_WIFI_AP":             29,
	"MORPHEUS_COMMAND_LED_OPERATION_SUCCESS":        30,
	"MORPHEUS_COMMAND_PUSH_DATA_AFTER_SET_TIMEZONE": 31,
	"MORPHEUS_COMMAND_SET_COUNTRY_CODE":             32,
	"MORPHEUS_COMMAND_SET_SERVER_IP":                33,
	"MORPHEUS_COMMAND_CONNECTION_STATE":             34,
	"MORPHEUS_COMMAND_PILL_PROX_DATA":               35,
}

func (x MorpheusCommand_CommandType) Enum() *MorpheusCommand_CommandType {
	p := new(MorpheusCommand_CommandType)
	*p = x
	return p
}
func (x MorpheusCommand_CommandType) String() string {
	return proto.EnumName(MorpheusCommand_CommandType_name, int32(x))
}
func (x *MorpheusCommand_CommandType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MorpheusCommand_CommandType_value, data, "MorpheusCommand_CommandType")
	if err != nil {
		return err
	}
	*x = MorpheusCommand_CommandType(value)
	return nil
}
func (MorpheusCommand_CommandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor2, []int{2, 0}
}

type WifiEndpoint struct {
	Ssid             *string              `protobuf:"bytes,1,req,name=ssid" json:"ssid,omitempty"`
	Bssid            []byte               `protobuf:"bytes,2,opt,name=bssid" json:"bssid,omitempty"`
	Rssi             *int32               `protobuf:"varint,4,req,name=rssi" json:"rssi,omitempty"`
	SecurityType     *WifiEndpointSecType `protobuf:"varint,5,req,name=security_type,enum=WifiEndpointSecType" json:"security_type,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *WifiEndpoint) Reset()                    { *m = WifiEndpoint{} }
func (m *WifiEndpoint) String() string            { return proto.CompactTextString(m) }
func (*WifiEndpoint) ProtoMessage()               {}
func (*WifiEndpoint) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *WifiEndpoint) GetSsid() string {
	if m != nil && m.Ssid != nil {
		return *m.Ssid
	}
	return ""
}

func (m *WifiEndpoint) GetBssid() []byte {
	if m != nil {
		return m.Bssid
	}
	return nil
}

func (m *WifiEndpoint) GetRssi() int32 {
	if m != nil && m.Rssi != nil {
		return *m.Rssi
	}
	return 0
}

func (m *WifiEndpoint) GetSecurityType() WifiEndpointSecType {
	if m != nil && m.SecurityType != nil {
		return *m.SecurityType
	}
	return WifiEndpoint_SL_SCAN_SEC_TYPE_OPEN
}

type PillData struct {
	DeviceId            *string `protobuf:"bytes,1,req,name=device_id" json:"device_id,omitempty"`
	BatteryLevel        *int32  `protobuf:"varint,2,opt,name=battery_level" json:"battery_level,omitempty"`
	Uptime              *int32  `protobuf:"varint,3,opt,name=uptime" json:"uptime,omitempty"`
	MotionDataEntrypted []byte  `protobuf:"bytes,4,opt,name=motion_data_entrypted" json:"motion_data_entrypted,omitempty"`
	FirmwareVersion     *int32  `protobuf:"varint,5,opt,name=firmware_version" json:"firmware_version,omitempty"`
	Timestamp           *uint64 `protobuf:"varint,6,req,name=timestamp" json:"timestamp,omitempty"`
	Rssi                *int32  `protobuf:"varint,7,opt,name=rssi" json:"rssi,omitempty"`
	XXX_unrecognized    []byte  `json:"-"`
}

func (m *PillData) Reset()                    { *m = PillData{} }
func (m *PillData) String() string            { return proto.CompactTextString(m) }
func (*PillData) ProtoMessage()               {}
func (*PillData) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *PillData) GetDeviceId() string {
	if m != nil && m.DeviceId != nil {
		return *m.DeviceId
	}
	return ""
}

func (m *PillData) GetBatteryLevel() int32 {
	if m != nil && m.BatteryLevel != nil {
		return *m.BatteryLevel
	}
	return 0
}

func (m *PillData) GetUptime() int32 {
	if m != nil && m.Uptime != nil {
		return *m.Uptime
	}
	return 0
}

func (m *PillData) GetMotionDataEntrypted() []byte {
	if m != nil {
		return m.MotionDataEntrypted
	}
	return nil
}

func (m *PillData) GetFirmwareVersion() int32 {
	if m != nil && m.FirmwareVersion != nil {
		return *m.FirmwareVersion
	}
	return 0
}

func (m *PillData) GetTimestamp() uint64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *PillData) GetRssi() int32 {
	if m != nil && m.Rssi != nil {
		return *m.Rssi
	}
	return 0
}

type MorpheusCommand struct {
	Version      *int32                       `protobuf:"varint,1,req,name=version" json:"version,omitempty"`
	Type         *MorpheusCommand_CommandType `protobuf:"varint,2,req,name=type,enum=MorpheusCommand_CommandType" json:"type,omitempty"`
	DeviceId     *string                      `protobuf:"bytes,3,opt,name=deviceId" json:"deviceId,omitempty"`
	AccountId    *string                      `protobuf:"bytes,4,opt,name=accountId" json:"accountId,omitempty"`
	Error        *ErrorType                   `protobuf:"varint,5,opt,name=error,enum=ErrorType" json:"error,omitempty"`
	WifiName     *string                      `protobuf:"bytes,6,opt,name=wifiName" json:"wifiName,omitempty"`
	WifiSSID     *string                      `protobuf:"bytes,7,opt,name=wifiSSID" json:"wifiSSID,omitempty"`
	WifiPassword *string                      `protobuf:"bytes,8,opt,name=wifiPassword" json:"wifiPassword,omitempty"`
	// optional int32 batteryLevel = 9;  // Removed field, dont reuse
	// optional int32 uptime = 10;       // Removed field, dont reuse
	// optional int32 motionData = 11;   // This field is removed, to make thing explicit, please DONOT remove the line
	// optional bytes motionDataEntrypted = 12;   // This field is removed, to make thing explicit, please DONOT remove the line
	FirmwareVersion     *int32               `protobuf:"varint,13,opt,name=firmwareVersion" json:"firmwareVersion,omitempty"`
	WifiScanResult      []*WifiEndpoint      `protobuf:"bytes,14,rep,name=wifi_scan_result" json:"wifi_scan_result,omitempty"`
	SecurityType        *WifiEndpointSecType `protobuf:"varint,15,opt,name=security_type,enum=WifiEndpointSecType" json:"security_type,omitempty"`
	PillData            *PillData            `protobuf:"bytes,16,opt,name=pill_data" json:"pill_data,omitempty"`
	WifiConnectionState *WifiConnectionState `protobuf:"varint,17,opt,name=wifi_connection_state,enum=WifiConnectionState" json:"wifi_connection_state,omitempty"`
	BleBondCount        *int32               `protobuf:"varint,18,opt,name=ble_bond_count" json:"ble_bond_count,omitempty"`
	CountryCode         *string              `protobuf:"bytes,19,opt,name=country_code" json:"country_code,omitempty"`
	AesKey              []byte               `protobuf:"bytes,20,opt,name=aes_key" json:"aes_key,omitempty"`
	TopVersion          *string              `protobuf:"bytes,21,opt,name=top_version" json:"top_version,omitempty"`
	ServerIp            *uint32              `protobuf:"varint,22,opt,name=server_ip" json:"server_ip,omitempty"`
	SocketErrorCode     *uint32              `protobuf:"varint,23,opt,name=socket_error_code" json:"socket_error_code,omitempty"`
	HttpResponseCode    *string              `protobuf:"bytes,24,opt,name=http_response_code" json:"http_response_code,omitempty"`
	AppVersion          *int32               `protobuf:"varint,25,opt,name=app_version" json:"app_version,omitempty"`
	XXX_unrecognized    []byte               `json:"-"`
}

func (m *MorpheusCommand) Reset()                    { *m = MorpheusCommand{} }
func (m *MorpheusCommand) String() string            { return proto.CompactTextString(m) }
func (*MorpheusCommand) ProtoMessage()               {}
func (*MorpheusCommand) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *MorpheusCommand) GetVersion() int32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *MorpheusCommand) GetType() MorpheusCommand_CommandType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return MorpheusCommand_MORPHEUS_COMMAND_SET_TIME
}

func (m *MorpheusCommand) GetDeviceId() string {
	if m != nil && m.DeviceId != nil {
		return *m.DeviceId
	}
	return ""
}

func (m *MorpheusCommand) GetAccountId() string {
	if m != nil && m.AccountId != nil {
		return *m.AccountId
	}
	return ""
}

func (m *MorpheusCommand) GetError() ErrorType {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ErrorType_TIME_OUT
}

func (m *MorpheusCommand) GetWifiName() string {
	if m != nil && m.WifiName != nil {
		return *m.WifiName
	}
	return ""
}

func (m *MorpheusCommand) GetWifiSSID() string {
	if m != nil && m.WifiSSID != nil {
		return *m.WifiSSID
	}
	return ""
}

func (m *MorpheusCommand) GetWifiPassword() string {
	if m != nil && m.WifiPassword != nil {
		return *m.WifiPassword
	}
	return ""
}

func (m *MorpheusCommand) GetFirmwareVersion() int32 {
	if m != nil && m.FirmwareVersion != nil {
		return *m.FirmwareVersion
	}
	return 0
}

func (m *MorpheusCommand) GetWifiScanResult() []*WifiEndpoint {
	if m != nil {
		return m.WifiScanResult
	}
	return nil
}

func (m *MorpheusCommand) GetSecurityType() WifiEndpointSecType {
	if m != nil && m.SecurityType != nil {
		return *m.SecurityType
	}
	return WifiEndpoint_SL_SCAN_SEC_TYPE_OPEN
}

func (m *MorpheusCommand) GetPillData() *PillData {
	if m != nil {
		return m.PillData
	}
	return nil
}

func (m *MorpheusCommand) GetWifiConnectionState() WifiConnectionState {
	if m != nil && m.WifiConnectionState != nil {
		return *m.WifiConnectionState
	}
	return WifiConnectionState_NO_WLAN_CONNECTED
}

func (m *MorpheusCommand) GetBleBondCount() int32 {
	if m != nil && m.BleBondCount != nil {
		return *m.BleBondCount
	}
	return 0
}

func (m *MorpheusCommand) GetCountryCode() string {
	if m != nil && m.CountryCode != nil {
		return *m.CountryCode
	}
	return ""
}

func (m *MorpheusCommand) GetAesKey() []byte {
	if m != nil {
		return m.AesKey
	}
	return nil
}

func (m *MorpheusCommand) GetTopVersion() string {
	if m != nil && m.TopVersion != nil {
		return *m.TopVersion
	}
	return ""
}

func (m *MorpheusCommand) GetServerIp() uint32 {
	if m != nil && m.ServerIp != nil {
		return *m.ServerIp
	}
	return 0
}

func (m *MorpheusCommand) GetSocketErrorCode() uint32 {
	if m != nil && m.SocketErrorCode != nil {
		return *m.SocketErrorCode
	}
	return 0
}

func (m *MorpheusCommand) GetHttpResponseCode() string {
	if m != nil && m.HttpResponseCode != nil {
		return *m.HttpResponseCode
	}
	return ""
}

func (m *MorpheusCommand) GetAppVersion() int32 {
	if m != nil && m.AppVersion != nil {
		return *m.AppVersion
	}
	return 0
}

type BatchedPillData struct {
	Pills            []*PillData `protobuf:"bytes,1,rep,name=pills" json:"pills,omitempty"`
	DeviceId         *string     `protobuf:"bytes,2,req,name=device_id" json:"device_id,omitempty"`
	Prox             []*PillData `protobuf:"bytes,3,rep,name=prox" json:"prox,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *BatchedPillData) Reset()                    { *m = BatchedPillData{} }
func (m *BatchedPillData) String() string            { return proto.CompactTextString(m) }
func (*BatchedPillData) ProtoMessage()               {}
func (*BatchedPillData) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *BatchedPillData) GetPills() []*PillData {
	if m != nil {
		return m.Pills
	}
	return nil
}

func (m *BatchedPillData) GetDeviceId() string {
	if m != nil && m.DeviceId != nil {
		return *m.DeviceId
	}
	return ""
}

func (m *BatchedPillData) GetProx() []*PillData {
	if m != nil {
		return m.Prox
	}
	return nil
}

func init() {
	proto.RegisterType((*WifiEndpoint)(nil), "wifi_endpoint")
	proto.RegisterType((*PillData)(nil), "pill_data")
	proto.RegisterType((*MorpheusCommand)(nil), "MorpheusCommand")
	proto.RegisterType((*BatchedPillData)(nil), "batched_pill_data")
	proto.RegisterEnum("ErrorType", ErrorType_name, ErrorType_value)
	proto.RegisterEnum("WifiConnectionState", WifiConnectionState_name, WifiConnectionState_value)
	proto.RegisterEnum("WifiEndpointSecType", WifiEndpointSecType_name, WifiEndpointSecType_value)
	proto.RegisterEnum("MorpheusCommand_CommandType", MorpheusCommand_CommandType_name, MorpheusCommand_CommandType_value)
}

var fileDescriptor2 = []byte{
	// 1442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x57, 0xdd, 0x72, 0xda, 0xc8,
	0x12, 0x0e, 0x06, 0x6c, 0x33, 0x18, 0x18, 0xc6, 0x7f, 0x72, 0x7e, 0x1d, 0x72, 0x92, 0xf8, 0xe4,
	0xd4, 0xa1, 0xea, 0xb8, 0xea, 0x3c, 0x80, 0x90, 0x06, 0x50, 0x59, 0x48, 0x5a, 0x49, 0xd8, 0x61,
	0xf7, 0x62, 0x0a, 0x83, 0x52, 0xa6, 0x02, 0x48, 0x25, 0xc9, 0xc9, 0xfa, 0x11, 0xf6, 0x29, 0xf6,
	0x72, 0xdf, 0x63, 0x9f, 0x64, 0x1f, 0x63, 0x2f, 0xb7, 0x67, 0x10, 0x60, 0x9b, 0x71, 0xae, 0x6c,
	0x75, 0x7f, 0xd3, 0xd3, 0x3f, 0x5f, 0xf7, 0x34, 0x88, 0xcc, 0xc2, 0x38, 0xba, 0x09, 0x6e, 0x13,
	0x76, 0x3d, 0x0d, 0x9a, 0x51, 0x1c, 0xa6, 0x61, 0xe3, 0xaf, 0x1c, 0xaa, 0x7c, 0x9f, 0x7c, 0x99,
	0xb0, 0x60, 0x3e, 0x8e, 0xc2, 0xc9, 0x3c, 0x25, 0x7b, 0xa8, 0x90, 0x24, 0x93, 0xb1, 0x92, 0x3b,
	0xdd, 0x3a, 0x2b, 0x91, 0x0a, 0x2a, 0x5e, 0x8b, 0xcf, 0xad, 0xd3, 0xdc, 0xd9, 0x1e, 0x57, 0xc6,
	0xf0, 0xa9, 0x14, 0x40, 0x59, 0x24, 0x4d, 0x54, 0x49, 0x82, 0xd1, 0x6d, 0x3c, 0x49, 0xef, 0x58,
	0x7a, 0x17, 0x05, 0x4a, 0x11, 0xc4, 0xd5, 0xf3, 0xe3, 0xe6, 0x03, 0x8b, 0x4d, 0xc0, 0x08, 0x75,
	0x23, 0x45, 0xbb, 0xcb, 0xff, 0xc9, 0x09, 0x3a, 0xf4, 0x4c, 0xe6, 0x69, 0xaa, 0xc5, 0x3c, 0xaa,
	0x31, 0x7f, 0xe0, 0x50, 0x66, 0x3b, 0xd4, 0xc2, 0xcf, 0x88, 0x82, 0x0e, 0x36, 0x54, 0x57, 0xd4,
	0xc1, 0x39, 0xb9, 0xc6, 0x51, 0xf1, 0x96, 0xd4, 0x1c, 0x68, 0xce, 0x71, 0xbe, 0xf1, 0x7b, 0x0e,
	0x95, 0xa2, 0xc9, 0x74, 0xca, 0xc6, 0xc3, 0x74, 0x48, 0xea, 0xa8, 0x34, 0x0e, 0xbe, 0x4d, 0x46,
	0x01, 0x5b, 0xc5, 0x78, 0x88, 0x2a, 0xd7, 0xc3, 0x34, 0x0d, 0xe2, 0x3b, 0x36, 0x0d, 0xbe, 0x05,
	0x53, 0x11, 0x6b, 0x91, 0x54, 0xd1, 0xf6, 0x6d, 0x94, 0x4e, 0x66, 0x81, 0x92, 0x17, 0xdf, 0xaf,
	0xd0, 0xe1, 0x2c, 0x4c, 0x27, 0xe1, 0x5c, 0x18, 0x82, 0xf0, 0xd2, 0xf8, 0x2e, 0x4a, 0x83, 0x31,
	0x24, 0x83, 0xa7, 0x46, 0x41, 0xf8, 0xcb, 0x24, 0x9e, 0x7d, 0x1f, 0xc6, 0x01, 0xfb, 0x16, 0xc4,
	0x09, 0x00, 0x21, 0x1f, 0xfc, 0x20, 0x5c, 0xc9, 0xcd, 0x24, 0xe9, 0x70, 0x16, 0x29, 0xdb, 0x70,
	0x65, 0x61, 0x95, 0xc7, 0x1d, 0x0e, 0x68, 0xfc, 0x56, 0x43, 0xb5, 0x5e, 0x56, 0x1b, 0x2d, 0x9c,
	0xcd, 0x86, 0xf3, 0x31, 0xa9, 0xa1, 0x9d, 0xa5, 0x95, 0x9c, 0x48, 0xf6, 0x27, 0x54, 0x10, 0x39,
	0xde, 0x12, 0x39, 0x7e, 0xd9, 0x7c, 0x74, 0xa0, 0x99, 0xfd, 0xf5, 0x79, 0x72, 0x31, 0xda, 0x5d,
	0x04, 0x69, 0x8c, 0x85, 0xf3, 0x25, 0xee, 0xc3, 0x70, 0x34, 0x0a, 0x6f, 0xe7, 0xa9, 0xb1, 0x70,
	0xb8, 0x04, 0x29, 0x2b, 0x06, 0x71, 0x1c, 0xc6, 0xc2, 0xcb, 0xea, 0x39, 0x6a, 0x52, 0xfe, 0xb5,
	0x3c, 0xcf, 0x4b, 0x68, 0x0d, 0x21, 0xf8, 0x6d, 0x01, 0xce, 0x24, 0x9e, 0x67, 0xe8, 0xc2, 0xe9,
	0x12, 0x39, 0x40, 0x7b, 0x5c, 0xe2, 0x0c, 0x93, 0xe4, 0x7b, 0x18, 0x8f, 0x95, 0x5d, 0x21, 0x3d,
	0x46, 0xb5, 0x65, 0x16, 0x2e, 0x33, 0xf7, 0x2b, 0x22, 0x09, 0x67, 0x08, 0x0b, 0x56, 0x24, 0xa3,
	0xe1, 0x9c, 0xc5, 0x41, 0x72, 0x3b, 0x4d, 0x95, 0xea, 0x69, 0xfe, 0xac, 0x7c, 0x5e, 0x7d, 0x48,
	0x97, 0x4d, 0x56, 0xd5, 0x84, 0x7f, 0x4f, 0xb1, 0x0a, 0xea, 0xb2, 0x2e, 0xaf, 0x82, 0x01, 0x5b,
	0x86, 0x58, 0xd6, 0x05, 0xff, 0x3f, 0x3a, 0x14, 0x07, 0x47, 0xe1, 0x7c, 0x1e, 0x8c, 0x44, 0xfd,
	0xa0, 0x12, 0x69, 0xa0, 0xd4, 0x85, 0xd9, 0xa3, 0xa6, 0x54, 0x4b, 0x8e, 0x50, 0x15, 0xba, 0x84,
	0x5d, 0x87, 0xf3, 0x31, 0x13, 0x79, 0x53, 0x88, 0x88, 0x03, 0xc2, 0x16, 0x9f, 0x40, 0x96, 0x51,
	0x38, 0x0e, 0x94, 0x7d, 0x11, 0x36, 0x54, 0x6b, 0x18, 0x24, 0xec, 0x6b, 0x70, 0xa7, 0x1c, 0x08,
	0x36, 0xec, 0xa3, 0x72, 0x1a, 0x46, 0x2b, 0x22, 0x1c, 0x2e, 0x8b, 0x90, 0x04, 0x31, 0xc8, 0xd8,
	0x24, 0x52, 0x8e, 0x40, 0x54, 0x81, 0x22, 0xd4, 0x93, 0x70, 0xf4, 0x35, 0x48, 0x99, 0xa8, 0xc5,
	0xc2, 0xe6, 0xb1, 0x50, 0x3d, 0x47, 0xe4, 0x26, 0x4d, 0x23, 0x9e, 0xac, 0x28, 0x9c, 0x27, 0xc1,
	0x42, 0xa7, 0x08, 0x4b, 0x60, 0x7e, 0x18, 0xad, 0xcd, 0x9f, 0x08, 0x1a, 0xfd, 0x89, 0x50, 0xf9,
	0x3e, 0x0b, 0x5e, 0xa1, 0x93, 0x9e, 0xed, 0x3a, 0x5d, 0xda, 0xf7, 0x98, 0x66, 0xf7, 0x7a, 0xaa,
	0xa5, 0x43, 0x73, 0xf8, 0xcc, 0x37, 0x7a, 0x14, 0xda, 0x4c, 0xa6, 0xee, 0x2c, 0xd5, 0x39, 0xf2,
	0x01, 0x35, 0xa4, 0xa7, 0xaf, 0x8c, 0xb6, 0xc1, 0xa8, 0xa5, 0x3b, 0xb6, 0x61, 0xf9, 0xd0, 0x79,
	0x32, 0x5c, 0x67, 0x03, 0x97, 0x27, 0x6f, 0xd0, 0x0b, 0xa9, 0x3d, 0xd5, 0x54, 0xdd, 0x9e, 0x87,
	0x0b, 0x52, 0x40, 0x67, 0x0d, 0x28, 0x92, 0xff, 0xa0, 0x8f, 0x9b, 0x16, 0xae, 0x0c, 0x5f, 0xeb,
	0x32, 0xdf, 0x66, 0x8e, 0x6a, 0xb8, 0x86, 0xd5, 0x61, 0x3d, 0x5b, 0xa7, 0x78, 0x1b, 0xda, 0xe5,
	0xc3, 0x0f, 0xc0, 0x96, 0xed, 0xf6, 0x54, 0x73, 0x81, 0xdd, 0x21, 0xef, 0xd0, 0x9b, 0x4d, 0xac,
	0xaf, 0xba, 0x8b, 0x20, 0xf8, 0x50, 0xc1, 0xbb, 0xa4, 0x81, 0x5e, 0x4b, 0x40, 0xb6, 0xb3, 0xc6,
	0x94, 0xa4, 0x18, 0x1e, 0x82, 0x4e, 0x2f, 0x0d, 0x8d, 0x32, 0x43, 0xc7, 0x88, 0x7c, 0x44, 0xef,
	0x36, 0x30, 0xd4, 0x55, 0x3d, 0x2a, 0x22, 0xa0, 0x3a, 0x73, 0xba, 0xb6, 0x45, 0x71, 0x99, 0xbc,
	0x46, 0xcf, 0x37, 0x80, 0x1c, 0xc2, 0x1c, 0xc3, 0x34, 0xf1, 0x1e, 0xf0, 0xe3, 0x48, 0x62, 0xc8,
	0xb5, 0x5d, 0x5c, 0x91, 0xe6, 0x52, 0x9c, 0xf5, 0xa8, 0xe5, 0x51, 0x5c, 0x25, 0xa7, 0xe8, 0xe5,
	0x06, 0xa0, 0x6f, 0xad, 0xcd, 0xd7, 0xa4, 0x7e, 0xae, 0x04, 0x7a, 0xbb, 0xcf, 0x5a, 0xb4, 0x63,
	0x58, 0x18, 0xcb, 0xfd, 0x04, 0x1b, 0x4c, 0x57, 0x7d, 0x15, 0xd7, 0xa5, 0xd9, 0x15, 0xfa, 0x2e,
	0x85, 0x14, 0xb7, 0xa8, 0xea, 0x63, 0xf2, 0x34, 0x68, 0x7d, 0xd3, 0xbe, 0x34, 0xbd, 0x6d, 0x55,
	0xf3, 0x6d, 0x77, 0xc0, 0x5c, 0x0a, 0x64, 0xc2, 0x07, 0xd2, 0xc0, 0x84, 0x21, 0xaf, 0xab, 0x5e,
	0x50, 0x0f, 0x1f, 0xca, 0xab, 0x3d, 0xb0, 0xb4, 0x7b, 0x55, 0x3a, 0x82, 0x71, 0xf5, 0xaf, 0x4d,
	0x33, 0xbc, 0x30, 0xac, 0x65, 0x52, 0x90, 0x58, 0x16, 0xd5, 0x7c, 0xaa, 0xe3, 0x63, 0xf2, 0x1e,
	0xbd, 0xfd, 0x01, 0xb2, 0x65, 0x5b, 0x3a, 0xc0, 0x14, 0x69, 0xb7, 0x99, 0x50, 0xed, 0x56, 0xdf,
	0x1b, 0xe0, 0x13, 0x69, 0xc1, 0xb8, 0xda, 0x77, 0x0d, 0xc7, 0x19, 0xe0, 0xe7, 0xe4, 0xdf, 0xe8,
	0xbd, 0x14, 0x00, 0x6f, 0xa6, 0xab, 0xfa, 0x86, 0x6d, 0x41, 0x26, 0x0c, 0x10, 0xe0, 0x17, 0xd2,
	0x82, 0x88, 0x97, 0x91, 0x33, 0x15, 0xbf, 0x94, 0x7a, 0xcc, 0x59, 0x6a, 0xd1, 0xcf, 0x59, 0xdb,
	0xaa, 0x0e, 0x7e, 0x25, 0xed, 0xa0, 0x87, 0x37, 0x7a, 0x7d, 0x4d, 0xa3, 0x9e, 0x87, 0x5f, 0x93,
	0xff, 0xa1, 0xff, 0x6e, 0x26, 0xa1, 0xef, 0x75, 0x05, 0x07, 0x98, 0xda, 0xf6, 0xa9, 0xbb, 0x1a,
	0x3d, 0x3f, 0x73, 0x7a, 0xbf, 0x91, 0x7a, 0xc1, 0x21, 0x9a, 0xdd, 0xb7, 0x7c, 0x28, 0xa8, 0xc6,
	0x7b, 0xf3, 0x54, 0xde, 0x76, 0x00, 0xf3, 0xa8, 0x7b, 0x09, 0x46, 0x0d, 0x07, 0xbf, 0x95, 0x9a,
	0xca, 0x4a, 0x24, 0xdc, 0xf4, 0x55, 0x9f, 0xe2, 0xc6, 0xd3, 0x1c, 0x73, 0x5c, 0xfb, 0xf3, 0x82,
	0xad, 0xef, 0x1a, 0xbf, 0xa0, 0x3a, 0x2c, 0x03, 0xa3, 0x9b, 0x60, 0xcc, 0xd6, 0x6f, 0x08, 0x3c,
	0x95, 0xfc, 0x23, 0x81, 0xa7, 0x38, 0xff, 0xe8, 0x79, 0x79, 0xb0, 0x4f, 0x6c, 0x89, 0x7d, 0x42,
	0x41, 0x05, 0x58, 0xae, 0x7e, 0x85, 0x97, 0xf7, 0x11, 0xf8, 0xd3, 0x1f, 0x79, 0x54, 0x5a, 0xbf,
	0xb2, 0x7b, 0x68, 0x97, 0xe7, 0x83, 0xd9, 0x7d, 0x1f, 0xc6, 0x71, 0x1d, 0x55, 0x2c, 0xea, 0x5f,
	0xd9, 0xee, 0x45, 0xd6, 0xc5, 0x39, 0xbe, 0xd4, 0x64, 0x9c, 0x54, 0x4d, 0x97, 0xaa, 0xfa, 0x20,
	0x1b, 0x11, 0x30, 0x75, 0x8f, 0xd1, 0x3e, 0x8c, 0x55, 0xea, 0x5a, 0xea, 0xa2, 0xcf, 0xb2, 0x33,
	0x79, 0xbe, 0x22, 0x65, 0x67, 0xb8, 0xb8, 0xc5, 0xe7, 0x4a, 0xbb, 0x0f, 0x0d, 0x5d, 0x80, 0x97,
	0x0b, 0x67, 0x1a, 0xcb, 0x66, 0x3d, 0x0a, 0xa9, 0x18, 0x60, 0xbe, 0xd5, 0x9c, 0xac, 0x0c, 0x6d,
	0x70, 0x69, 0x9b, 0x9b, 0x03, 0xf4, 0x72, 0x8c, 0x33, 0xc3, 0x62, 0xae, 0x6a, 0x75, 0xf8, 0xd0,
	0x04, 0xe7, 0xae, 0x4c, 0x20, 0xd5, 0xbd, 0x44, 0x2f, 0x7c, 0xd8, 0x85, 0x85, 0xaa, 0xce, 0x0d,
	0xf0, 0x41, 0x6b, 0xb7, 0x7c, 0x15, 0x0e, 0x41, 0x99, 0x4a, 0x9c, 0xe3, 0xe2, 0xc4, 0xca, 0x9a,
	0x4e, 0x4d, 0xea, 0xd3, 0xe5, 0x65, 0x08, 0xde, 0x5c, 0x62, 0x58, 0x97, 0xaa, 0x69, 0xe8, 0x4c,
	0xd5, 0x04, 0x0f, 0x78, 0x33, 0x96, 0xf9, 0xa4, 0x6b, 0xdb, 0x6e, 0x16, 0xd2, 0x82, 0x57, 0xd9,
	0x19, 0x31, 0x05, 0xa1, 0x7c, 0xbe, 0xdd, 0xea, 0xb7, 0xc1, 0x30, 0x27, 0xcd, 0x52, 0x57, 0x79,
	0xa0, 0xd3, 0xe9, 0x7d, 0x5d, 0x95, 0xc7, 0x9d, 0x51, 0xe8, 0x5e, 0x00, 0xbc, 0x1c, 0xbc, 0x1a,
	0xb5, 0x4f, 0x7f, 0xe7, 0x9e, 0x58, 0x1b, 0x78, 0x70, 0x90, 0x91, 0xfb, 0xa1, 0x83, 0xbd, 0x67,
	0xf0, 0x22, 0xd7, 0x1e, 0xa4, 0xc3, 0xea, 0x40, 0x01, 0x09, 0xaa, 0x3e, 0x02, 0x6e, 0xc1, 0x26,
	0xb5, 0x67, 0x38, 0x30, 0xae, 0xa0, 0xb3, 0xe9, 0x25, 0x48, 0xf2, 0x5c, 0xa2, 0x5b, 0x1e, 0x9f,
	0x60, 0xb6, 0xc9, 0x25, 0xa2, 0x54, 0x9e, 0xad, 0x5d, 0x88, 0x6e, 0x58, 0x9e, 0x2c, 0x72, 0x9c,
	0x4b, 0x7f, 0xea, 0x53, 0x8f, 0xb3, 0x1f, 0xde, 0xd4, 0x6d, 0xd8, 0xce, 0x4b, 0x6b, 0xc0, 0x0e,
	0x27, 0x94, 0x07, 0x5b, 0x30, 0x8f, 0x11, 0xaa, 0x00, 0x97, 0x77, 0xa9, 0x69, 0xda, 0xec, 0x82,
	0x0e, 0x16, 0xb2, 0x12, 0xec, 0xb4, 0x88, 0x5f, 0xb5, 0xca, 0x38, 0x60, 0x32, 0x03, 0x4b, 0x59,
	0xb9, 0xd5, 0x44, 0xca, 0x28, 0x9c, 0x35, 0x6f, 0x82, 0xe9, 0x34, 0x6c, 0x26, 0xb0, 0x87, 0x45,
	0xb7, 0xcd, 0x61, 0x34, 0x69, 0xc2, 0x3a, 0xd4, 0x22, 0x5e, 0x00, 0x9b, 0x48, 0xb6, 0x64, 0x38,
	0xfc, 0x17, 0x44, 0xf2, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2a, 0xc7, 0x09, 0xaa, 0x57, 0x0c,
	0x00, 0x00,
}
