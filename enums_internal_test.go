package stablediffusion

import (
	"strings"
	"testing"

	"github.com/example/stablediffusion/bindings"
)

// 常量数值必须与 C 头文件中的枚举一一对应。
func TestSamplerSchedulerConstantValues(t *testing.T) {
	cases := []struct {
		name string
		got  int
		want int
	}{
		{"SamplerEuler", int(SamplerEuler), int(bindings.EULER_SAMPLE_METHOD)},
		{"SamplerEulerA", int(SamplerEulerA), int(bindings.EULER_A_SAMPLE_METHOD)},
		{"SamplerRes2s", int(SamplerRes2s), int(bindings.RES_2S_SAMPLE_METHOD)},
		{"SchedulerDiscrete", int(SchedulerDiscrete), int(bindings.DISCRETE_SCHEDULER)},
		{"SchedulerBongTangent", int(SchedulerBongTangent), int(bindings.BONG_TANGENT_SCHEDULER)},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}

	// iota 回归保护：Euler 必须为 0（曾经因为多个枚举共用一个 const 块导致为 4）
	if int(bindings.EULER_SAMPLE_METHOD) != 0 {
		t.Errorf("EULER_SAMPLE_METHOD = %d, want 0", int(bindings.EULER_SAMPLE_METHOD))
	}
	if int(bindings.DISCRETE_SCHEDULER) != 0 {
		t.Errorf("DISCRETE_SCHEDULER = %d, want 0", int(bindings.DISCRETE_SCHEDULER))
	}
	if int(bindings.SAMPLE_METHOD_COUNT) != len(Samplers()) {
		t.Errorf("SAMPLE_METHOD_COUNT = %d, want %d",
			int(bindings.SAMPLE_METHOD_COUNT), len(Samplers()))
	}
	if int(bindings.SCHEDULER_COUNT) != len(Schedulers()) {
		t.Errorf("SCHEDULER_COUNT = %d, want %d",
			int(bindings.SCHEDULER_COUNT), len(Schedulers()))
	}
}

func TestSamplerSchedulerNames(t *testing.T) {
	if SamplerName(SamplerEulerA) != "euler_a" {
		t.Errorf("SamplerName(EulerA) = %q, want euler_a", SamplerName(SamplerEulerA))
	}
	if SamplerName(SamplerRes2s) != "res_2s" {
		t.Errorf("SamplerName(Res2s) = %q, want res_2s", SamplerName(SamplerRes2s))
	}
	if SchedulerName(SchedulerBongTangent) != "bong_tangent" {
		t.Errorf("SchedulerName(BongTangent) = %q, want bong_tangent",
			SchedulerName(SchedulerBongTangent))
	}
	// 未知枚举值名称为空
	if SamplerName(Sampler(bindings.SAMPLE_METHOD_COUNT)) != "" {
		t.Error("unknown sampler should have empty name")
	}
}

func TestParseSampler(t *testing.T) {
	s, err := ParseSampler("EULER_A") // 大小写不敏感
	if err != nil {
		t.Fatalf("ParseSampler returned error: %v", err)
	}
	if s != SamplerEulerA {
		t.Errorf("ParseSampler(EULER_A) = %d, want %d", s, SamplerEulerA)
	}

	if _, err := ParseSampler("not-a-sampler"); err == nil {
		t.Error("expected error for unknown sampler name")
	} else if !strings.Contains(err.Error(), "unknown sampler") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseScheduler(t *testing.T) {
	s, err := ParseScheduler("karras")
	if err != nil {
		t.Fatalf("ParseScheduler returned error: %v", err)
	}
	if s != SchedulerKarras {
		t.Errorf("ParseScheduler(karras) = %d, want %d", s, SchedulerKarras)
	}

	if _, err := ParseScheduler("nope"); err == nil {
		t.Error("expected error for unknown scheduler name")
	}
}

func TestValidateSamplerUnknownValue(t *testing.T) {
	err := validateSampler(Sampler(9999))
	if err == nil {
		t.Fatal("expected error for unknown sampler value")
	}
	if !strings.Contains(err.Error(), "unknown sampler") {
		t.Errorf("error should mention unknown sampler, got: %v", err)
	}
}

func TestValidateSchedulerUnknownValue(t *testing.T) {
	err := validateScheduler(Scheduler(9999))
	if err == nil {
		t.Fatal("expected error for unknown scheduler value")
	}
	if !strings.Contains(err.Error(), "unknown scheduler") {
		t.Errorf("error should mention unknown scheduler, got: %v", err)
	}
}

// 模拟“Go 绑定认识该枚举、但旧版动态库不支持”的场景，错误信息必须提示升级库。
func TestValidateSamplerOldLibrary(t *testing.T) {
	orig := samplerSupported
	samplerSupported = func(s bindings.SampleMethod) bool { return false }
	defer func() { samplerSupported = orig }()

	err := validateSampler(SamplerRes2s)
	if err == nil {
		t.Fatal("expected error when library does not support sampler")
	}
	if !strings.Contains(err.Error(), "upgrade the library") {
		t.Errorf("error should ask to upgrade library, got: %v", err)
	}
	if !strings.Contains(err.Error(), "res_2s") {
		t.Errorf("error should mention sampler name, got: %v", err)
	}
}

func TestValidateSchedulerOldLibrary(t *testing.T) {
	orig := schedulerSupported
	schedulerSupported = func(s bindings.Scheduler) bool { return false }
	defer func() { schedulerSupported = orig }()

	err := validateScheduler(SchedulerBongTangent)
	if err == nil {
		t.Fatal("expected error when library does not support scheduler")
	}
	if !strings.Contains(err.Error(), "upgrade the library") {
		t.Errorf("error should ask to upgrade library, got: %v", err)
	}
	if !strings.Contains(err.Error(), "bong_tangent") {
		t.Errorf("error should mention scheduler name, got: %v", err)
	}
}

func TestValidateSamplerConfig(t *testing.T) {
	// Mock 库支持全部已知枚举，合法配置应通过校验
	if err := validateSamplerConfig(SamplerConfig{
		Method:    SamplerEulerA,
		Scheduler: SchedulerKarras,
	}); err != nil {
		t.Errorf("valid config should pass validation, got: %v", err)
	}

	// 采样器非法
	err := validateSamplerConfig(SamplerConfig{
		Method:    Sampler(42),
		Scheduler: SchedulerKarras,
	})
	if err == nil || !strings.Contains(err.Error(), "unknown sampler") {
		t.Errorf("expected unknown sampler error, got: %v", err)
	}

	// 调度器非法
	err = validateSamplerConfig(SamplerConfig{
		Method:    SamplerEulerA,
		Scheduler: Scheduler(42),
	})
	if err == nil || !strings.Contains(err.Error(), "unknown scheduler") {
		t.Errorf("expected unknown scheduler error, got: %v", err)
	}
}

func TestSamplersSchedulersLists(t *testing.T) {
	if len(Samplers()) != 14 {
		t.Errorf("expected 14 samplers, got %d", len(Samplers()))
	}
	if len(Schedulers()) != 11 {
		t.Errorf("expected 11 schedulers, got %d", len(Schedulers()))
	}

	// 返回的切片应是副本，修改不影响内部列表
	s := Samplers()
	s[0] = Sampler(123)
	if Samplers()[0] != SamplerEuler {
		t.Error("Samplers() should return a defensive copy")
	}
}
