package bindings

import (
	"errors"
	"strings"
	"testing"
)

// 枚举值必须与 stable-diffusion.h 中的 enum sample_method_t 完全一致
func TestSampleMethodEnumValues(t *testing.T) {
	cases := []struct {
		method SampleMethod
		want   int
	}{
		{EULER_SAMPLE_METHOD, 0},
		{EULER_A_SAMPLE_METHOD, 1},
		{HEUN_SAMPLE_METHOD, 2},
		{DPM2_SAMPLE_METHOD, 3},
		{DPMPP2S_A_SAMPLE_METHOD, 4},
		{DPMPP2M_SAMPLE_METHOD, 5},
		{DPMPP2Mv2_SAMPLE_METHOD, 6},
		{IPNDM_SAMPLE_METHOD, 7},
		{IPNDM_V_SAMPLE_METHOD, 8},
		{LCM_SAMPLE_METHOD, 9},
		{DDIM_TRAILING_SAMPLE_METHOD, 10},
		{TCD_SAMPLE_METHOD, 11},
		{RES_MULTISTEP_SAMPLE_METHOD, 12},
		{RES_2S_SAMPLE_METHOD, 13},
		{SAMPLE_METHOD_COUNT, 14},
	}
	for _, c := range cases {
		if int(c.method) != c.want {
			t.Errorf("sample method %v = %d, want %d (must match stable-diffusion.h)", c.method, int(c.method), c.want)
		}
	}
}

// 枚举值必须与 stable-diffusion.h 中的 enum scheduler_t 完全一致
func TestSchedulerEnumValues(t *testing.T) {
	cases := []struct {
		scheduler Scheduler
		want      int
	}{
		{DISCRETE_SCHEDULER, 0},
		{KARRAS_SCHEDULER, 1},
		{EXPONENTIAL_SCHEDULER, 2},
		{AYS_SCHEDULER, 3},
		{GITS_SCHEDULER, 4},
		{SGM_UNIFORM_SCHEDULER, 5},
		{SIMPLE_SCHEDULER, 6},
		{SMOOTHSTEP_SCHEDULER, 7},
		{KL_OPTIMAL_SCHEDULER, 8},
		{LCM_SCHEDULER, 9},
		{BONG_TANGENT_SCHEDULER, 10},
		{SCHEDULER_COUNT, 11},
	}
	for _, c := range cases {
		if int(c.scheduler) != c.want {
			t.Errorf("scheduler %v = %d, want %d (must match stable-diffusion.h)", c.scheduler, int(c.scheduler), c.want)
		}
	}
}

// String 返回的名称必须与 C 库 sample_method_to_str / scheduler_to_str 一致
func TestEnumStringNames(t *testing.T) {
	if got := EULER_A_SAMPLE_METHOD.String(); got != "euler_a" {
		t.Errorf("EULER_A_SAMPLE_METHOD.String() = %q, want %q", got, "euler_a")
	}
	if got := DPMPP2Mv2_SAMPLE_METHOD.String(); got != "dpm++2mv2" {
		t.Errorf("DPMPP2Mv2_SAMPLE_METHOD.String() = %q, want %q", got, "dpm++2mv2")
	}
	if got := KARRAS_SCHEDULER.String(); got != "karras" {
		t.Errorf("KARRAS_SCHEDULER.String() = %q, want %q", got, "karras")
	}
	// 未知值不应 panic
	if got := SampleMethod(99).String(); !strings.Contains(got, "99") {
		t.Errorf("SampleMethod(99).String() = %q, want it to contain the numeric value", got)
	}
	if got := Scheduler(-1).String(); !strings.Contains(got, "-1") {
		t.Errorf("Scheduler(-1).String() = %q, want it to contain the numeric value", got)
	}
}

