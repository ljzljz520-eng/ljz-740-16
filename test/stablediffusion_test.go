package test

import (
	"strings"
	"testing"

	"github.com/example/stablediffusion"
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

// 测试用例2.7：未知采样器枚举值必须在调用底层库之前被拒绝
func TestGenerateImageRejectsUnknownSampler(t *testing.T) {
	// ctx 为 nil 也没关系：校验发生在任何底层调用之前，且不需要上下文
	var ctx *stablediffusion.Context

	cfg := stablediffusion.GenerationConfig{
		Prompt: "test",
		Sampler: stablediffusion.SamplerConfig{
			Method:    stablediffusion.Sampler(9999),
			Scheduler: stablediffusion.SchedulerKarras,
		},
	}

	_, err := ctx.GenerateImage(cfg)
	if err == nil {
		t.Fatal("expected error for unknown sampler, got nil")
	}
	if !strings.Contains(err.Error(), "unknown sampler") {
		t.Errorf("error should mention unknown sampler, got: %v", err)
	}
}

// 测试用例2.8：未知调度器枚举值必须在调用底层库之前被拒绝
func TestGenerateImageRejectsUnknownScheduler(t *testing.T) {
	var ctx *stablediffusion.Context

	cfg := stablediffusion.GenerationConfig{
		Prompt: "test",
		Sampler: stablediffusion.SamplerConfig{
			Method:    stablediffusion.SamplerEulerA,
			Scheduler: stablediffusion.Scheduler(9999),
		},
	}

	_, err := ctx.GenerateImage(cfg)
	if err == nil {
		t.Fatal("expected error for unknown scheduler, got nil")
	}
	if !strings.Contains(err.Error(), "unknown scheduler") {
		t.Errorf("error should mention unknown scheduler, got: %v", err)
	}
}

// 测试用例2.9：GenerateVideo 同样拒绝未知采样器
func TestGenerateVideoRejectsUnknownSampler(t *testing.T) {
	var ctx *stablediffusion.Context

	cfg := stablediffusion.VideoGenerationConfig{
		Prompt: "test",
		Sampler: stablediffusion.SamplerConfig{
			Method:    stablediffusion.Sampler(-1),
			Scheduler: stablediffusion.SchedulerKarras,
		},
	}

	_, err := ctx.GenerateVideo(cfg)
	if err == nil {
		t.Fatal("expected error for unknown sampler, got nil")
	}
	if !strings.Contains(err.Error(), "unknown sampler") {
		t.Errorf("error should mention unknown sampler, got: %v", err)
	}
}

// 测试用例2.10：字符串解析枚举
func TestParseSamplerAndScheduler(t *testing.T) {
	sampler, err := stablediffusion.ParseSampler("euler_a")
	if err != nil {
		t.Fatalf("ParseSampler failed: %v", err)
	}
	if sampler != stablediffusion.SamplerEulerA {
		t.Errorf("ParseSampler(euler_a) = %d, want %d", sampler, stablediffusion.SamplerEulerA)
	}

	scheduler, err := stablediffusion.ParseScheduler("karras")
	if err != nil {
		t.Fatalf("ParseScheduler failed: %v", err)
	}
	if scheduler != stablediffusion.SchedulerKarras {
		t.Errorf("ParseScheduler(karras) = %d, want %d", scheduler, stablediffusion.SchedulerKarras)
	}

	if _, err := stablediffusion.ParseSampler("bogus"); err == nil {
		t.Error("expected error when parsing bogus sampler")
	}
	if _, err := stablediffusion.ParseScheduler("bogus"); err == nil {
		t.Error("expected error when parsing bogus scheduler")
	}
}
