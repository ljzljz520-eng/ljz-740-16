package bindings

import (
	"fmt"
	"strings"
	"sync"
)

/*
采样方法（SampleMethod）与调度器（Scheduler）枚举定义。

与 stable-diffusion.h 中的 enum sample_method_t / enum scheduler_t 一一对应。

注意：每个枚举类型使用独立的 const 块，保证 iota 从 0 开始，
使 Go 常量的数值与 C 头文件中的枚举值严格一致。
（此前所有枚举混写在同一个 const 块中，iota 会持续递增，
  导致 EULER_SAMPLE_METHOD 等常量的实际值与 C 库不一致。）
*/

// SampleMethod 采样方法（采样器），对应 C 的 enum sample_method_t
type SampleMethod int

const (
	EULER_SAMPLE_METHOD SampleMethod = iota
	EULER_A_SAMPLE_METHOD
	HEUN_SAMPLE_METHOD
	DPM2_SAMPLE_METHOD
	DPMPP2S_A_SAMPLE_METHOD
	DPMPP2M_SAMPLE_METHOD
	DPMPP2Mv2_SAMPLE_METHOD
	IPNDM_SAMPLE_METHOD
	IPNDM_V_SAMPLE_METHOD
	LCM_SAMPLE_METHOD
	DDIM_TRAILING_SAMPLE_METHOD
	TCD_SAMPLE_METHOD
	RES_MULTISTEP_SAMPLE_METHOD
	RES_2S_SAMPLE_METHOD
	SAMPLE_METHOD_COUNT
)

// Scheduler 调度器，对应 C 的 enum scheduler_t
type Scheduler int

const (
	DISCRETE_SCHEDULER Scheduler = iota
	KARRAS_SCHEDULER
	EXPONENTIAL_SCHEDULER
	AYS_SCHEDULER
	GITS_SCHEDULER
	SGM_UNIFORM_SCHEDULER
	SIMPLE_SCHEDULER
	SMOOTHSTEP_SCHEDULER
	KL_OPTIMAL_SCHEDULER
	LCM_SCHEDULER
	BONG_TANGENT_SCHEDULER
	SCHEDULER_COUNT
)

// noneStr 是 C 库中 sd_*_name 系列函数对未知枚举值的返回值（NONE_STR）
const noneStr = "NONE"

// sampleMethodNames 与 stable-diffusion.cpp 中 sample_method_to_str 表保持一致
// （使用数组而非切片，使下面的编译期长度自检生效）
var sampleMethodNames = [...]string{
	"euler",
	"euler_a",
	"heun",
	"dpm2",
	"dpm++2s_a",
	"dpm++2m",
	"dpm++2mv2",
	"ipndm",
	"ipndm_v",
	"lcm",
	"ddim_trailing",
	"tcd",
	"res_multistep",
	"res_2s",
}

// schedulerNames 与 stable-diffusion.cpp 中 scheduler_to_str 表保持一致
var schedulerNames = [...]string{
	"discrete",
	"karras",
	"exponential",
	"ays",
	"gits",
	"sgm_uniform",
	"simple",
	"smoothstep",
	"kl_optimal",
	"lcm",
	"bong_tangent",
}

// 编译期自检：名称表长度必须与枚举数量一致
var _ = [1]struct{}{}[int(SAMPLE_METHOD_COUNT)-len(sampleMethodNames)]
var _ = [1]struct{}{}[int(SCHEDULER_COUNT)-len(schedulerNames)]

// ---------------------------------------------------------------------------
// 错误类型
// ---------------------------------------------------------------------------

// UnknownSampleMethodError 表示传入了 Go 绑定中不存在的采样方法（未知值）
type UnknownSampleMethodError struct {
	Value string
}

func (e *UnknownSampleMethodError) Error() string {
	return fmt.Sprintf("unknown sample method %q: supported values are: %s",
		e.Value, strings.Join(sampleMethodNames[:], ", "))
}

