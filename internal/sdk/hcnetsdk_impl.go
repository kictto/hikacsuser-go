package sdk

/*
#cgo windows CFLAGS: -I./
#cgo windows LDFLAGS: -L${SRCDIR}/../../lib/win64 -lHCNetSDK

#cgo linux CFLAGS: -I./
#cgo linux LDFLAGS: -L${SRCDIR}/../../lib/linux64 -lhcnetsdk -ldl

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#ifdef _WIN32
#include <windows.h>
#else
#include <dlfcn.h>
#endif

typedef struct {
    char sDeviceAddress[129];
    unsigned char byUseTransport;
    unsigned short wPort;
    char sUserName[64];
    char sPassword[64];
    unsigned char bUseAsynLogin;
    char byRes[128];
    unsigned char byLoginMode;
    unsigned char byHttps;
    unsigned char byProxyType;    //0-不使用代理，1-使用标准代理，2-使用EHome代理
    unsigned char byUseUTCTime;   //0-不进行转换，默认,1-接口上输入输出全部使用UTC时间,2-接口上输入输出全部使用本地时间
    int iProxyID;                 //代理服务器序号
    unsigned char byVerifyMode;   //认证方式，0-不认证，1-双向认证，2-单向认证
    char byRes2[119];
} NET_DVR_USER_LOGIN_INFO;

typedef struct {
    char sSerialNumber[48];
    unsigned char byAlarmInPortNum;
    unsigned char byAlarmOutPortNum;
    unsigned char byDiskNum;
    unsigned char byDVRType;
    unsigned char byChanNum;
    unsigned char byStartChan;
    unsigned char byIPChanNum;
    unsigned char byRes1[24];
} NET_DVR_DEVICEINFO_V30;

typedef struct {
    unsigned int dwSize;
    NET_DVR_DEVICEINFO_V30 struDeviceV30;
    unsigned char bySupportLock;
    unsigned char byRetryLoginTime;
    unsigned char byPasswordLevel;
    unsigned char byRes1;
    unsigned int dwSurplusLockTime;
    unsigned char byCharEncodeType;
    unsigned char bySupportDev5;
    unsigned char bySupport;
    unsigned char byLoginMode;
    unsigned int dwOEMCode;
    int iResidualValidity;
    unsigned char byResidualValidity;
    unsigned char bySingleStartDTalkChan;
    unsigned char bySingleDTalkChanNums;
    unsigned char byPassWordResetLevel;
    unsigned char bySupportStreamEncrypt;
    unsigned char byMarketType;
    char byRes2[238];
} NET_DVR_DEVICEINFO_V40;

typedef struct {
    char sPath[256];
} NET_DVR_LOCAL_SDK_PATH;

typedef struct {
    unsigned int dwSize;
    unsigned char byLevel;           // 布防优先级，0-一等级（高），1-二等级（中），2-三等级（低）
    unsigned char byAlarmInfoType;   // 报警信息上传方式：0-老报警信息（NET_DVR_ALARMINFO），1-新报警信息(NET_DVR_ALARMINFO_V30)
    unsigned char byRetAlarmTypeV40; // V40报警信息类型
    unsigned char byRetDevInfoVersion; // V40报警信息对应设备信息版本号
    unsigned char byRetVQDAlarmType;   // VQD报警上传类型（用于报警类型区分）
    unsigned char byFaceAlarmDetection; // 人脸报警信息类型
    unsigned char bySupport;
    unsigned char byBrokenNetHttp;
    unsigned short wSeverityFilter; // 严重程度，用于SMART IPC
    unsigned char bySnapTimes;       // 设备联动抓图次数
    unsigned char bySnapSeq;         // 设备联动抓图序号
    unsigned char byRelRecordChan;
    unsigned char byRes1[12];
    unsigned char byChannel;         // 触发报警的通道号
    unsigned char byRes[35];
} NET_DVR_SETUPALARM_PARAM;

// SDK函数声明
int NET_DVR_Init();
int NET_DVR_Cleanup();
int NET_DVR_SetLogToFile(int iLogLevel, char* strLogDir, int bAutoDel);
int NET_DVR_SetSDKInitCfg(int enumType, void* lpInBuff);
int NET_DVR_Login_V40(NET_DVR_USER_LOGIN_INFO* pLoginInfo, NET_DVR_DEVICEINFO_V40* lpDeviceInfo);
int NET_DVR_Logout(int lUserID);
int NET_DVR_GetDVRConfig(int lUserID, unsigned int dwCommand, int lChannel, void* lpOutBuffer, unsigned int dwOutBufferSize, unsigned int* lpBytesReturned);
int NET_DVR_SetDVRConfig(int lUserID, unsigned int dwCommand, int lChannel, void* lpInBuffer, unsigned int dwInBufferSize);
long NET_DVR_StartRemoteConfig(int lUserID, unsigned int dwCommand, void* lpInBuffer, unsigned int dwInBufferSize, void* fRemoteConfigCallback, void* pUserData);
int NET_DVR_StopRemoteConfig(long lHandle);
int NET_DVR_SendWithRecvRemoteConfig(long lHandle, void* lpInBuffer, unsigned int dwInBufferSize, void* lpOutBuffer, unsigned int dwOutBufferSize, unsigned int* dwOutDataLen);
int NET_DVR_GetNextRemoteConfig(long lHandle, void* lpOutBuff, unsigned int dwOutBuffSize, unsigned int* lpOutDataLen);
unsigned int NET_DVR_GetLastError();
char* NET_DVR_GetErrorMsg(int* pErrorNo);
int NET_DVR_ControlGateway(int lUserID, int lGatewayIndex, unsigned int dwStaic);
int NET_DVR_SetupAlarmChan_V41(int lUserID, NET_DVR_SETUPALARM_PARAM* lpSetupParam);
int NET_DVR_CloseAlarmChan_V30(int lAlarmHandle);

// 图片查找和获取相关函数声明
int NET_DVR_FindPicture(int lUserID, void* pFindParam);
int NET_DVR_FindNextPicture(int lFindHandle, void* lpFindData);
int NET_DVR_CloseFindPicture(int lFindHandle);
int NET_DVR_GetPicture_V50(int lUserID, void* pPicParam, void* pBuffer);

// 定义报警回调函数类型
typedef void (*MSGCallBack)(int lCommand, void* pAlarmer, char* pAlarmInfo, unsigned int dwBufLen, void* pUser);

// 声明设置报警回调函数的C函数
int NET_DVR_SetDVRMessageCallBack_V50(int iIndex, MSGCallBack fMessageCallBack, void* pUser);

// 声明一个导出给C的回调函数桥接
extern int goMSGCallbackBridge(int lCommand, void* pAlarmer, void* pAlarmInfo, unsigned int dwBufLen, void* pUser);

// 定义报警设备信息结构体
typedef struct {
    unsigned char byUserIDValid;               // userid是否有效 0-无效，1-有效
    unsigned char bySerialValid;               // 序列号是否有效 0-无效，1-有效
    unsigned char byVersionValid;              // 版本号是否有效 0-无效，1-有效
    unsigned char byDeviceNameValid;           // 设备名字是否有效 0-无效，1-有效
    unsigned char byMacAddrValid;              // MAC地址是否有效 0-无效，1-有效
    unsigned char byLinkPortValid;             // login端口是否有效 0-无效，1-有效
    unsigned char byDeviceIPValid;             // 设备IP是否有效 0-无效，1-有效
    unsigned char bySocketIPValid;             // socket ip是否有效 0-无效，1-有效
    int lUserID;                     // NET_DVR_Login()返回值, 布防时有效
    char sSerialNumber[48];  // 序列号
    unsigned int dwDeviceVersion;              // 版本信息 高16位表示主版本，低16位表示次版本
    char sDeviceName[32];          // 设备名字
    unsigned char byMacAddr[6];      // MAC地址
    unsigned short wLinkPort;                   // link port
    char sDeviceIP[128];              // IP地址
    char sSocketIP[128];              // 报警主动上传时的socket IP地址
    unsigned char byIpProtocol;                // Ip协议 0-IPV4, 1-IPV6
    unsigned char byRes2[11];
} NET_DVR_ALARMER;

// 定义布防参数结构体

*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"
)

var (
	once        sync.Once
	sdkInstance *HCNetSDKImpl
)

// GetSDKInstance 获取SDK实例（单例模式）
func GetSDKInstance() (*HCNetSDKImpl, error) {
	var err error
	once.Do(func() {
		sdkInstance = &HCNetSDKImpl{}
		err = sdkInstance.initialize()
	})

	if err != nil {
		return nil, err
	}

	return sdkInstance, nil
}

// HCNetSDKImpl SDK实现
type HCNetSDKImpl struct {
	initialized bool
}

// initialize 初始化SDK
func (sdk *HCNetSDKImpl) initialize() error {
	// 设置库路径
	if err := sdk.setLibraryPath(); err != nil {
		return err
	}

	// 调用SDK初始化
	if !sdk.NET_DVR_Init() {
		return fmt.Errorf("NET_DVR_Init failed, error code: %d", sdk.NET_DVR_GetLastError())
	}

	sdk.initialized = true
	return nil
}

// setLibraryPath 设置SDK库路径
func (sdk *HCNetSDKImpl) setLibraryPath() error {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// 优先尝试从可执行文件目录查找lib
	execDir := filepath.Dir(execPath)
	libPath := filepath.Join(execDir, "lib")

	// 如果可执行文件目录下没有lib，则尝试从程序所在目录或源码目录查找
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		// 获取当前工作目录
		workDir, err := os.Getwd()
		if err == nil {
			libPath = filepath.Join(workDir, "lib")
			// 如果工作目录也没有lib，则使用相对于源码的固定路径
			if _, err := os.Stat(libPath); os.IsNotExist(err) {
				_, callerFile, _, _ := runtime.Caller(0)
				// 从源码文件位置确定lib目录
				sourceDir := filepath.Dir(callerFile)
				libPath = filepath.Join(sourceDir, "..", "..", "lib")
			}
		}
	}

	// 根据操作系统和架构选择合适的库路径
	var osDir string
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			osDir = "win64"
		} else {
			osDir = "win32"
		}
	} else if runtime.GOOS == "linux" {
		if runtime.GOARCH == "amd64" {
			osDir = "linux64"
		} else {
			osDir = "linux32"
		}
	} else {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	sdkPath := filepath.Join(libPath, osDir)

	// 打印找到的库路径便于调试
	fmt.Printf("使用SDK路径: %s\n", sdkPath)

	// Linux系统设置额外的库路径
	if runtime.GOOS == "linux" {
		cryptoPath := filepath.Join(sdkPath, "libcrypto.so.1.1")
		sslPath := filepath.Join(sdkPath, "libssl.so.1.1")

		cCryptoPath := C.CString(cryptoPath)
		cSSLPath := C.CString(sslPath)
		cSDKPath := C.CString(sdkPath)

		defer C.free(unsafe.Pointer(cCryptoPath))
		defer C.free(unsafe.Pointer(cSSLPath))
		defer C.free(unsafe.Pointer(cSDKPath))

		var struComPath C.NET_DVR_LOCAL_SDK_PATH
		C.strcpy((*C.char)(unsafe.Pointer(&struComPath.sPath[0])), cSDKPath)

		// 设置OpenSSL路径
		if C.NET_DVR_SetSDKInitCfg(3, unsafe.Pointer(cCryptoPath)) == 0 {
			return fmt.Errorf("set libcrypto path failed")
		}

		if C.NET_DVR_SetSDKInitCfg(4, unsafe.Pointer(cSSLPath)) == 0 {
			return fmt.Errorf("set libssl path failed")
		}

		// 设置SDK组件路径
		if C.NET_DVR_SetSDKInitCfg(2, unsafe.Pointer(&struComPath)) == 0 {
			return fmt.Errorf("set SDK path failed")
		}
	}

	return nil
}

// NET_DVR_Init 初始化SDK
func (sdk *HCNetSDKImpl) NET_DVR_Init() bool {
	return C.NET_DVR_Init() != 0
}

// NET_DVR_Cleanup 释放SDK资源
func (sdk *HCNetSDKImpl) NET_DVR_Cleanup() bool {
	return C.NET_DVR_Cleanup() != 0
}

// NET_DVR_SetLogToFile 设置日志
func (sdk *HCNetSDKImpl) NET_DVR_SetLogToFile(iLogLevel int, strLogDir string, bAutoDel bool) bool {
	cLogDir := C.CString(strLogDir)
	defer C.free(unsafe.Pointer(cLogDir))

	var autoDel int
	if bAutoDel {
		autoDel = 1
	}

	return C.NET_DVR_SetLogToFile(C.int(iLogLevel), cLogDir, C.int(autoDel)) != 0
}

// NET_DVR_SetSDKInitCfg 设置SDK初始化参数
func (sdk *HCNetSDKImpl) NET_DVR_SetSDKInitCfg(enumType int, lpInBuff unsafe.Pointer) bool {
	return C.NET_DVR_SetSDKInitCfg(C.int(enumType), lpInBuff) != 0
}

// NET_DVR_Login_V40 用户登录V40
func (sdk *HCNetSDKImpl) NET_DVR_Login_V40(pLoginInfo *NET_DVR_USER_LOGIN_INFO, lpDeviceInfo *NET_DVR_DEVICEINFO_V40) int {
	// 创建C结构体
	var cLoginInfo C.NET_DVR_USER_LOGIN_INFO
	var cDeviceInfo C.NET_DVR_DEVICEINFO_V40

	// 复制设备地址
	C.memcpy(unsafe.Pointer(&cLoginInfo.sDeviceAddress[0]),
		unsafe.Pointer(&pLoginInfo.SDeviceAddress[0]),
		C.sizeof_char*C.size_t(NET_DVR_DEV_ADDRESS_MAX_LEN))

	// 设置传输模式
	cLoginInfo.byUseTransport = C.uchar(pLoginInfo.ByUseTransport)

	// 复制用户名
	C.memcpy(unsafe.Pointer(&cLoginInfo.sUserName[0]),
		unsafe.Pointer(&pLoginInfo.SUserName[0]),
		C.sizeof_char*C.size_t(NET_DVR_LOGIN_USERNAME_MAX_LEN))

	// 复制密码
	C.memcpy(unsafe.Pointer(&cLoginInfo.sPassword[0]),
		unsafe.Pointer(&pLoginInfo.SPassword[0]),
		C.sizeof_char*C.size_t(NET_DVR_LOGIN_PASSWD_MAX_LEN))

	// 设置其他参数
	cLoginInfo.wPort = C.ushort(pLoginInfo.WPort)
	if pLoginInfo.BUseAsynLogin {
		cLoginInfo.bUseAsynLogin = 1
	} else {
		cLoginInfo.bUseAsynLogin = 0
	}
	cLoginInfo.byLoginMode = C.uchar(pLoginInfo.ByLoginMode)
	cLoginInfo.byHttps = C.uchar(pLoginInfo.ByHttps)

	// 设置新增的字段
	cLoginInfo.byProxyType = C.uchar(pLoginInfo.ByProxyType)
	cLoginInfo.byUseUTCTime = C.uchar(pLoginInfo.ByUseUTCTime)
	cLoginInfo.iProxyID = C.int(pLoginInfo.IProxyID)
	cLoginInfo.byVerifyMode = C.uchar(pLoginInfo.ByVerifyMode)

	// 调用C函数
	userID := C.NET_DVR_Login_V40(&cLoginInfo, &cDeviceInfo)

	// 复制设备信息
	lpDeviceInfo.StructSize = uint32(cDeviceInfo.dwSize)

	// 复制序列号
	C.memcpy(unsafe.Pointer(&lpDeviceInfo.DeviceInfo.SSerialNumber[0]),
		unsafe.Pointer(&cDeviceInfo.struDeviceV30.sSerialNumber[0]),
		C.sizeof_char*C.size_t(SERIALNO_LEN))

	// 复制其他字段
	lpDeviceInfo.DeviceInfo.ByAlarmInPortNum = byte(cDeviceInfo.struDeviceV30.byAlarmInPortNum)
	lpDeviceInfo.DeviceInfo.ByAlarmOutPortNum = byte(cDeviceInfo.struDeviceV30.byAlarmOutPortNum)
	lpDeviceInfo.DeviceInfo.ByDiskNum = byte(cDeviceInfo.struDeviceV30.byDiskNum)
	lpDeviceInfo.DeviceInfo.ByDVRType = byte(cDeviceInfo.struDeviceV30.byDVRType)
	lpDeviceInfo.DeviceInfo.ByChanNum = byte(cDeviceInfo.struDeviceV30.byChanNum)
	lpDeviceInfo.DeviceInfo.ByStartChan = byte(cDeviceInfo.struDeviceV30.byStartChan)
	lpDeviceInfo.DeviceInfo.ByIPChanNum = byte(cDeviceInfo.struDeviceV30.byIPChanNum)

	// 设置新增字段
	lpDeviceInfo.BySupportLock = byte(cDeviceInfo.bySupportLock)
	lpDeviceInfo.ByRetryLoginTime = byte(cDeviceInfo.byRetryLoginTime)
	lpDeviceInfo.ByPasswordLevel = byte(cDeviceInfo.byPasswordLevel)
	lpDeviceInfo.ByRes1 = byte(cDeviceInfo.byRes1)
	lpDeviceInfo.DwSurplusLockTime = uint32(cDeviceInfo.dwSurplusLockTime)
	lpDeviceInfo.ByCharEncodeType = byte(cDeviceInfo.byCharEncodeType)
	lpDeviceInfo.BySupportDev5 = byte(cDeviceInfo.bySupportDev5)
	lpDeviceInfo.BySupport = byte(cDeviceInfo.bySupport)
	lpDeviceInfo.ByLoginMode = byte(cDeviceInfo.byLoginMode)
	lpDeviceInfo.DwOEMCode = uint32(cDeviceInfo.dwOEMCode)
	lpDeviceInfo.IResidualValidity = int32(cDeviceInfo.iResidualValidity)
	lpDeviceInfo.ByResidualValidity = byte(cDeviceInfo.byResidualValidity)
	lpDeviceInfo.BySingleStartDTalkChan = byte(cDeviceInfo.bySingleStartDTalkChan)
	lpDeviceInfo.BySingleDTalkChanNums = byte(cDeviceInfo.bySingleDTalkChanNums)
	lpDeviceInfo.ByPassWordResetLevel = byte(cDeviceInfo.byPassWordResetLevel)
	lpDeviceInfo.BySupportStreamEncrypt = byte(cDeviceInfo.bySupportStreamEncrypt)
	lpDeviceInfo.ByMarketType = byte(cDeviceInfo.byMarketType)

	return int(userID)
}

// NET_DVR_Logout 用户注销
func (sdk *HCNetSDKImpl) NET_DVR_Logout(lUserID int) bool {
	return C.NET_DVR_Logout(C.int(lUserID)) != 0
}

// NET_DVR_GetDVRConfig 获取设备参数
func (sdk *HCNetSDKImpl) NET_DVR_GetDVRConfig(lUserID int, dwCommand uint32, lChannel int,
	lpOutBuffer unsafe.Pointer, dwOutBufferSize uint32, lpBytesReturned *uint32) bool {
	var cBytesReturned C.uint

	// 打印调试信息
	fmt.Printf("调用NET_DVR_GetDVRConfig: userID=%d, command=%d, channel=%d, bufferSize=%d, outBufferAddr=%v\n",
		lUserID, dwCommand, lChannel, dwOutBufferSize, lpOutBuffer)

	// 确保输出缓冲区指针有效
	if lpOutBuffer == nil {
		fmt.Println("错误: 输出缓冲区指针为空")
		return false
	}

	// 确保返回字节数指针有效
	if lpBytesReturned == nil {
		fmt.Println("错误: 返回字节数指针为空")
		return false
	}

	// 调用C函数
	result := C.NET_DVR_GetDVRConfig(C.int(lUserID), C.uint(dwCommand), C.int(lChannel),
		lpOutBuffer, C.uint(dwOutBufferSize), &cBytesReturned)

	// 更新返回字节数
	*lpBytesReturned = uint32(cBytesReturned)

	// 获取错误码
	errCode := C.NET_DVR_GetLastError()

	// 打印结果信息
	fmt.Printf("NET_DVR_GetDVRConfig结果: %v, 返回字节数: %d, 错误码: %d\n",
		result != 0, cBytesReturned, errCode)

	// 如果失败，获取错误信息
	if result == 0 {
		var errNo int32
		errMsg := sdk.NET_DVR_GetErrorMsg(&errNo)
		fmt.Printf("错误信息: %s (错误码: %d)\n", errMsg, errNo)
	}

	return result != 0
}

// NET_DVR_SetDVRConfig 设置设备参数
func (sdk *HCNetSDKImpl) NET_DVR_SetDVRConfig(lUserID int, dwCommand uint32, lChannel int,
	lpInBuffer unsafe.Pointer, dwInBufferSize uint32) bool {
	return C.NET_DVR_SetDVRConfig(C.int(lUserID), C.uint(dwCommand), C.int(lChannel),
		lpInBuffer, C.uint(dwInBufferSize)) != 0
}

// NET_DVR_StartRemoteConfig 启动远程配置
func (sdk *HCNetSDKImpl) NET_DVR_StartRemoteConfig(lUserID int, dwCommand uint32,
	lpInBuffer unsafe.Pointer, dwInBufferSize uint32,
	fRemoteConfigCallback uintptr, pUserData unsafe.Pointer) int64 {
	handle := C.NET_DVR_StartRemoteConfig(C.int(lUserID), C.uint(dwCommand),
		lpInBuffer, C.uint(dwInBufferSize),
		unsafe.Pointer(fRemoteConfigCallback), pUserData)

	return int64(handle)
}

// NET_DVR_StopRemoteConfig 停止远程配置
func (sdk *HCNetSDKImpl) NET_DVR_StopRemoteConfig(lHandle int64) bool {
	return C.NET_DVR_StopRemoteConfig(C.long(lHandle)) != 0
}

// NET_DVR_SendWithRecvRemoteConfig 发送接收数据
func (sdk *HCNetSDKImpl) NET_DVR_SendWithRecvRemoteConfig(lHandle int64, lpInBuff unsafe.Pointer, dwInBuffSize uint32,
	lpOutBuff unsafe.Pointer, dwOutBuffSize uint32, dwOutDataLen *uint32) int {
	var cOutDataLen C.uint
	result := C.NET_DVR_SendWithRecvRemoteConfig(C.long(lHandle), lpInBuff, C.uint(dwInBuffSize),
		lpOutBuff, C.uint(dwOutBuffSize), &cOutDataLen)

	if dwOutDataLen != nil {
		*dwOutDataLen = uint32(cOutDataLen)
	}

	return int(result)
}

// NET_DVR_GetNextRemoteConfig 获取下一个配置项
func (sdk *HCNetSDKImpl) NET_DVR_GetNextRemoteConfig(lHandle int64, lpOutBuff unsafe.Pointer, dwOutBuffSize uint32, dwOutDataLen *uint32) int {
	var cOutDataLen C.uint
	result := C.NET_DVR_GetNextRemoteConfig(C.long(lHandle), lpOutBuff, C.uint(dwOutBuffSize), &cOutDataLen)

	if dwOutDataLen != nil {
		*dwOutDataLen = uint32(cOutDataLen)
	}

	return int(result)
}

// NET_DVR_GetLastError 获取最后错误码
func (sdk *HCNetSDKImpl) NET_DVR_GetLastError() uint32 {
	return uint32(C.NET_DVR_GetLastError())
}

// NET_DVR_GetErrorMsg 获取错误码对应的错误信息
func (sdk *HCNetSDKImpl) NET_DVR_GetErrorMsg(pErrorNo *int32) string {
	var cErrorNo C.int
	if pErrorNo != nil {
		cErrorNo = C.int(*pErrorNo)
	}

	cMsg := C.NET_DVR_GetErrorMsg(&cErrorNo)
	if cMsg == nil {
		return ""
	}

	if pErrorNo != nil {
		*pErrorNo = int32(cErrorNo)
	}

	return C.GoString(cMsg)
}

// NET_DVR_ControlGateway 远程控制门操作
func (sdk *HCNetSDKImpl) NET_DVR_ControlGateway(lUserID int, lGatewayIndex int, dwStaic uint32) bool {
	return C.NET_DVR_ControlGateway(C.int(lUserID), C.int(lGatewayIndex), C.uint(dwStaic)) != 0
}

// NET_DVR_SetupAlarmChan_V41 建立报警上传通道，布防V41
func (sdk *HCNetSDKImpl) NET_DVR_SetupAlarmChan_V41(lUserID int, lpSetupParam *NET_DVR_SETUPALARM_PARAM) int {
	// 创建C结构体
	var cSetupParam C.NET_DVR_SETUPALARM_PARAM

	// 复制参数
	cSetupParam.dwSize = C.uint(lpSetupParam.DwSize)
	cSetupParam.byLevel = C.uchar(lpSetupParam.ByLevel)
	cSetupParam.byAlarmInfoType = C.uchar(lpSetupParam.ByAlarmInfoType)
	cSetupParam.byRetAlarmTypeV40 = C.uchar(lpSetupParam.ByRetAlarmTypeV40)
	cSetupParam.byRetDevInfoVersion = C.uchar(lpSetupParam.ByRetDevInfoVersion)
	cSetupParam.byRetVQDAlarmType = C.uchar(lpSetupParam.ByRetVQDAlarmType)
	cSetupParam.byFaceAlarmDetection = C.uchar(lpSetupParam.ByFaceAlarmDetection)
	cSetupParam.bySupport = C.uchar(lpSetupParam.BySupport)
	cSetupParam.byBrokenNetHttp = C.uchar(lpSetupParam.ByBrokenNetHttp)
	cSetupParam.wSeverityFilter = C.ushort(lpSetupParam.WSeverityFilter)
	cSetupParam.bySnapTimes = C.uchar(lpSetupParam.BySnapTimes)
	cSetupParam.bySnapSeq = C.uchar(lpSetupParam.BySnapSeq)
	cSetupParam.byRelRecordChan = C.uchar(lpSetupParam.ByRelRecordChan)

	// 调用C函数
	return int(C.NET_DVR_SetupAlarmChan_V41(C.int(lUserID), &cSetupParam))
}

// NET_DVR_CloseAlarmChan_V30 关闭报警上传通道，撤防V30
func (sdk *HCNetSDKImpl) NET_DVR_CloseAlarmChan_V30(lAlarmHandle int) bool {
	return C.NET_DVR_CloseAlarmChan_V30(C.int(lAlarmHandle)) != 0
}

const NET_DVR_CAPTURE_FACE_INFO = 2510 // 采集人脸信息

//export goMSGCallbackBridge
func goMSGCallbackBridge(lCommand C.int, pAlarmer unsafe.Pointer, pAlarmInfo unsafe.Pointer, dwBufLen C.uint, pUser unsafe.Pointer) C.int {
	// 将C结构体指针转换为Go结构体指针
	cAlarmer := (*C.NET_DVR_ALARMER)(pAlarmer)
	var goAlarmer NET_DVR_ALARMER

	// 复制字段 - 注意：这里直接访问C结构体字段
	goAlarmer.ByUserIDValid = byte(cAlarmer.byUserIDValid)
	goAlarmer.BySerialValid = byte(cAlarmer.bySerialValid)
	goAlarmer.ByVersionValid = byte(cAlarmer.byVersionValid)
	goAlarmer.ByDeviceNameValid = byte(cAlarmer.byDeviceNameValid)
	goAlarmer.ByMacAddrValid = byte(cAlarmer.byMacAddrValid)
	goAlarmer.ByLinkPortValid = byte(cAlarmer.byLinkPortValid)
	goAlarmer.ByDeviceIPValid = byte(cAlarmer.byDeviceIPValid)
	goAlarmer.BySocketIPValid = byte(cAlarmer.bySocketIPValid)
	goAlarmer.LUserID = int32(cAlarmer.lUserID)
	goAlarmer.DwDeviceVersion = uint32(cAlarmer.dwDeviceVersion)
	goAlarmer.WLinkPort = uint16(cAlarmer.wLinkPort)
	goAlarmer.ByIpProtocol = byte(cAlarmer.byIpProtocol)

	// 使用C.memcpy复制字节数组
	C.memcpy(unsafe.Pointer(&goAlarmer.SSerialNumber[0]), unsafe.Pointer(&cAlarmer.sSerialNumber[0]), C.sizeof_char*C.size_t(SERIALNO_LEN))
	C.memcpy(unsafe.Pointer(&goAlarmer.SDeviceName[0]), unsafe.Pointer(&cAlarmer.sDeviceName[0]), C.sizeof_char*C.size_t(NAME_LEN))
	C.memcpy(unsafe.Pointer(&goAlarmer.ByMacAddr[0]), unsafe.Pointer(&cAlarmer.byMacAddr[0]), C.sizeof_char*C.size_t(MACADDR_LEN))
	C.memcpy(unsafe.Pointer(&goAlarmer.SDeviceIP[0]), unsafe.Pointer(&cAlarmer.sDeviceIP[0]), C.sizeof_char*C.size_t(128))
	C.memcpy(unsafe.Pointer(&goAlarmer.SSocketIP[0]), unsafe.Pointer(&cAlarmer.sSocketIP[0]), C.sizeof_char*C.size_t(128))
	C.memcpy(unsafe.Pointer(&goAlarmer.ByRes2[0]), unsafe.Pointer(&cAlarmer.byRes2[0]), C.sizeof_char*C.size_t(11))

	// 获取并调用注册的回调函数 (假设索引为0，实际应用可能需要根据pUser查找)
	if callback, ok := getMSGCallback(0); ok {
		result := callback(int(lCommand), &goAlarmer, pAlarmInfo, uint32(dwBufLen), pUser)
		if result {
			return 1 // 返回1表示成功处理
		}
	}

	return 0 // 返回0表示未处理或处理失败
}

// NET_DVR_SetDVRMessageCallBack_V50 设置报警回调函数，V50版本
func (sdk *HCNetSDKImpl) NET_DVR_SetDVRMessageCallBack_V50(iIndex int, fMessageCallBack MSGCallBack_V31, pUser unsafe.Pointer) bool {
	// 注册回调函数
	registerMSGCallback(iIndex, fMessageCallBack)

	// 调用C函数设置回调
	// 直接传递导出的Go函数作为C函数指针
	return C.NET_DVR_SetDVRMessageCallBack_V50(C.int(iIndex), (C.MSGCallBack)(C.goMSGCallbackBridge), pUser) != 0
}

// 全局变量，用于保存回调函数
var (
	msgCallbacks    = make(map[int]MSGCallBack_V31)
	msgCallbackLock sync.Mutex
)

// 注册回调函数
func registerMSGCallback(index int, callback MSGCallBack_V31) {
	msgCallbackLock.Lock()
	defer msgCallbackLock.Unlock()
	msgCallbacks[index] = callback
}

// 获取回调函数
func getMSGCallback(index int) (MSGCallBack_V31, bool) {
	msgCallbackLock.Lock()
	defer msgCallbackLock.Unlock()
	callback, ok := msgCallbacks[index]
	return callback, ok
}

// NET_DVR_FindPicture 查找图片
func (sdk *HCNetSDKImpl) NET_DVR_FindPicture(lUserID int, pFindParam *NET_DVR_FIND_PICTURE_PARAM) int {
	return int(C.NET_DVR_FindPicture(C.int(lUserID), unsafe.Pointer(pFindParam)))
}

// NET_DVR_FindNextPicture 获取下一张图片
func (sdk *HCNetSDKImpl) NET_DVR_FindNextPicture(lFindHandle int, lpFindData *NET_DVR_FIND_PICTURE) int {
	return int(C.NET_DVR_FindNextPicture(C.int(lFindHandle), unsafe.Pointer(lpFindData)))
}

// NET_DVR_CloseFindPicture 关闭查找图片句柄
func (sdk *HCNetSDKImpl) NET_DVR_CloseFindPicture(lFindHandle int) bool {
	return C.NET_DVR_CloseFindPicture(C.int(lFindHandle)) != 0
}

// NET_DVR_GetPicture_V50 获取图片
func (sdk *HCNetSDKImpl) NET_DVR_GetPicture_V50(lUserID int, pPicParam *NET_DVR_FIND_PICTURE, pParam *NET_DVR_GETPIC_PARAM) bool {
	return C.NET_DVR_GetPicture_V50(C.int(lUserID), unsafe.Pointer(pPicParam), unsafe.Pointer(pParam)) != 0
}
