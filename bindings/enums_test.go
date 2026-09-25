package bindings

import "testing"

func TestSampleMethodString(t *testing.T) {
	cases := map[SampleMethod]string{
		EULER_SAMPLE_METHOD:         "euler",
		EULER_A_SAMPLE_METHOD:       "euler_a",
		RES_MULTISTEP_SAMPLE_METHOD: "res_multistep",
		RES_2S_SAMPLE_METHOD:        "res_2s",
	}
	for m, want := range cases {
		if got := m.String(); got != want {
			t.Errorf("(%d).String() = %q, want %q", m, got, want)
		}
	}

	// 越界值
	if (SAMPLE_METHOD_COUNT).String() != "" {
		t.Error("SAMPLE_METHOD_COUNT.String() should be empty")
	}
	if SampleMethod(-1).String() != "" {
		t.Error("(-1).String() should be empty")
	}
}

func TestSchedulerString(t *testing.T) {
	cases := map[Scheduler]string{
		DISCRETE_SCHEDULER:     "discrete",
		KARRAS_SCHEDULER:       "karras",
		BONG_TANGENT_SCHEDULER: "bong_tangent",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("(%d).String() = %q, want %q", s, got, want)
		}
	}
	if (SCHEDULER_COUNT).String() != "" {
		t.Error("SCHEDULER_COUNT.String() should be empty")
	}
}

func TestValidRanges(t *testing.T) {
	if !EULER_SAMPLE_METHOD.Valid() || !RES_2S_SAMPLE_METHOD.Valid() {
		t.Error("edge samplers should be valid")
	}
	if SAMPLE_METHOD_COUNT.Valid() {
		t.Error("SAMPLE_METHOD_COUNT should be invalid")
	}
	if !BONG_TANGENT_SCHEDULER.Valid() {
		t.Error("last scheduler should be valid")
	}
	if SCHEDULER_COUNT.Valid() {
		t.Error("SCHEDULER_COUNT should be invalid")
	}
}

func TestParseSampleMethod(t *testing.T) {
	m, ok := ParseSampleMethod("Euler_A")
	if !ok || m != EULER_A_SAMPLE_METHOD {
		t.Errorf("ParseSampleMethod(Euler_A) = %d,%v", m, ok)
	}
	if _, ok := ParseSampleMethod("???"); ok {
		t.Error("ParseSampleMethod(???) should fail")
	}
}

func TestParseScheduler(t *testing.T) {
	s, ok := ParseScheduler("karras")
	if !ok || s != KARRAS_SCHEDULER {
		t.Errorf("ParseScheduler(karras) = %d,%v", s, ok)
	}
	if _, ok := ParseScheduler("???"); ok {
		t.Error("ParseScheduler(???) should fail")
	}
}

// Mock 库支持全部已知枚举，且不支持越界值。
func TestSupportDetectionWithMock(t *testing.T) {
	if !IsSampleMethodSupported(EULER_A_SAMPLE_METHOD) {
		t.Error("euler_a should be supported by mock library")
	}
	if IsSampleMethodSupported(SAMPLE_METHOD_COUNT) {
		t.Error("out-of-range sampler should not be supported")
	}
	if !IsSchedulerSupported(KARRAS_SCHEDULER) {
		t.Error("karras should be supported by mock library")
	}
	if IsSchedulerSupported(SCHEDULER_COUNT) {
		t.Error("out-of-range scheduler should not be supported")
	}
}

// 枚举值回归：与 stable-diffusion.h 对齐
func TestEnumValuesMatchCHeader(t *testing.T) {
	samplers := []SampleMethod{
		EULER_SAMPLE_METHOD, EULER_A_SAMPLE_METHOD, HEUN_SAMPLE_METHOD,
		DPM2_SAMPLE_METHOD, DPMPP2S_A_SAMPLE_METHOD, DPMPP2M_SAMPLE_METHOD,
		DPMPP2Mv2_SAMPLE_METHOD, IPNDM_SAMPLE_METHOD, IPNDM_V_SAMPLE_METHOD,
		LCM_SAMPLE_METHOD, DDIM_TRAILING_SAMPLE_METHOD, TCD_SAMPLE_METHOD,
		RES_MULTISTEP_SAMPLE_METHOD, RES_2S_SAMPLE_METHOD,
	}
	for i, m := range samplers {
		if int(m) != i {
			t.Errorf("sampler %d has value %d, want %d", i, m, i)
		}
	}
	schedulers := []Scheduler{
		DISCRETE_SCHEDULER, KARRAS_SCHEDULER, EXPONENTIAL_SCHEDULER,
		AYS_SCHEDULER, GITS_SCHEDULER, SGM_UNIFORM_SCHEDULER,
		SIMPLE_SCHEDULER, SMOOTHSTEP_SCHEDULER, KL_OPTIMAL_SCHEDULER,
		LCM_SCHEDULER, BONG_TANGENT_SCHEDULER,
	}
	for i, s := range schedulers {
		if int(s) != i {
			t.Errorf("scheduler %d has value %d, want %d", i, s, i)
		}
	}
}