// UnknownSchedulerError 表示传入了 Go 绑定中不存在的调度器（未知值）
type UnknownSchedulerError struct {
	Value string
}

func (e *UnknownSchedulerError) Error() string {
	return fmt.Sprintf("unknown scheduler %q: supported values are: %s",
		e.Value, strings.Join(schedulerNames[:], ", "))
}

// UnsupportedSampleMethodError 表示该采样方法在 Go 绑定中已知，
// 但当前加载的 stable-diffusion 动态库版本过旧、尚不支持它
type UnsupportedSampleMethodError struct {
	Method          SampleMethod
	LibrarySupports int // 当前库实际支持的采样方法数量
}

func (e *UnsupportedSampleMethodError) Error() string {
	return fmt.Sprintf("sample method %q is not supported by the loaded stable-diffusion library "+
		"(it only supports %d sample methods): please upgrade the libstable-diffusion library file "+
		"(libstable-diffusion.so / libstable-diffusion.dylib / stable-diffusion.dll)",
		e.Method.String(), e.LibrarySupports)
}

// UnsupportedSchedulerError 表示该调度器在 Go 绑定中已知，
// 但当前加载的 stable-diffusion 动态库版本过旧、尚不支持它
type UnsupportedSchedulerError struct {
	Scheduler       Scheduler
	LibrarySupports int // 当前库实际支持的调度器数量
}

func (e *UnsupportedSchedulerError) Error() string {
	return fmt.Sprintf("scheduler %q is not supported by the loaded stable-diffusion library "+
		"(it only supports %d schedulers): please upgrade the libstable-diffusion library file "+
		"(libstable-diffusion.so / libstable-diffusion.dylib / stable-diffusion.dll)",
		e.Scheduler.String(), e.LibrarySupports)
}

// ---------------------------------------------------------------------------
// SampleMethod 方法
// ---------------------------------------------------------------------------

// String 返回采样方法在 C 库中的标准名称（如 "euler_a"）
func (m SampleMethod) String() string {
	if m >= 0 && int(m) < len(sampleMethodNames) {
		return sampleMethodNames[m]
	}
	return fmt.Sprintf("SampleMethod(%d)", int(m))
}

// IsValid 判断是否为 Go 绑定中已知的采样方法
func (m SampleMethod) IsValid() bool {
	return m >= 0 && m < SAMPLE_METHOD_COUNT
}

