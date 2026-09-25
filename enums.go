package stablediffusion

import (
	"fmt"
	"strings"

	"github.com/example/stablediffusion/bindings"
)

// 本文件将采样器与调度器整理为高层 Go 枚举常量，使用者无需直接引用
// bindings 包即可完成配置。所有常量的底层数值与 stable-diffusion.h 中
// 的 C 枚举一一对应，因此与直接使用 bindings 常量完全兼容。
//
// 传入未知值、或当前加载的动态库版本过旧不支持某个枚举时，GenerateImage /
// GenerateVideo 会在真正调用库之前返回错误；后者的错误信息会明确提示
// 需要升级库文件。

// Sampler 表示采样器（采样方法）。
// 它是 bindings.SampleMethod 的类型别名，可与底层常量互换使用。
type Sampler = bindings.SampleMethod

// Scheduler 表示调度器（噪声时间表）。
// 它是 bindings.Scheduler 的类型别名，可与底层常量互换使用。
type Scheduler = bindings.Scheduler

// 支持的采样器。名称与 C 库 sd_sample_method_name 返回的规范字符串一致。
const (
	SamplerEuler        Sampler = bindings.EULER_SAMPLE_METHOD
	SamplerEulerA       Sampler = bindings.EULER_A_SAMPLE_METHOD
	SamplerHeun         Sampler = bindings.HEUN_SAMPLE_METHOD
	SamplerDPM2         Sampler = bindings.DPM2_SAMPLE_METHOD
	SamplerDPMpp2SA     Sampler = bindings.DPMPP2S_A_SAMPLE_METHOD
	SamplerDPMpp2M      Sampler = bindings.DPMPP2M_SAMPLE_METHOD
	SamplerDPMpp2Mv2    Sampler = bindings.DPMPP2Mv2_SAMPLE_METHOD
	SamplerIPNDM        Sampler = bindings.IPNDM_SAMPLE_METHOD
	SamplerIPNDMv       Sampler = bindings.IPNDM_V_SAMPLE_METHOD
	SamplerLCM          Sampler = bindings.LCM_SAMPLE_METHOD
	SamplerDDIMTrailing Sampler = bindings.DDIM_TRAILING_SAMPLE_METHOD
	SamplerTCD          Sampler = bindings.TCD_SAMPLE_METHOD
	SamplerResMultistep Sampler = bindings.RES_MULTISTEP_SAMPLE_METHOD
	SamplerRes2s        Sampler = bindings.RES_2S_SAMPLE_METHOD
)

// 支持的调度器。名称与 C 库 sd_scheduler_name 返回的规范字符串一致。
const (
	SchedulerDiscrete    Scheduler = bindings.DISCRETE_SCHEDULER
	SchedulerKarras      Scheduler = bindings.KARRAS_SCHEDULER
	SchedulerExponential Scheduler = bindings.EXPONENTIAL_SCHEDULER
	SchedulerAys         Scheduler = bindings.AYS_SCHEDULER
	SchedulerGits        Scheduler = bindings.GITS_SCHEDULER
	SchedulerSgmUniform  Scheduler = bindings.SGM_UNIFORM_SCHEDULER
	SchedulerSimple      Scheduler = bindings.SIMPLE_SCHEDULER
	SchedulerSmoothstep  Scheduler = bindings.SMOOTHSTEP_SCHEDULER
	SchedulerKlOptimal   Scheduler = bindings.KL_OPTIMAL_SCHEDULER
	SchedulerLcm         Scheduler = bindings.LCM_SCHEDULER
	SchedulerBongTangent Scheduler = bindings.BONG_TANGENT_SCHEDULER
)

// allSamplers 按枚举顺序列出本包已知的全部采样器。
var allSamplers = []Sampler{
	SamplerEuler,
	SamplerEulerA,
	SamplerHeun,
	SamplerDPM2,
	SamplerDPMpp2SA,
	SamplerDPMpp2M,
	SamplerDPMpp2Mv2,
	SamplerIPNDM,
	SamplerIPNDMv,
	SamplerLCM,
	SamplerDDIMTrailing,
	SamplerTCD,
	SamplerResMultistep,
	SamplerRes2s,
}

// allSchedulers 按枚举顺序列出本包已知的全部调度器。
var allSchedulers = []Scheduler{
	SchedulerDiscrete,
	SchedulerKarras,
	SchedulerExponential,
	SchedulerAys,
	SchedulerGits,
	SchedulerSgmUniform,
	SchedulerSimple,
	SchedulerSmoothstep,
	SchedulerKlOptimal,
	SchedulerLcm,
	SchedulerBongTangent,
}

