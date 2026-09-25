package bindings

import "strings"

// 本文件整理采样器（SampleMethod）与调度器（Scheduler）的枚举信息：
//   - 底层 C 库使用的规范字符串名称（与 stable-diffusion.cpp 中的
//     sample_method_names / scheduler_names 一致）
//   - String() 方法
//   - 对“当前加载的动态库是否支持该枚举”的运行时探测
//
// 较新版本的 stable-diffusion.cpp 会引入新的采样器/调度器，其枚举值追加在
// 末尾。若用户链接的是旧版动态库，传入新枚举值时 C 侧无法识别。此时
// IsXxxSupported 返回 false，高层 API 据此提示升级库文件。

// sampleMethodNames 按枚举值顺序排列，必须与 stable-diffusion.h 中的
// enum sample_method_t 严格对应。
var sampleMethodNames = [...]string{
	EULER_SAMPLE_METHOD:         "euler",
	EULER_A_SAMPLE_METHOD:       "euler_a",
	HEUN_SAMPLE_METHOD:          "heun",
	DPM2_SAMPLE_METHOD:          "dpm2",
	DPMPP2S_A_SAMPLE_METHOD:     "dpmpp_2s_a",
	DPMPP2M_SAMPLE_METHOD:       "dpmpp_2m",
	DPMPP2Mv2_SAMPLE_METHOD:     "dpmpp_2m_v2",
	IPNDM_SAMPLE_METHOD:         "ipndm",
	IPNDM_V_SAMPLE_METHOD:       "ipndm_v",
	LCM_SAMPLE_METHOD:           "lcm",
	DDIM_TRAILING_SAMPLE_METHOD: "ddim_trailing",
	TCD_SAMPLE_METHOD:           "tcd",
	RES_MULTISTEP_SAMPLE_METHOD: "res_multistep",
	RES_2S_SAMPLE_METHOD:        "res_2s",
}

// schedulerNames 按枚举值顺序排列，必须与 stable-diffusion.h 中的
// enum scheduler_t 严格对应。
var schedulerNames = [...]string{
	DISCRETE_SCHEDULER:     "discrete",
	KARRAS_SCHEDULER:       "karras",
	EXPONENTIAL_SCHEDULER:  "exponential",
	AYS_SCHEDULER:          "ays",
	GITS_SCHEDULER:         "gits",
	SGM_UNIFORM_SCHEDULER:  "sgm_uniform",
	SIMPLE_SCHEDULER:       "simple",
	SMOOTHSTEP_SCHEDULER:   "smoothstep",
	KL_OPTIMAL_SCHEDULER:   "kl_optimal",
	LCM_SCHEDULER:          "lcm",
	BONG_TANGENT_SCHEDULER: "bong_tangent",
}

// String 返回采样器在 C 库中的规范名称（如 "euler_a"）。
// 对未知枚举值返回空字符串。
func (m SampleMethod) String() string {
	if int(m) < 0 || int(m) >= len(sampleMethodNames) {
		return ""
	}
	return sampleMethodNames[m]
}

// String 返回调度器在 C 库中的规范名称（如 "karras"）。
// 对未知枚举值返回空字符串。
func (s Scheduler) String() string {
	if int(s) < 0 || int(s) >= len(schedulerNames) {
		return ""
	}
	return schedulerNames[s]
}

// Valid 判断采样器枚举值是否为本包已知的合法值。
// 合法只代表 Go 绑定层面认识该枚举，不代表当前动态库一定支持，
// 动态库支持情况请使用 IsSampleMethodSupported 判断。
func (m SampleMethod) Valid() bool {
	return int(m) >= 0 && int(m) < int(SAMPLE_METHOD_COUNT)
}

// Valid 判断调度器枚举值是否为本包已知的合法值。
// 合法只代表 Go 绑定层面认识该枚举，不代表当前动态库一定支持，
// 动态库支持情况请使用 IsSchedulerSupported 判断。
func (s Scheduler) Valid() bool {
	return int(s) >= 0 && int(s) < int(SCHEDULER_COUNT)
}

// IsSampleMethodSupported 判断当前加载的 stable-diffusion 动态库是否支持
// 指定采样器。
//
// 判定方式：调用库导出的 sd_sample_method_name，旧版库对不认识的枚举值
// 会返回 "unknown"（或空串）。
func IsSampleMethodSupported(m SampleMethod) bool {
	if !m.Valid() {
		return false
	}
	name := GetSampleMethodName(m)
	return name != "" && name != "unknown"
}

// IsSchedulerSupported 判断当前加载的 stable-diffusion 动态库是否支持
// 指定调度器。
//
// 判定方式：调用库导出的 sd_scheduler_name，旧版库对不认识的枚举值
// 会返回 "unknown"（或空串）。
func IsSchedulerSupported(s Scheduler) bool {
	if !s.Valid() {
		return false
	}
	name := GetSchedulerName(s)
	return name != "" && name != "unknown"
}

// ParseSampleMethod 将 C 库规范名称（如 "euler_a"）转换为采样器枚举。
// 大小写不敏感；返回的布尔值表示名称是否已知。
func ParseSampleMethod(name string) (SampleMethod, bool) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	for m, n := range sampleMethodNames {
		if n == normalized {
			return SampleMethod(m), true
		}
	}
	return SAMPLE_METHOD_COUNT, false
}

// ParseScheduler 将 C 库规范名称（如 "karras"）转换为调度器枚举。
// 大小写不敏感；返回的布尔值表示名称是否已知。
func ParseScheduler(name string) (Scheduler, bool) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	for s, n := range schedulerNames {
		if n == normalized {
			return Scheduler(s), true
		}
	}
	return SCHEDULER_COUNT, false
}