// Validate 校验采样方法是否可用：
//   - 未知枚举值返回 *UnknownSampleMethodError；
//   - 已知但当前加载的库版本过旧不支持时，返回 *UnsupportedSampleMethodError
//     （错误信息中提示需要升级库文件）。
func (m SampleMethod) Validate() error {
	if !m.IsValid() {
		return &UnknownSampleMethodError{Value: fmt.Sprintf("%d", int(m))}
	}
	if n := librarySampleMethodCount(); int(m) >= n {
		return &UnsupportedSampleMethodError{Method: m, LibrarySupports: n}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Scheduler 方法
// ---------------------------------------------------------------------------

// String 返回调度器在 C 库中的标准名称（如 "karras"）
func (s Scheduler) String() string {
	if s >= 0 && int(s) < len(schedulerNames) {
		return schedulerNames[s]
	}
	return fmt.Sprintf("Scheduler(%d)", int(s))
}

// IsValid 判断是否为 Go 绑定中已知的调度器
func (s Scheduler) IsValid() bool {
	return s >= 0 && s < SCHEDULER_COUNT
}

// Validate 校验调度器是否可用：
//   - 未知枚举值返回 *UnknownSchedulerError；
//   - 已知但当前加载的库版本过旧不支持时，返回 *UnsupportedSchedulerError
//     （错误信息中提示需要升级库文件）。
func (s Scheduler) Validate() error {
	if !s.IsValid() {
		return &UnknownSchedulerError{Value: fmt.Sprintf("%d", int(s))}
	}
	if n := librarySchedulerCount(); int(s) >= n {
		return &UnsupportedSchedulerError{Scheduler: s, LibrarySupports: n}
	}
	return nil
}

// ---------------------------------------------------------------------------
// 解析函数
// ---------------------------------------------------------------------------

// ParseSampleMethod 将字符串（如 "euler_a"）解析为 SampleMethod。
// 未知值返回 *UnknownSampleMethodError；
// 已知但当前库版本不支持时返回 *UnsupportedSampleMethodError（提示升级库文件）。
func ParseSampleMethod(s string) (SampleMethod, error) {
	for i, name := range sampleMethodNames {
		if s == name {
			m := SampleMethod(i)
			if err := m.Validate(); err != nil {
				return m, err
			}
			return m, nil
		}
	}
	return 0, &UnknownSampleMethodError{Value: s}
}

// ParseScheduler 将字符串（如 "karras"）解析为 Scheduler。
// 未知值返回 *UnknownSchedulerError；
// 已知但当前库版本不支持时返回 *UnsupportedSchedulerError（提示升级库文件）。
func ParseScheduler(s string) (Scheduler, error) {
	for i, name := range schedulerNames {
		if s == name {
			sch := Scheduler(i)
			if err := sch.Validate(); err != nil {
				return sch, err
			}
			return sch, nil
		}
	}
	return 0, &UnknownSchedulerError{Value: s}
}

// ---------------------------------------------------------------------------
// 支持项查询
// ---------------------------------------------------------------------------

// SupportedSampleMethods 返回当前加载的库实际支持的全部采样方法
func SupportedSampleMethods() []SampleMethod {
	n := librarySampleMethodCount()
	out := make([]SampleMethod, n)
	for i := range out {
		out[i] = SampleMethod(i)
	}
	return out
}

// SupportedSchedulers 返回当前加载的库实际支持的全部调度器
func SupportedSchedulers() []Scheduler {
	n := librarySchedulerCount()
	out := make([]Scheduler, n)
	for i := range out {
		out[i] = Scheduler(i)
	}
	return out
}

// ---------------------------------------------------------------------------
// 底层库能力探测
// ---------------------------------------------------------------------------

var (
	libSupportMu          sync.Mutex
	libSampleMethodCountN = -1 // 缓存：库实际支持的采样方法数量，-1 表示未探测
	libSchedulerCountN    = -1 // 缓存：库实际支持的调度器数量，-1 表示未探测
)

// librarySampleMethodCount 探测当前加载的库实际支持的采样方法数量。
// 原理：C 库的 sd_sample_method_name 对 >= SAMPLE_METHOD_COUNT 的值返回 "NONE"。
// 旧版本库的 SAMPLE_METHOD_COUNT 可能小于 Go 绑定编译时的值。
func librarySampleMethodCount() int {
	libSupportMu.Lock()
	defer libSupportMu.Unlock()
	if libSampleMethodCountN < 0 {
		n := 0
		for n < int(SAMPLE_METHOD_COUNT) && GoString(sdSampleMethodName(SampleMethod(n))) != noneStr {
			n++
		}
		libSampleMethodCountN = n
	}
	return libSampleMethodCountN
}

// librarySchedulerCount 探测当前加载的库实际支持的调度器数量，原理同上。
func librarySchedulerCount() int {
	libSupportMu.Lock()
	defer libSupportMu.Unlock()
	if libSchedulerCountN < 0 {
		n := 0
		for n < int(SCHEDULER_COUNT) && GoString(sdSchedulerName(Scheduler(n))) != noneStr {
			n++
		}
		libSchedulerCountN = n
	}
	return libSchedulerCountN
}

// resetLibrarySupportCache 重置库能力探测缓存（仅用于测试）
func resetLibrarySupportCache() {
	libSupportMu.Lock()
	defer libSupportMu.Unlock()
	libSampleMethodCountN = -1
	libSchedulerCountN = -1
}