// Samplers 返回当前 Go 绑定已知的全部采样器（枚举顺序）。
// 注意：这代表“绑定支持”，不代表已加载的动态库一定全部支持，
// 动态库支持情况见 IsSamplerSupported。
func Samplers() []Sampler {
	out := make([]Sampler, len(allSamplers))
	copy(out, allSamplers)
	return out
}

// Schedulers 返回当前 Go 绑定已知的全部调度器（枚举顺序）。
// 注意：这代表“绑定支持”，不代表已加载的动态库一定全部支持，
// 动态库支持情况见 IsSchedulerSupported。
func Schedulers() []Scheduler {
	out := make([]Scheduler, len(allSchedulers))
	copy(out, allSchedulers)
	return out
}

// SamplerName 返回采样器的规范名称（如 "euler_a"）；未知值返回空字符串。
func SamplerName(s Sampler) string {
	return s.String()
}

// SchedulerName 返回调度器的规范名称（如 "karras"）；未知值返回空字符串。
func SchedulerName(s Scheduler) string {
	return s.String()
}

// IsSamplerSupported 判断当前加载的 stable-diffusion 动态库是否支持该采样器。
func IsSamplerSupported(s Sampler) bool {
	return samplerSupported(s)
}

// IsSchedulerSupported 判断当前加载的 stable-diffusion 动态库是否支持该调度器。
func IsSchedulerSupported(s Scheduler) bool {
	return schedulerSupported(s)
}

// 以下两个间接变量便于在测试中模拟“旧版库不支持新枚举”的场景。
var (
	samplerSupported   = bindings.IsSampleMethodSupported
	schedulerSupported = bindings.IsSchedulerSupported
)

// supportedSamplerNames 返回全部已知采样器名称，用于错误提示。
func supportedSamplerNames() string {
	names := make([]string, len(allSamplers))
	for i, s := range allSamplers {
		names[i] = s.String()
	}
	return strings.Join(names, ", ")
}

// supportedSchedulerNames 返回全部已知调度器名称，用于错误提示。
func supportedSchedulerNames() string {
	names := make([]string, len(allSchedulers))
	for i, s := range allSchedulers {
		names[i] = s.String()
	}
	return strings.Join(names, ", ")
}

// validateSampler 校验采样器：先判断枚举值是否已知，再判断当前动态库是否支持。
func validateSampler(s Sampler) error {
	if !s.Valid() {
		return fmt.Errorf("unknown sampler value %d; supported samplers are: %s",
			int(s), supportedSamplerNames())
	}
	if !samplerSupported(s) {
		return fmt.Errorf(
			"sampler %q is not supported by the loaded stable-diffusion library; "+
				"please upgrade the library file "+
				"(libstable-diffusion.so / libstable-diffusion.dylib / stable-diffusion.dll) "+
				"to a newer version that includes sampler %q",
			s.String(), s.String())
	}
	return nil
}

// validateScheduler 校验调度器：先判断枚举值是否已知，再判断当前动态库是否支持。
func validateScheduler(s Scheduler) error {
	if !s.Valid() {
		return fmt.Errorf("unknown scheduler value %d; supported schedulers are: %s",
			int(s), supportedSchedulerNames())
	}
	if !schedulerSupported(s) {
		return fmt.Errorf(
			"scheduler %q is not supported by the loaded stable-diffusion library; "+
				"please upgrade the library file "+
				"(libstable-diffusion.so / libstable-diffusion.dylib / stable-diffusion.dll) "+
				"to a newer version that includes scheduler %q",
			s.String(), s.String())
	}
	return nil
}

// validateSamplerConfig 校验一份采样器配置中的采样器与调度器。
func validateSamplerConfig(c SamplerConfig) error {
	if err := validateSampler(c.Method); err != nil {
		return err
	}
	if err := validateScheduler(c.Scheduler); err != nil {
		return err
	}
	return nil
}

// ParseSampler 将采样器名称（如 "euler_a"，大小写不敏感）解析为采样器枚举。
// 名称未知或当前动态库不支持时返回错误。
func ParseSampler(name string) (Sampler, error) {
	s, ok := bindings.ParseSampleMethod(name)
	if !ok {
		return s, fmt.Errorf("unknown sampler %q; supported samplers are: %s",
			name, supportedSamplerNames())
	}
	if err := validateSampler(s); err != nil {
		return s, err
	}
	return s, nil
}

// ParseScheduler 将调度器名称（如 "karras"，大小写不敏感）解析为调度器枚举。
// 名称未知或当前动态库不支持时返回错误。
func ParseScheduler(name string) (Scheduler, error) {
	s, ok := bindings.ParseScheduler(name)
	if !ok {
		return s, fmt.Errorf("unknown scheduler %q; supported schedulers are: %s",
			name, supportedSchedulerNames())
	}
	if err := validateScheduler(s); err != nil {
		return s, err
	}
	return s, nil
}