// 解析合法名称：与 String() 往返一致
func TestParseRoundTrip(t *testing.T) {
	for i := 0; i < int(SAMPLE_METHOD_COUNT); i++ {
		m := SampleMethod(i)
		got, err := ParseSampleMethod(m.String())
		if err != nil {
			t.Fatalf("ParseSampleMethod(%q) returned error: %v", m.String(), err)
		}
		if got != m {
			t.Errorf("ParseSampleMethod(%q) = %v, want %v", m.String(), got, m)
		}
	}
	for i := 0; i < int(SCHEDULER_COUNT); i++ {
		s := Scheduler(i)
		got, err := ParseScheduler(s.String())
		if err != nil {
			t.Fatalf("ParseScheduler(%q) returned error: %v", s.String(), err)
		}
		if got != s {
			t.Errorf("ParseScheduler(%q) = %v, want %v", s.String(), got, s)
		}
	}
}

// 传入未知值必须返回错误，且错误信息中列出支持的取值
func TestParseUnknownValues(t *testing.T) {
	_, err := ParseSampleMethod("not_a_sampler")
	var unknownSampler *UnknownSampleMethodError
	if !errors.As(err, &unknownSampler) {
		t.Fatalf("ParseSampleMethod(unknown) error type = %T, want *UnknownSampleMethodError", err)
	}
	if !strings.Contains(err.Error(), "not_a_sampler") || !strings.Contains(err.Error(), "euler_a") {
		t.Errorf("error message should contain the bad value and supported list, got: %v", err)
	}

	_, err = ParseScheduler("not_a_scheduler")
	var unknownScheduler *UnknownSchedulerError
	if !errors.As(err, &unknownScheduler) {
		t.Fatalf("ParseScheduler(unknown) error type = %T, want *UnknownSchedulerError", err)
	}
	if !strings.Contains(err.Error(), "not_a_scheduler") || !strings.Contains(err.Error(), "karras") {
		t.Errorf("error message should contain the bad value and supported list, got: %v", err)
	}
}

// 越界枚举值校验必须返回未知值错误
func TestValidateOutOfRange(t *testing.T) {
	var unknownSampler *UnknownSampleMethodError
	if err := SampleMethod(99).Validate(); !errors.As(err, &unknownSampler) {
		t.Errorf("SampleMethod(99).Validate() error = %v, want UnknownSampleMethodError", err)
	}
	if err := SampleMethod(-1).Validate(); !errors.As(err, &unknownSampler) {
		t.Errorf("SampleMethod(-1).Validate() error = %v, want UnknownSampleMethodError", err)
	}

	var unknownScheduler *UnknownSchedulerError
	if err := Scheduler(99).Validate(); !errors.As(err, &unknownScheduler) {
		t.Errorf("Scheduler(99).Validate() error = %v, want UnknownSchedulerError", err)
	}
	if err := Scheduler(-1).Validate(); !errors.As(err, &unknownScheduler) {
		t.Errorf("Scheduler(-1).Validate() error = %v, want UnknownSchedulerError", err)
	}
}

