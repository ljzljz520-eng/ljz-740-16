package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/example/stablediffusion"
	"github.com/example/stablediffusion/bindings"
)

// 测试用例2.1：系统信息获取
func TestSystemInfo(t *testing.T) {
	// 预期结果：能够获取到系统信息字符串
	info := stablediffusion.GetSystemInfo()
	if info == "" {
		t.Errorf("Expected non-empty system info, got empty string")
	}
	t.Logf("System Info: %s", info)
}

// 测试用例2.2：版本信息获取
func TestVersionInfo(t *testing.T) {
	// 预期结果：能够获取到版本信息和提交信息
	version := stablediffusion.GetVersion()
	commit := stablediffusion.GetCommit()
	
	if version == "" {
		t.Errorf("Expected non-empty version, got empty string")
	}
	
	if commit == "" {
		t.Errorf("Expected non-empty commit, got empty string")
	}
	
	t.Logf("Version: %s", version)
	t.Logf("Commit: %s", commit)
}

// 测试用例2.3：创建上下文（模拟）
func TestCreateContext(t *testing.T) {
	// 注意：这个测试需要实际的模型文件，这里只是测试API调用
	// 预期结果：在没有真实模型文件以模拟错误拦截及处理。
	modelPath := "test_model.gguf"
	
	_, err := stablediffusion.NewContext(stablediffusion.DefaultContextOptions(modelPath))
	if err == nil {
		t.Logf("Context created successfully (unexpected without real model)")
	} else {
		t.Logf("Expected error without model file: %v", err)
	}
}

// 测试用例2.4：超分辨率器创建（模拟）
func TestCreateUpscaler(t *testing.T) {
	// 注意：这个测试需要实际的超分辨率模型文件
	// 预期结果：在没有真实模型文件以模拟错误拦截及处理。
	modelPath := "test_esrgan.gguf"
	
	_, err := stablediffusion.NewUpscaler(modelPath)
	if err == nil {
		t.Logf("Upscaler created successfully (unexpected without real model)")
	} else {
		t.Logf("Expected error without model file: %v", err)
	}
}

// 测试用例2.5：模型转换（模拟）
func TestConvertModel(t *testing.T) {
	// 注意：这个测试需要实际的模型文件
	// 预期结果：能够在没有文件数据的情况下报错。
	inputPath := "test_model.gguf"
	vaePath := ""
	outputPath := "test_model_f16.gguf"
	
	err := stablediffusion.ConvertModel(inputPath, vaePath, outputPath, 1) // 1 代表 SD_TYPE_F16
	if err == nil {
		t.Logf("Model converted successfully (unexpected without real model)")
	} else {
		t.Logf("Expected error without model file: %v", err)
	}
}

// 测试用例2.6：图像生成配置验证
func TestImageGenerationConfig(t *testing.T) {
	// 测试配置结构的创建和设置
	// 预期结果：字段初始化完成无误差（例如宽高等于 512）。
	
	cfg := stablediffusion.GenerationConfig{
		Prompt:         "A beautiful cat",
		NegativePrompt: "ugly, blurry, bad art",
		Width:          512,
		Height:         512,
		Seed:           42,
		Strength:       0.8,
		BatchCount:     1,
		ClipSkip:       1,
	}
	
	// 验证配置
	if cfg.Prompt != "A beautiful cat" {
		t.Errorf("Expected prompt 'A beautiful cat', got '%s'", cfg.Prompt)
	}
	
	if cfg.Width != 512 || cfg.Height != 512 {
		t.Errorf("Expected size 512x512, got %dx%d", cfg.Width, cfg.Height)
	}
	
	t.Logf("Image generation config created successfully")
}

// 测试用例2.7：采样器/调度器字符串解析
func TestParseSamplerScheduler(t *testing.T) {
	// 预期结果：合法名称解析为对应枚举，且与 String() 往返一致
	m, err := stablediffusion.ParseSampleMethod("euler_a")
	if err != nil {
		t.Fatalf("ParseSampleMethod(euler_a) returned error: %v", err)
	}
	if m != bindings.EULER_A_SAMPLE_METHOD {
		t.Errorf("ParseSampleMethod(euler_a) = %v, want EULER_A_SAMPLE_METHOD", m)
	}

	s, err := stablediffusion.ParseScheduler("karras")
	if err != nil {
		t.Fatalf("ParseScheduler(karras) returned error: %v", err)
	}
	if s != bindings.KARRAS_SCHEDULER {
		t.Errorf("ParseScheduler(karras) = %v, want KARRAS_SCHEDULER", s)
	}
}

// 测试用例2.8：未知采样器/调度器返回错误
func TestParseUnknownSamplerScheduler(t *testing.T) {
	// 预期结果：返回错误，错误信息中包含支持的取值列表
	_, err := stablediffusion.ParseSampleMethod("not_a_sampler")
	var unknownSampler *bindings.UnknownSampleMethodError
	if !errors.As(err, &unknownSampler) {
		t.Fatalf("error type = %T, want *bindings.UnknownSampleMethodError", err)
	}
	if !strings.Contains(err.Error(), "euler_a") {
		t.Errorf("error message should list supported samplers, got: %v", err)
	}

	_, err = stablediffusion.ParseScheduler("not_a_scheduler")
	var unknownScheduler *bindings.UnknownSchedulerError
	if !errors.As(err, &unknownScheduler) {
		t.Fatalf("error type = %T, want *bindings.UnknownSchedulerError", err)
	}
	if !strings.Contains(err.Error(), "karras") {
		t.Errorf("error message should list supported schedulers, got: %v", err)
	}
}

// 测试用例2.9：GenerateImage 对未知采样器/调度器返回错误
func TestGenerateImageValidatesSampler(t *testing.T) {
	// 预期结果：在调用底层库之前就拦截非法枚举值并返回错误
	ctx := &stablediffusion.Context{}

	_, err := ctx.GenerateImage(stablediffusion.GenerationConfig{
		Prompt: "test",
		Sampler: stablediffusion.SamplerConfig{
			Method:    bindings.SampleMethod(99),
			Scheduler: bindings.KARRAS_SCHEDULER,
		},
	})
	var unknownSampler *bindings.UnknownSampleMethodError
	if !errors.As(err, &unknownSampler) {
		t.Errorf("GenerateImage with invalid method error = %v, want UnknownSampleMethodError", err)
	}

	_, err = ctx.GenerateImage(stablediffusion.GenerationConfig{
		Prompt: "test",
		Sampler: stablediffusion.SamplerConfig{
			Method:    bindings.EULER_A_SAMPLE_METHOD,
			Scheduler: bindings.Scheduler(-1),
		},
	})
	var unknownScheduler *bindings.UnknownSchedulerError
	if !errors.As(err, &unknownScheduler) {
		t.Errorf("GenerateImage with invalid scheduler error = %v, want UnknownSchedulerError", err)
	}
}

// 测试用例2.10：支持项查询
func TestSupportedSamplersSchedulers(t *testing.T) {
	// 预期结果：Mock 实现下返回全部枚举值
	if got := len(stablediffusion.SupportedSampleMethods()); got != int(bindings.SAMPLE_METHOD_COUNT) {
		t.Errorf("len(SupportedSampleMethods()) = %d, want %d", got, int(bindings.SAMPLE_METHOD_COUNT))
	}
	if got := len(stablediffusion.SupportedSchedulers()); got != int(bindings.SCHEDULER_COUNT) {
		t.Errorf("len(SupportedSchedulers()) = %d, want %d", got, int(bindings.SCHEDULER_COUNT))
	}
}