// 模拟底层库版本过旧：只支持前 10 个采样方法、前 8 个调度器。
// 此时使用较新的枚举值，错误信息必须提示升级库文件。
func TestLibraryTooOld(t *testing.T) {
	origSampleMethodName := sdSampleMethodName
	origSchedulerName := sdSchedulerName
	defer func() {
		sdSampleMethodName = origSampleMethodName
		sdSchedulerName = origSchedulerName
		resetLibrarySupportCache()
	}()

	sdSampleMethodName = func(m SampleMethod) *byte {
		if m >= 0 && m < 10 { // 旧库只支持到 LCM_SAMPLE_METHOD
			return CString(sampleMethodNames[m])
		}
		return CString(noneStr)
	}
	sdSchedulerName = func(s Scheduler) *byte {
		if s >= 0 && s < 8 { // 旧库只支持到 SMOOTHSTEP_SCHEDULER
			return CString(schedulerNames[s])
		}
		return CString(noneStr)
	}
	resetLibrarySupportCache()

	// 旧库支持的值仍然可用
	if err := LCM_SAMPLE_METHOD.Validate(); err != nil {
		t.Errorf("LCM_SAMPLE_METHOD should be supported by the old library, got: %v", err)
	}

	// 旧库不支持的采样方法：错误需提示升级库文件
	err := TCD_SAMPLE_METHOD.Validate()
	var unsupportedSampler *UnsupportedSampleMethodError
	if !errors.As(err, &unsupportedSampler) {
		t.Fatalf("TCD_SAMPLE_METHOD.Validate() error type = %T, want *UnsupportedSampleMethodError", err)
	}
	if !strings.Contains(err.Error(), "upgrade") {
		t.Errorf("error message should mention upgrading the library, got: %v", err)
	}
	if !strings.Contains(err.Error(), "tcd") {
		t.Errorf("error message should contain the sampler name, got: %v", err)
	}

	// Parse 路径同样应返回升级提示
	_, err = ParseSampleMethod("res_2s")
	if !errors.As(err, &unsupportedSampler) {
		t.Fatalf("ParseSampleMethod(res_2s) error type = %T, want *UnsupportedSampleMethodError", err)
	}
	if !strings.Contains(err.Error(), "upgrade") {
		t.Errorf("error message should mention upgrading the library, got: %v", err)
	}

	// 旧库不支持的调度器：错误需提示升级库文件
	err = BONG_TANGENT_SCHEDULER.Validate()
	var unsupportedScheduler *UnsupportedSchedulerError
	if !errors.As(err, &unsupportedScheduler) {
		t.Fatalf("BONG_TANGENT_SCHEDULER.Validate() error type = %T, want *UnsupportedSchedulerError", err)
	}
	if !strings.Contains(err.Error(), "upgrade") {
		t.Errorf("error message should mention upgrading the library, got: %v", err)
	}

	// 支持项查询应只返回旧库支持的子集
	if got := len(SupportedSampleMethods()); got != 10 {
		t.Errorf("len(SupportedSampleMethods()) = %d, want 10", got)
	}
	if got := len(SupportedSchedulers()); got != 8 {
		t.Errorf("len(SupportedSchedulers()) = %d, want 8", got)
	}
}

// Mock 实现下支持项查询应返回全部枚举值
func TestSupportedLists(t *testing.T) {
	resetLibrarySupportCache()
	if got := len(SupportedSampleMethods()); got != int(SAMPLE_METHOD_COUNT) {
		t.Errorf("len(SupportedSampleMethods()) = %d, want %d", got, int(SAMPLE_METHOD_COUNT))
	}
	if got := len(SupportedSchedulers()); got != int(SCHEDULER_COUNT) {
		t.Errorf("len(SupportedSchedulers()) = %d, want %d", got, int(SCHEDULER_COUNT))
	}
	// 每个支持项都应通过校验
	for _, m := range SupportedSampleMethods() {
		if err := m.Validate(); err != nil {
			t.Errorf("supported sample method %v failed validation: %v", m, err)
		}
	}
	for _, s := range SupportedSchedulers() {
		if err := s.Validate(); err != nil {
			t.Errorf("supported scheduler %v failed validation: %v", s, err)
		}
	}
}

// Mock 的名称转换函数应与真实库行为一致
func TestMockNameConversion(t *testing.T) {
	if got := GetSampleMethodName(EULER_A_SAMPLE_METHOD); got != "euler_a" {
		t.Errorf("GetSampleMethodName(EULER_A) = %q, want %q", got, "euler_a")
	}
	if got := GetSchedulerName(KARRAS_SCHEDULER); got != "karras" {
		t.Errorf("GetSchedulerName(KARRAS) = %q, want %q", got, "karras")
	}
	// 越界值返回 "NONE"（与 C 库 NONE_STR 一致）
	if got := GetSampleMethodName(SAMPLE_METHOD_COUNT); got != noneStr {
		t.Errorf("GetSampleMethodName(COUNT) = %q, want %q", got, noneStr)
	}
	// 未知字符串返回 *_COUNT（与 C 库 str_to_* 行为一致）
	if got := StrToSampleMethod("bogus"); got != SAMPLE_METHOD_COUNT {
		t.Errorf("StrToSampleMethod(bogus) = %v, want SAMPLE_METHOD_COUNT", got)
	}
	if got := StrToScheduler("bogus"); got != SCHEDULER_COUNT {
		t.Errorf("StrToScheduler(bogus) = %v, want SCHEDULER_COUNT", got)
	}
	if got := StrToSampleMethod("euler_a"); got != EULER_A_SAMPLE_METHOD {
		t.Errorf("StrToSampleMethod(euler_a) = %v, want EULER_A_SAMPLE_METHOD", got)
	}
}
